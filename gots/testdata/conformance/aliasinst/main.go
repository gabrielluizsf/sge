// An alias of an instance of a generic struct, whose class is written
// after the type declarations: the alias must not reach it before.
package main

import "fmt"

type Pair[K comparable, V any] struct {
	Key K
	Val V
}

type Entry = Pair[string, int]

func main() {
	e := Entry{"a", 1}
	var p Pair[string, int] = e
	fmt.Println(e, p, fmt.Sprintf("%T", e))
}
