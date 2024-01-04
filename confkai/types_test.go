package confkai_test

import (
	"testing"

	"github.com/confkai/confkai/confkai"
)

func TestAssertType(t *testing.T) {

	err := ShouldNotPanic(t, func() {
		confkai.AssertType[string](confkai.Value("abc"))()
	})
	if err != nil {
		t.Error(err)
	}
	err = ShouldNotPanic(t, func() {
		confkai.AssertType[int](confkai.Value(123))()
	})
	if err != nil {
		t.Error(err)
	}
	err = ShouldPanic(t, func() {
		confkai.AssertType[string](confkai.Value([]byte("123")))()
	})
	if err != nil {
		t.Error(err)
	}
	err = ShouldPanic(t, func() {
		confkai.AssertType[string](confkai.Value(123))()
	})
	if err != nil {
		t.Error(err)
	}
}

func TestTypeOf(t *testing.T) {

	_, err := confkai.TypeOf[string](confkai.Value("abc"))()
	if err != nil {
		t.Errorf("should not return an error: %s", err)
	}
	_, err = confkai.TypeOf[int](confkai.Value(123))()
	if err != nil {
		t.Errorf("should not return an error: %s", err)
	}
	_, err = confkai.TypeOf[string](confkai.Value([]byte("123")))()
	if err == nil {
		t.Errorf("should return an error")
	}
	_, err = confkai.TypeOf[string](confkai.Value(123))()
	if err == nil {
		t.Errorf("should return an error")
	}
}
