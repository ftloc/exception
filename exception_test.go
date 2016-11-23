package exception_test

import (
	"github.com/ftloc/exception"

	"errors"
	"fmt"
	"path"
	"testing"
)

func ExampleTry() {
	exception.Try(func() {
		// do something
		if 1 == 2 {
			// oh noes, something is wrong
			exception.Throw("1 should not be == 2")
		}

		exception.Throw(struct{ a int }{a: 1})
	}).Catch(func(s string) {
		fmt.Printf("Caught a string: %s\n", s)
	}).CatchAll(func(i interface{}) {
		fmt.Printf("Caught something: %+v\n", i)
	}).Finally(func() {
		fmt.Println("Something might have been wrong ... who knows ...")
	})
	// Output: Caught something: {a:1}
	// Something might have been wrong ... who knows ...
}

func TestTry(t *testing.T) {
	exception.Try(func() {}).Finally(func() {})
}

func TestCatch(t *testing.T) {
	type test struct{}
	called := false
	exception.Try(func() { exception.Throw(test{}) }).Catch(func(t test) {
		called = true
	}).Finally(func() {})
	if !called {
		t.Fail()
	}

	thrown := false
	inner := func() {
		defer func() {
			if r := recover(); r != nil {
				thrown = true
			}
		}()
		exception.Try(func() {}).Catch(func() {})
	}
	inner()
	if !thrown {
		t.Fail()
	}
}

func TestCatchAll(t *testing.T) {
	exception.Try(func() { exception.Throw(1) }).CatchAll(func(i interface{}) {}).Finally(func() {})
}

func TestIgnore(t *testing.T) {
	exception.Try(func() { exception.Throw(1) }).Ignore().Finally(func() {})
}

func TestFinally(t *testing.T) {
	called := false
	exception.Try(func() {}).Finally(func() { called = true })
	if !called {
		t.Fail()
	}

	called = false
	callorder := ""
	exception.Try(func() { exception.Throw(1) }).Catch(func(i int) {
		callorder += "C"
	}).Finally(func() {
		called = true
		callorder += "F"
	})
	if !called {
		t.Fail()
	}
	if callorder != "CF" {
		t.Fail()
	}

	thrown := false
	inner := func() {
		called = false
		defer func() {
			if r := recover(); r != nil {
				thrown = true
			}
		}()
		exception.Try(func() {
			exception.Throw(1)
		}).Finally(func() { called = true })
	}
	inner()

	if !(thrown && called) {
		t.Fail()
	}
}

func TestGo(t *testing.T) {
	called := false
	exception.Try(func() {
		called = true
	}).Go()
	if !called {
		t.Fatal("The function was never called.")
	}
}

func TestThrowOnFalse(t *testing.T) {
	exception.ThrowOnFalse(true, 1)
	called := false
	exception.Try(func() { exception.ThrowOnFalse(false, 1) }).Catch(func(i int) {
		called = true
	}).Finally(func() {})
	exception.ThrowOnFalse(called, "Function was not called :(")
}

func TestThrowOnError(t *testing.T) {
	exception.ThrowOnError(nil, 1)
	called := false
	exception.Try(func() { exception.ThrowOnError(errors.New("test"), 2) }).Catch(func(i int) {
		called = true
	}).Finally(func() {})
	if !called {
		t.Fail()
	}
}

func throw1() interface{} {
	return 1
}

func throw2() interface{} {
	return 2
}

func TestThrowOnFalseFn(t *testing.T) {
	called := false
	exception.Try(func() {
		exception.ThrowOnFalseFn(true, throw1)
		exception.ThrowOnFalseFn(false, throw2)
	}).Catch(func(i int) {
		if i != 2 {
			t.Fail()
		}
		called = true
	}).Finally(func() {})
	if !called {
		t.Fail()
	}
}

func TestThrowOnErrorFn(t *testing.T) {
	called := false
	exception.Try(func() {
		exception.ThrowOnErrorFn(nil, throw1)
		exception.ThrowOnErrorFn(errors.New("test"), throw2)
	}).Catch(func(i int) {
		if i != 2 {
			t.Fail()
		}
		called = true
	}).Finally(func() {})
	if !called {
		t.Fail()
	}
}

func TestGetThrowerWithoutException(t *testing.T) {
	ok, _, _ := exception.GetThrower()
	if ok {
		t.Fatal("Got thrower without exception.")
	}
}

func TestGetThrowerInFinallyWithoutException(t *testing.T) {
	exception.Try(func() {
	}).Finally(func() {
		ok, _, _ := exception.GetThrower()
		if ok {
			t.Fatal("Got thrower without exception.")
		}
	})
}

// not sure why we disallow that
func TestGetThrowerWithoutExceptionHandler(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			ok, _, _ := exception.GetThrower()
			if ok {
				t.Fatal("Got thrower without exception.")
			}
		}
	}()
	panic("call recover")
}

func TestGetThrowerWithException(t *testing.T) {
	type ts struct{}
	exception.Try(func() {
		exception.Throw(ts{})
	}).Catch(func(ts) {
		ok, f, _ := exception.GetThrower()
		if !ok {
			t.Fatal("Got no thrower in Catch.")
		}
		_, file := path.Split(f)
		if file != "exception_test.go" {
			t.Fatal("Wrong file identified.")
		}
	}).Finally(func() {
		ok, f, _ := exception.GetThrower()
		if !ok {
			t.Fatal("Got no thrower in Finally.")
		}
		_, file := path.Split(f)
		if file != "exception_test.go" {
			t.Fatal("Wrong file identified.")
		}
	})
}
