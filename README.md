[![GoDoc](https://godoc.org/github.com/ftloc/exception?status.svg)](https://godoc.org/github.com/ftloc/exception)
[![Build Status](https://travis-ci.org/ftloc/exception.svg?branch=master)](https://travis-ci.org/ftloc/exception)
[![Coverage Status](https://coveralls.io/repos/github/ftloc/exception/badge.svg?branch=master)](https://coveralls.io/github/ftloc/exception?branch=master)
# exception
exception handling for golang

# usage
## basics
Basically Try/Catch are just some function wrappers that recover from panics.
```go
exception.Try(func(){
	// something that might go wrong
    exception.Throw(fmt.Errorf("Some error"))
}).Catch(func(e error){
	log.Warningf("An error occured: %s", e)
}).Go()
```

## multiple errors
The handler function is chosen by the function signature, that takes the the given object as first and only argument.
```go
type Exception struct{
	Message string
}
exception.Try(func(){
	// something that might go wrong
    exception.Throw(Exception{Message:"This is an exception"})
}).Catch(func(e Exception){
	log.Warningf("An exception occured: %s", e.Message)
}).Catch(func(e error){
	log.Warningf("An error occured: %s", e)
}).Go()
```

## unknown errors
If you call into other people's code and don't know what they might throw in a future version of their code, you can catch that with a catchall function.
```go
exception.Try(func(){
	// something that might go wrong
    exception.Throw(fmt.Errorf("Some error"))
}).CatchAll(func(e interface{}){
	log.Warningf("An error occured: %+v", e)
}).Go()
```

## panic
Sometimes you have some code that calls panic.
```go
exception.Try(func(){
	// something that might go wrong
    panic(fmt.Errorf("Some error"))
}).CatchAll(func(e interface{}){
	log.Warningf("An error occured: %+v", e)
}).Go()
```

## ignore
You just want to try over and over?
```go
stop := false
for !stop {
	exception.Try(func(){
		// something that might go wrong
	    panic(fmt.Errorf("Some error"))
        stop = true
	}).Ignore().Go()
}
```

## finally
If you need to run some cleanup code anyways
```go
exception.Try(func(){
	// something that might go wrong
    panic(fmt.Errorf("Some error"))
}).Ignore().Finally(func(){
	// some cleanup
})
```
