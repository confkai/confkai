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
type Valuer interface {
	Value() (any, error)
}

// ValuerFunc is a function that implements
// Valuer. You can use this to create your own
// Valuers. For Example:
//
//	confkai.ValuerFunc(func() (any, error) {
//		result := os.GetEnv("foo_bar")
//		return result, nil
//	})
type ValuerFunc func() (any, error)

// Value() is ValueFunc's implementation of Valuer
func (f ValuerFunc) Value() (any, error) {
	return f()
}

// Value returns a Valuer that always returns the
// given value.
func Value(v any) Valuer {
	return ValuerFunc(func() (any, error) {
		return v, nil
	})
}

// FirstOf returns the first value from the list
// whose Valuer does not return an error. Returns an
// error itself if no valid value is found.
func FirstOf(vals ...Valuer) Valuer {
	return ValuerFunc(func() (any, error) {
		errs := make([]error, 0)
		for _, v := range vals {
			val, err := v.Value()
			if err == nil {
				return val, nil
			}
			errs = append(errs, err)
		}
		return nil, fmt.Errorf("no valuers in FirstOf returned a value, %w", errors.Join(errs...))
	})
}

// thread safe cache storing
var cachedMap = sync.Map{}

// Cached is a simple thread-safe caching mechanism to wrap values with. The
// child value will only be called once if successful. That value wil
// be stored and returned for any subsequent calls.
func Cached(value Valuer) Valuer {
	return ValuerFunc(func() (any, error) {
		hashable := _getFunctionName(value)
		val, ok := cachedMap.Load(hashable)
		if ok {
			// Cache hit
			return val, nil
		}
		// Cache miss
		val, err := value.Value()
		if err != nil {
			return nil, err
		}
		cachedMap.Store(hashable, val)
		return val, nil
	})
}

func _getFunctionName(temp interface{}) string {
	strs := (runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name())
	return strs
}

// Valuers are lazy loading by default. Eager calls Value() on
// the passed Valuer immediately and wraps the result as another
// valuer, effectively "Eager Loading" the  passed value.
func Eager(value Valuer) Valuer {
	val, err := value.Value()
	return ValuerFunc(func() (any, error) {
		return val, err
	})
}
