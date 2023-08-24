package main

import (
	"log"
	"strconv"
)

type a struct {
	name string
}

type A interface {
	Name()
	Foo(string) string
}

func (a *a) Foo(s string) string {
	return a.name + s
}

type B interface {
	A
	Bar(int) string
}

type c struct {
	a
}

func (c c) Name() string {
	return c.name
}

func (c c) Bar(i int) string {
	return c.name + strconv.Itoa(i)
}

var C = c{
	a: a{name: "c"},
}

func test(v interface{}) string {
	vi := v.(B)
	return vi.Bar(1)
}

func main() {
	log.Println(test(C))
}
