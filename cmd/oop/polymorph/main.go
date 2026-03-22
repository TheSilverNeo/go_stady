package main

import "fmt"

type C interface {
	test()
}

type A struct {
}

func (a *A) test() {
	fmt.Println("A")
}

type B struct {
}

func (b *B) test() {
	fmt.Println("B")
}

func polymorph(c C) {
	c.test()
}

func main() {
	a := A{}
	b := B{}

	some := []C{&a, &b}

	for _, cc := range some {
		polymorph(cc)
	}
}
