package exception

import (
	"github.com/ftloc/caller"
	"reflect"
	"strconv"
)

type (
	// Tryer is the interface that is exposed that encapsules the packages functionality. You may exchange the implementation easily if you deem it necessarily
	Tryer interface {
		Catch(interface{}) Tryer
		CatchAll(func(interface{})) Tryer
		Finally(func())
	}
	tryer struct {
		mainfn   func()
		catches  map[reflect.Type]interface{}
		catchall func(interface{})
	}
)

// Try creates a Tryer object. The given function will be called, when finally on the Tryer object is called.
func Try(mainfn func()) Tryer {
	return &tryer{
		mainfn:   mainfn,
		catches:  make(map[reflect.Type]interface{}),
		catchall: func(i interface{}) { panic(i) },
	}
}

// Catch catches exceptions of the type that the given function takes as a first (and only) argument.
func (t *tryer) Catch(fn interface{}) Tryer {
	va := reflect.ValueOf(fn)
	if va.Kind() != reflect.Func || va.Type().NumIn() != 1 {
		panic("Catch needs a function with exactly one parameter (got: " + va.Type().String() + " with " + strconv.Itoa(va.Type().NumIn()) + " arguments)")
	}

	typ := va.Type().In(0)
	t.catches[typ] = fn
	return t
}

// CatchAll catches all exceptions and panics that occur within the tried function, that are not specifically caught.
func (t *tryer) CatchAll(fn func(interface{})) Tryer {
	t.catchall = fn
	return t
}

// Finally initiates the call to the tried function and is always called after
// the function was executed, no matter if an exception occured or not.
func (t *tryer) Finally(finfn func()) {
	defer func() {
		defer finfn()
		if r := recover(); r != nil {
			tyo := reflect.TypeOf(r)
			fn, ok := t.catches[tyo]
			if !ok {
				t.catchall(r)
				return
			}
			caller.CallWith(fn, r)
		}
	}()
	t.mainfn()
}

// Throw an exception. Any type qualifies as an exception.
func Throw(i interface{}) {
	panic(i)
}

// ThrowOnFalse throws an exception if the bool b equals false
func ThrowOnFalse(b bool, i interface{}) {
	if !b {
		Throw(i)
	}
}

// ThrowOnFalseFn throws an exception (produced by f) if the bool b equals false
func ThrowOnFalseFn(b bool, f func() interface{}) {
	if !b {
		Throw(f())
	}
}

// ThrowOnError throws an exception if e is not nil
func ThrowOnError(e error, i interface{}) {
	if nil != e {
		Throw(i)
	}
}

// ThrowOnErrorFn throws an exception (produced by f) if e is not nil
func ThrowOnErrorFn(e error, f func() interface{}) {
	if nil != e {
		Throw(f())
	}
}
