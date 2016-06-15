package exception

import (
	"errors"
	"testing"
)

func TestTry(t *testing.T) {
	Try(func() {}).Finally(func() {})
}

func TestCatch(t *testing.T) {
	type test struct{}
	called := false
	Try(func() { Throw(test{}) }).Catch(func(t test) {
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
		Try(func() {}).Catch(func() {})
	}
	inner()
	if !thrown {
		t.Fail()
	}
}

func TestCatchAll(t *testing.T) {
	Try(func() { Throw(1) }).CatchAll(func(i interface{}) {}).Finally(func() {})
}

func TestFinally(t *testing.T) {
	called := false
	Try(func() {}).Finally(func() { called = true })
	if !called {
		t.Fail()
	}

	called = false
	callorder := ""
	Try(func() { Throw(1) }).Catch(func(i int) {
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
		Try(func() {
			Throw(1)
		}).Finally(func() { called = true })
	}
	inner()

	if !(thrown && called) {
		t.Fail()
	}
}

func TestThrowOnFalse(t *testing.T) {
	ThrowOnFalse(true, 1)
	called := false
	Try(func() { ThrowOnFalse(false, 1) }).Catch(func(i int) {
		called = true
	}).Finally(func() {})
	ThrowOnFalse(called, "Function was not called :(")
}

func TestThrowOnError(t *testing.T) {
	ThrowOnError(nil, 1)
	called := false
	Try(func() { ThrowOnError(errors.New("test"), 2) }).Catch(func(i int) {
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
	Try(func() {
		ThrowOnFalseFn(true, throw1)
		ThrowOnFalseFn(false, throw2)
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
	Try(func() {
		ThrowOnErrorFn(nil, throw1)
		ThrowOnErrorFn(errors.New("test"), throw2)
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
