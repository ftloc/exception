package exception

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/ftloc/caller"
)

type (
	// Tryer is the interface that is exposed that encapsules the packages functionality. You may exchange the implementation easily if you deem it necessarily
	Tryer interface {
		Catch(interface{}) Tryer
		CatchAll(func(interface{})) Tryer
		Ignore() Tryer
		Finally(func())
		Go()
	}
	tryer struct {
		mainfn   func()
		catches  map[reflect.Type]interface{}
		catchall func(interface{})
	}
)

var (
	pkgPath = reflect.TypeOf(tryer{}).PkgPath()
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

// Ignore catches all exceptions and silently ignores them
func (t *tryer) Ignore() Tryer {
	t.catchall = func(interface{}) {}
	return t
}

// Finally initiates the call to the tried function and is always called after
// the function was executed, no matter if an exception occurred or not.
func (t *tryer) Finally(finfn func()) {
	defer func() {
		defer finfn()
		if r := recover(); r != nil {
			t.findHandler(r)
		}
	}()
	t.mainfn()
}

// This function implements the search for the correct handler
func (t *tryer) findHandler(r interface{}) {
	tyo := reflect.TypeOf(r)

	if fn, ok := t.catches[tyo]; ok {
		caller.CallWith(fn, r)
		return
	}

	for tt, fn := range t.catches {
		if tyo.ConvertibleTo(tt) {
			caller.CallWith(fn, r)
			return
		}
	}

	t.catchall(r)
}

// Go is a shorthand version of Finally(func(){})
func (t *tryer) Go() {
	t.Finally(func() {})
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

// GetThrower scans the callstack for the last panic / exception.Throw* call and returns the filename and line, returns false as first parameter if nothing paniced
func GetThrower() (bool, string, int) {
	cs := make([]uintptr, 20)
	amount := runtime.Callers(2, cs)
	usedThrower := false
	var pfield *runtime.Func

	for i := 0; i < amount; i++ {
		f := runtime.FuncForPC(cs[i])
		if f.Name() == "runtime.gopanic" {
			pfield = f
		} else if pfield != nil && strings.HasPrefix(f.Name(), pkgPath+".") {
			usedThrower = true
		} else if pfield != nil {
			if !usedThrower {
				return false, "", 0
			}
			fi, li := f.FileLine(f.Entry())
			return true, fi, li
		}
	}
	return false, "", 0
}
