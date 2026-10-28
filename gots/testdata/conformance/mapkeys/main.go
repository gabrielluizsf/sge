// Map keys compare as Go compares them: pointers by identity, even
// behind an interface, structs and arrays by contents, a struct's pointer
// fields by identity.
package main

import (
	"fmt"
	"sync"
)

type conn struct{ id int }

type pair struct {
	a, b int
	p    *conn
}

type closer interface{ Close() }

func (*conn) Close() {}

func main() {
	x, y := &conn{1}, &conn{1} // same contents, two pointers
	byIface := map[closer]string{}
	byIface[x] = "x"
	byIface[y] = "y"
	delete(byIface, x)
	_, hasX := byIface[x]
	fmt.Println(len(byIface), hasX, byIface[y])

	byAny := map[any]int{x: 1, y: 2, conn{1}: 3}
	fmt.Println(len(byAny), byAny[x], byAny[y], byAny[conn{1}])

	byPair := map[pair]int{{1, 2, x}: 1}
	byPair[pair{1, 2, y}] = 2
	byPair[pair{1, 2, x}] += 10
	fmt.Println(len(byPair), byPair[pair{1, 2, x}], byPair[pair{1, 2, y}])

	byArr := map[[2]*conn]bool{{x, y}: true}
	fmt.Println(byArr[[2]*conn{x, y}], byArr[[2]*conn{y, x}])
	byC := map[complex128]int{1 + 2i: 7}
	fmt.Println(byC[1+2i], byC[2+1i])

	var m sync.Map
	m.Store(x, "x")
	_, ok := m.Load(y)
	v, _ := m.Load(x)
	fmt.Println(ok, v)
}
