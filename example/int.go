package main

import (
	. "github.com/mattn/go-try/try"
)

func main() {
	Try(func() {
		panic(1)
	}).Catch(func(n int) {
		println("int exception:", n)
	})
}
