package confkai_test

import (
	"fmt"
	"testing"
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
