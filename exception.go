package exception

import (
	"reflect"
)

func CallWith(fn interface{}, parameters ...interface{}) {
	va := reflect.ValueOf(fn)
	if va.Kind() != reflect.Func {
		panic("Not a function")
	}

	s := make([]reflect.Value, 0)
	for _, p := range parameters {
		v := reflect.ValueOf(p)
		s = append(s, v)
	}

	va.Call(s)
}

type (
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
	Exception struct {
		Message string
	}
	OtherException struct {
		Message string
	}
)

func Try(mainfn func()) Tryer {
	return &tryer{
		mainfn:   mainfn,
		catches:  make(map[reflect.Type]interface{}),
		catchall: func(i interface{}) { panic(i) },
	}
}

func (t *tryer) Catch(fn interface{}) Tryer {
	va := reflect.ValueOf(fn)
	if va.Kind() != reflect.Func || va.Type().NumIn() == 0 {
		panic("Catch needs a function")
	}

	typ := va.Type().In(0)
	t.catches[typ] = fn
	return t
}

func (t *tryer) CatchAll(fn func(interface{})) Tryer {
	t.catchall = fn
	return t
}

func (to *tryer) Finally(finfn func()) {
	defer func() {
		if r := recover(); r != nil {
			t := reflect.TypeOf(r)
			fn, ok := to.catches[t]
			if !ok {
				to.catchall(r)
				return
			}
			CallWith(fn, r)
			finfn()
		}
	}()
	to.mainfn()
	finfn()
}

func Throw(i interface{}) {
	panic(i)
}

func ThrowOnFalse(b bool, i interface{}) {
	if !b {
		Throw(i)
	}
}

func ThrowOnError(e error, i interface{}) {
	if nil != e {
		Throw(i)
	}
}

func ThrowOnErrorFn(e error, f func() interface{}) {
	if nil != e {
		Throw(f())
	}
}
