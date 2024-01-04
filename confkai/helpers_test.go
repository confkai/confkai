package confkai_test

import (
	"fmt"
	"testing"

	"github.com/confkai/confkai/confkai"
)

func ShouldPanic(t *testing.T, fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = nil
		}
	}()

	fn()
	return fmt.Errorf("The code did not panic")
}

func ShouldNotPanic(t *testing.T, fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("The code did panicked")
		}
	}()
	fn()
	return nil
}

func pointerOf[T any](o T) *T {
	return &o
}

func failuerValue[T any](val T) confkai.Valuer[T] {
	return confkai.FuncValue(func() (t T, err error) {
		return t, fmt.Errorf("failed value")
	})
}
