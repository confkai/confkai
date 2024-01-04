package confkai_test

import (
	"os"
	"testing"
	"time"

	. "github.com/confkai/confkai/confkai"
)

func TestStructConfig1(t *testing.T) {

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

func TestStructConfig2(t *testing.T) {

	type MyConfig struct {
		Environment       func() string
		DatabaseName      func() string
		SlowMessage       func() string
		SlowMessageCached func() string
	}

	environment := "environment"

	os.Setenv("my_env", "staging")
	var config = MyConfig{
		Environment: RegisterTag(environment, Value(os.Getenv("my_env"))).Must(),
		DatabaseName: FirstOf(
			Tag(environment, "dev", Value("my_dev_db")),
			Tag(environment, "staging", Value("my_staging_db")),
			Tag(environment, "prod", Value("my_prod_db")),
		).Must(),
		SlowMessage: FuncValue(func() (string, error) {
			time.Sleep(time.Second)
			return "hello world", nil
		}).Must(),
		SlowMessageCached: Cached(FuncValue(func() (string, error) {
			time.Sleep(time.Second)
			return "hello universe", nil
		})).Must(),
	}

	if config.Environment() != "staging" {
		t.Error("Environment should equal staging")
	}

	if config.DatabaseName() != "my_staging_db" {
		t.Error("Environment should equal my_staging_db")
	}

	start := time.Now()
	if config.SlowMessage() != "hello world" {
		t.Errorf("message should equal 'hello world'")
	}
	if time.Since(start) < time.Second {
		t.Errorf("SlowMessage() should have blocked for at least 1 second")
	}

	start = time.Now()
	if config.SlowMessageCached() != "hello universe" {
		t.Errorf("message should equal 'hello universe'")
	}
	if time.Since(start) < time.Second {
		t.Errorf("SlowMessageCached() should have blocked for at least 1 second")
	}

	start = time.Now()
	if config.SlowMessageCached() != "hello universe" {
		t.Errorf("message should equal 'hello universe'")
	}
	if time.Since(start) > time.Millisecond {
		t.Errorf("SlowMessageCached() should not have blocked")
	}

}
