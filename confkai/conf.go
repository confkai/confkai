package confkai

import (
	"errors"
	"fmt"
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
// whose Valuer does not return an error
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
		val, ok := cachedMap.Load(value)
		if ok {
			// Cache hit
			return val, nil
		}
		// Cache miss
		val, err := value.Value()
		if err != nil {
			return nil, err
		}
		cachedMap.Store(value, val)
		return val, nil
	})
}
