package confkai_test

import (
	"testing"

	. "github.com/confkai/confkai/confkai"
)

func TestStructConfig(t *testing.T) {

	type TestType struct {
		Foo string
	}
	type TestStruct struct {
		A     func() string
		Afail func() string
		B     func() *string
		Bfail func() *string
		C     func() int
		Cfail func() int
		D     func() *int
		Dfail func() *int
		E     func() TestType
		Efail func() TestType
		F     func() *TestType
		Ffail func() *TestType
	}

	conf := TestStruct{
		A:     Value("A Foo").Must(),
		Afail: failuerValue("").Must(),
		B:     Value(pointerOf("A Foo")).Must(),
		Bfail: failuerValue(pointerOf("A Foo")).Must(),
		C:     Value(123).Must(),
		Cfail: failuerValue(123).Must(),
		D:     Value(pointerOf(123)).Must(),
		Dfail: failuerValue(pointerOf(123)).Must(),
		E:     Value(TestType{}).Must(),
		Efail: failuerValue(TestType{}).Must(),
		F:     Value(pointerOf(TestType{})).Must(),
		Ffail: failuerValue(pointerOf(TestType{})).Must(),
	}

	ShouldNotPanic(t, func() { conf.A() })
	ShouldNotPanic(t, func() { conf.B() })
	ShouldNotPanic(t, func() { conf.C() })
	ShouldNotPanic(t, func() { conf.D() })
	ShouldNotPanic(t, func() { conf.E() })
	ShouldNotPanic(t, func() { conf.F() })

	ShouldPanic(t, func() { conf.Afail() })
	ShouldPanic(t, func() { conf.Bfail() })
	ShouldPanic(t, func() { conf.Cfail() })
	ShouldPanic(t, func() { conf.Dfail() })
	ShouldPanic(t, func() { conf.Efail() })
	ShouldPanic(t, func() { conf.Ffail() })
}
