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
		A     Valuer[string]
		Afail Valuer[string]
		B     Valuer[*string]
		Bfail Valuer[*string]
		C     Valuer[int]
		Cfail Valuer[int]
		D     Valuer[*int]
		Dfail Valuer[*int]
		E     Valuer[TestType]
		Efail Valuer[TestType]
		F     Valuer[*TestType]
		Ffail Valuer[*TestType]
	}

	conf := TestStruct{
		A:     Value("A Foo"),
		Afail: failuerValue(""),
		B:     Value(pointerOf("A Foo")),
		Bfail: failuerValue(pointerOf("A Foo")),
		C:     Value(123),
		Cfail: failuerValue(123),
		D:     Value(pointerOf(123)),
		Dfail: failuerValue(pointerOf(123)),
		E:     Value(TestType{}),
		Efail: failuerValue(TestType{}),
		F:     Value(pointerOf(TestType{})),
		Ffail: failuerValue(pointerOf(TestType{})),
	}

	ShouldNotPanic(t, func() { conf.A.Must() })
	ShouldNotPanic(t, func() { conf.B.Must() })
	ShouldNotPanic(t, func() { conf.C.Must() })
	ShouldNotPanic(t, func() { conf.D.Must() })
	ShouldNotPanic(t, func() { conf.E.Must() })
	ShouldNotPanic(t, func() { conf.F.Must() })

	ShouldPanic(t, func() { conf.Afail.Must() })
	ShouldPanic(t, func() { conf.Bfail.Must() })
	ShouldPanic(t, func() { conf.Cfail.Must() })
	ShouldPanic(t, func() { conf.Dfail.Must() })
	ShouldPanic(t, func() { conf.Efail.Must() })
	ShouldPanic(t, func() { conf.Ffail.Must() })
}

func TestStructConfig2(t *testing.T) {

	type MyConfig struct {
		Environment       Valuer[string]
		DatabaseName      Valuer[string]
		SlowMessage       Valuer[string]
		SlowMessageCached Valuer[string]
	}

	environment := "environment"

	os.Setenv("my_env", "staging")
	var config = MyConfig{
		Environment: RegisterTag(environment, Value(os.Getenv("my_env"))),
		DatabaseName: FirstOf(
			Tag(environment, "dev", Value("my_dev_db")),
			Tag(environment, "staging", Value("my_staging_db")),
			Tag(environment, "prod", Value("my_prod_db")),
		),
		SlowMessage: FuncValue(func() (string, error) {
			time.Sleep(time.Second)
			return "hello world", nil
		}),
		SlowMessageCached: Cached(FuncValue(func() (string, error) {
			time.Sleep(time.Second)
			return "hello universe", nil
		})),
	}

	if config.Environment.Must() != "staging" {
		t.Error("Environment should equal staging")
	}

	if config.DatabaseName.Must() != "my_staging_db" {
		t.Error("Environment should equal my_staging_db")
	}

	start := time.Now()
	if config.SlowMessage.Must() != "hello world" {
		t.Errorf("message should equal 'hello world'")
	}
	if time.Since(start) < time.Second {
		t.Errorf("SlowMessage() should have blocked for at least 1 second")
	}

	start = time.Now()
	if config.SlowMessageCached.Must() != "hello universe" {
		t.Errorf("message should equal 'hello universe'")
	}
	if time.Since(start) < time.Second {
		t.Errorf("SlowMessageCached() should have blocked for at least 1 second")
	}

	start = time.Now()
	if config.SlowMessageCached.Must() != "hello universe" {
		t.Errorf("message should equal 'hello universe'")
	}
	if time.Since(start) > time.Millisecond {
		t.Errorf("SlowMessageCached() should not have blocked")
	}

}
