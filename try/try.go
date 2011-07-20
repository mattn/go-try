package try

import (
	"reflect"
	"fmt"
	"runtime"
)

const pkgName = "github.com/mattn/go-try/try"

type StackInfo struct {
	PC uintptr
	File string
	Line int
}

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
}

func Try(f func()) (r *CatchOrFinally) {
	defer func() {
		if e := recover(); e != nil {
			r = &CatchOrFinally{e}
		}
	}()
	reflect.ValueOf(f).Call([]reflect.Value{})
	return nil
}

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
			i := 1
			for {
				if p, f, l, o := runtime.Caller(i); o {
					rte.StackTrace = append(rte.StackTrace, StackInfo{p, f, l})
					i++
				} else {
					break
				}
			}
			ev := reflect.ValueOf(rte)
			reflect.ValueOf(f).Call([]reflect.Value{ev})
			return nil
		} else if lhs == rhs {
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(c.e)})
			return nil
		}
	}
	return c
}

func (c *CatchOrFinally) Finally(f interface{}) {
	reflect.ValueOf(f).Call([]reflect.Value{})
}

func Throw(e interface{}) {
	panic(e)
}
