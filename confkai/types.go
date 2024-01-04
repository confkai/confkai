package confkai

import (
	"fmt"
	"reflect"
)

// TypeOf returns a function that returns the given type
// or an error if the given Valuer's value cannot be casted
// as the given type.
func TypeOf[T any](v Valuer) func() (T, error) {
	return func() (T, error) {
		var t T
		val, err := v.Value()
		if err != nil {
			return t, nil
		}

		t, ok := val.(T)
		if !ok {
			return t, fmt.Errorf("%v is not a %s", val, reflect.TypeOf(t).String())
		}
		return t, nil
	}
}

// AssertType returns a function that asserts that the value
// is of the given type and panics if either the valuer returns
// an error, or if the type can not be asserted.
func AssertType[T any](v Valuer) func() T {
	return func() T {
		v, e := TypeOf[T](v)()
		if e != nil {
			panic(e)
		}
		return v
	}
}
