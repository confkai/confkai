package confkai

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

// Valuer is the base interface for any configuration type.
// With this base interface, we can wrap and compose different
// kinds of values.
type Valuer[T any] interface {
	// Value executes and returns the value, or an error
	// if one occurs.
	Value() (T, error)
	// Must returns a function that only returns
	// the value or panics if the underlying error from
	// Value() is not nil
	Must() func() T
}

// ValuerFunc is a function that implements
// Valuer. You can use this to create your own
// Valuers. For Example:
//
//	confkai.ValuerFunc(func() (any, error) {
//		result := os.GetEnv("foo_bar")
//		return result, nil
//	})
type ValuerFunc[T any] func() (T, error)

// Value() is ValueFunc's implementation of Valuer
func (f ValuerFunc[T]) Value() (T, error) {
	return f()
}

func (f ValuerFunc[T]) Must() func() T {
	return func() T {
		val, err := f()
		if err != nil {
			panic(err)
		}
		return val
	}
}

// FuncValue turns a function into a ValuerFunc, which implements Valuer
func FuncValue[T any](fn func() (T, error)) Valuer[T] {
	return ValuerFunc[T](fn)
}

// Value returns a Valuer that always returns the
// given value.
func Value[T any](v T) Valuer[T] {
	return ValuerFunc[T](func() (T, error) {
		return v, nil
	})
}

// FirstOf returns the first value from the list
// whose Valuer does not return an error. Returns an
// error itself if no valid value is found.
func FirstOf[T any](vals ...Valuer[T]) Valuer[T] {
	return ValuerFunc[T](func() (T, error) {
		var val T
		errs := make([]error, 0)
		for _, v := range vals {
			val, err := v.Value()
			if err == nil {
				return val, nil
			}
			errs = append(errs, err)
		}
		return val, fmt.Errorf("no valuers in FirstOf returned a value, %w", errors.Join(errs...))
	})
}

// thread safe cache storing
var cachedMap = sync.Map{}

// Cached is a simple thread-safe caching mechanism to wrap values with. The
// child value will only be called once if successful. That value wil
// be stored and returned for any subsequent calls.
func Cached[T any](value Valuer[T]) Valuer[T] {
	return ValuerFunc[T](func() (T, error) {
		hashable := _getFunctionName(value)
		var castedVal T
		val, ok := cachedMap.Load(hashable)
		if ok {
			// Cache hit
			castedVal, ok := val.(T)
			if !ok {
				return castedVal, fmt.Errorf("could not cast value of %v to type of '%s'", val, reflect.TypeOf(castedVal).String())
			}
			return val.(T), nil
		}
		// Cache miss
		castedVal, err := value.Value()
		if err != nil {
			return castedVal, err
		}
		cachedMap.Store(hashable, castedVal)
		return castedVal, nil
	})
}

func _getFunctionName(temp interface{}) string {
	strs := (runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name())
	return strs
}

// Valuers are lazy loading by default. Eager calls Value() on
// the passed Valuer immediately and wraps the result as another
// valuer, effectively "Eager Loading" the  passed value.
func Eager[T any](value Valuer[T]) Valuer[T] {
	val, err := value.Value()
	return ValuerFunc[T](func() (T, error) {
		return val, err
	})
}

var tagMap = sync.Map{}

// Register tag eagerly loads value and save it for use with
// the Valuer 'Tagged'. It only accepts string valuers.
func RegisterTag(tag string, value Valuer[string]) Valuer[string] {
	val, err := value.Value() //eagerly load the value
	if err != nil {
		panic(err)
	}
	tagMap.Store(tag, val) //store it
	return ValuerFunc[string](func() (string, error) {
		return value.Value()
	})
}

// Tag returns a valuer who will return an error if the
// the passed tag for key is not set, or the value set for the tag
// does not equal val.
func Tag[T any](key, val string, value Valuer[T]) Valuer[T] {
	return ValuerFunc[T](func() (T, error) {
		var emptyVal T
		tagValue, ok := tagMap.Load(key)
		if !ok {
			return emptyVal, fmt.Errorf("value for tag '%s' is not set", key)
		}
		if tagValue != val {
			return emptyVal, fmt.Errorf("%s is '%s', not '%s'", key, tagValue, val)
		}
		return value.Value()
	})
}
