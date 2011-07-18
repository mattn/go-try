package main

import (
	. "github.com/mattn/go-try/try"
)

func main() {
	Try(func() {
		//panic("foo")
		//panic(1)
		v := 0
		println(1 / v)
	}).Catch(func(n int) {
		println("int exception:", n)
	}).Catch(func(s string) {
		println("string exception:", s)
	}).Catch(func(e RuntimeError) {
		println("runtime exception:", e)
	}).Finally(func() {
		println("finalize")
	})
}
