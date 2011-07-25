package main

import (
	. "github.com/mattn/go-try/try"
)

func main() {
	Try(func() {
		panic("foo")
	}).Catch(func(s string) {
		println("string exception:", s)
	}).OrThrow()
}
