package confkai_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/confkai/confkai/confkai"
)

func TestValuerFunc(t *testing.T) {

	result, err := confkai.FuncValue(func() (string, error) {
		return "test", nil
	}).Value()
	if result != "test" {
		t.Error("result should equal test")
	}
	if err != nil {
		t.Error("err should be nil")
	}

	_, err = confkai.FuncValue(func() (string, error) {
		return "", fmt.Errorf("an error occurred")
	}).Value()
	if err == nil {
		t.Error("err should be not be nil")
	}

}

func TestValue(t *testing.T) {
	result, err := confkai.Value("hello world").Value()
	if err != nil {
		t.Errorf("err should be nil: %s", err)
	}
	if result != "hello world" {
		t.Error("result should be 'hello world'")
	}
}

func TestFirstOf(t *testing.T) {
	// When the first element errors
	result, err := confkai.FirstOf[string](
		confkai.FuncValue(func() (string, error) {
			return "", fmt.Errorf("test value error")
		}),
		confkai.Value("hello world"),
	).Value()
	if err != nil {
		t.Errorf("err should be nil: %s", err)
	}
	if result != "hello world" {
		t.Error("result should be 'hello world'")
	}

	// when the second element errors
	result, err = confkai.FirstOf[string](
		confkai.Value("hello world"),
		confkai.FuncValue(func() (string, error) {
			return "", fmt.Errorf("test value error")
		}),
	).Value()
	if err != nil {
		t.Errorf("err should be nil: %s", err)
	}
	if result != "hello world" {
		t.Error("result should be 'hello world'")
	}

	// when all elements error
	_, err = confkai.FirstOf[string](
		confkai.FuncValue(func() (string, error) {
			return "", fmt.Errorf("test value error")
		}),
		confkai.FuncValue(func() (string, error) {
			return "", fmt.Errorf("test value error")
		}),
	).Value()
	if err == nil {
		t.Errorf("err should not be nil")
	}

}

func TestCached(t *testing.T) {
	delay := time.Duration(3 * time.Second)
	aLongRunningValue := confkai.FuncValue(func() (string, error) {
		time.Sleep(delay)
		return "hello world", nil
	})

	start := time.Now()
	result, err := confkai.Cached[string](aLongRunningValue).Value()
	if err != nil {
		t.Errorf("err should be nil: %s", err)
	}
	if result != "hello world" {
		t.Error("result should be 'hello world'")
	}
	if time.Since(start) < (delay) {
		t.Error("Should have taken more than 5 seconds")
	}
	// test running again, it should be cached now
	start = time.Now()
	result, err = confkai.Cached[string](aLongRunningValue).Value()
	if err != nil {
		t.Errorf("err should be nil: %s", err)
	}
	if result != "hello world" {
		t.Error("result should be 'hello world'")
	}
	if time.Since(start) > (delay) {
		t.Error("Should have taken less than 5 seconds")
	}
}

func TestEager(t *testing.T) {

	delay := time.Duration(3 * time.Second)
	aLongRunningValuer := confkai.FuncValue(func() (string, error) {
		time.Sleep(delay)
		return "hello world", nil
	})

	start := time.Now()
	// since this is eager loading it should basically block for the
	// time that valuer takes to evaluate.
	myEagerConf := confkai.Eager[string](aLongRunningValuer)
	if time.Since(start) < delay {
		t.Error("Eager did not block as expected")
	}
	result, err := myEagerConf.Value()
	if err != nil {
		t.Errorf("err should be nil: %s", err)
	}
	if result != "hello world" {
		t.Error("result should be 'hello world'")
	}
}
