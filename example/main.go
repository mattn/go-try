package main

import (
	. "github.com/mattn/go-try/try"
	"fmt"
)

func main() {
	Try(func() {
		Try(func() {
			panic(1)
		}).Catch(func(n int) {
			println("int exception:", n)
		})

		Try(func() {
			panic("foo")
		}).Catch(func(s string) {
			println("string exception:", s)
		})

		v := 0
		println(1 / v)
	}).Catch(func(n int) {
		// not pass
		println("Catched int exception:", n)
	}).Catch(func(s string) {
		// not pass
		println("Catched string exception:", s)
	}).Catch(func(e RuntimeError) {
		fmt.Println("Catched runtime exception:", e)
		for _, st := range e.StackTrace {
			fmt.Printf("  %s:%d\n", st.File, st.Line)
		}
	}).Finally(func() {
		println("finalize")
	})
}
