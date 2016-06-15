package exception_test

import (
	"github.com/ftloc/exception"

	"errors"
	"testing"
)

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
