package try

import (
	"reflect"
	"fmt"
)

const pkgName = "github.com/mattn/go-try/try"

type RuntimeError string

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
		lhs := it.PkgPath() + "/" + it.Name()
		rhs := ct.PkgPath() + "/" + ct.Name()
		if rhs == "runtime/errorString" && lhs == pkgName+"/RuntimeError" {
			ev := reflect.ValueOf(RuntimeError(c.e.(fmt.Stringer).String()))
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
