package confkai

import (
	"errors"
	"fmt"
	"reflect"
)

type Valuer interface {
	Value() (any, error)
}

type ValuerFunc func() (any, error)

func (f ValuerFunc) Value() (any, error) {
	return f()
}

func Value(v any) Valuer {
	return ValuerFunc(func() (any, error) {
		return v, nil
	})
}

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

func TypeOf[T any](v Valuer) T {
	val, err := v.Value()
	if err != nil {
		panic(err)
	}

	var t T
	t, ok := val.(T)
	if !ok {
		panic(fmt.Sprintf("%v is not a %s", val, reflect.TypeOf(t).String()))
	}
	return t
}
