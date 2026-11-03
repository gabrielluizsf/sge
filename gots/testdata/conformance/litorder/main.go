// Program litorder checks that a keyed composite literal evaluates its
// elements in the order they are written, as Go does, not in the order
// the struct declares its fields.
package main

import "fmt"

type thing struct {
	a, b, c int
	inner   struct{ x, y int }
}

var n int

func next() int {
	n++
	return n
}

func main() {
	t := thing{c: next(), a: next(), b: next()}
	fmt.Println(t.a, t.b, t.c)
	u := &thing{b: next(), c: next()}
	fmt.Println(u.a, u.b, u.c)
	arr := [3]int{2: next(), 0: next(), 1: next()}
	fmt.Println(arr)
	m := map[string]int{"z": next(), "a": next()}
	fmt.Println(m["z"], m["a"])
}
