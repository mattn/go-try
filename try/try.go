package try

import (
	"reflect"
	"fmt"
	"runtime"
)

//StackInfo store code informations when catched exception.
type StackInfo struct {
	PC uintptr
	File string
	Line int
}

//RuntimeError is wrapper of runtime.errorString and stacktrace.
type RuntimeError struct {
	fmt.Stringer
	Message string
	StackTrace []StackInfo
}

func (rte RuntimeError) String() string {
	return rte.Message
}

type CatchOrFinally struct {
	e interface{}
	StackTrace []StackInfo
}

type OrThrowable struct {
	e interface{}
}

//Try call the function. And return interface that can call Catch or Finally.
func Try(f func()) (r *CatchOrFinally) {
	defer func() {
		if e := recover(); e != nil {
			r = &CatchOrFinally{}
			r.e = e
			i := 1
			for {
				if p, f, l, o := runtime.Caller(i); o {
					r.StackTrace = append(r.StackTrace, StackInfo{p, f, l})
					i++
				} else {
					break
				}
			}
		}
	}()
	reflect.ValueOf(f).Call([]reflect.Value{})
	return nil
}

//Catch call the exception handler. And return interface CatchOrFinally that
//can call Catch or Finally.
func (c *CatchOrFinally) Catch(f interface{}) (r *CatchOrFinally) {
	if c == nil || c.e == nil {
		return nil
	}
	rf := reflect.ValueOf(f)
	ft := rf.Type()
	if ft.NumIn() > 0 {
		it := ft.In(0)
		ct := reflect.TypeOf(c.e)
		lhs := it.String()
		rhs := ct.String()
		if rhs == "runtime.errorString" && lhs == "try.RuntimeError" {
			var rte RuntimeError
			rte.Message = c.e.(fmt.Stringer).String()
			rte.StackTrace = c.StackTrace
			ev := reflect.ValueOf(rte)
			reflect.ValueOf(f).Call([]reflect.Value{ev})
			return nil
		} else if lhs == rhs {
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(c.e)})
			return nil
		}
		println(lhs, rhs)
	}
	return c
}

//Finally always be called if defined.
func (c *CatchOrFinally) Finally(f interface{}) (r *OrThrowable) {
	reflect.ValueOf(f).Call([]reflect.Value{})
	return &OrThrowable{c.e}
}

//OrThrow throw error then never catch block entered.
func (c *CatchOrFinally) OrThrow() {
	if c != nil && c.e != nil {
		Throw(c.e)
	}
}

//OrThrow throw error then never catch block entered.
func (c *OrThrowable) OrThrow() {
	if c != nil && c.e != nil {
		Throw(c.e)
	}
}

//Throw is wrapper of panic().
func Throw(e interface{}) {
	panic(e)
}
