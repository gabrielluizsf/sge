// Multiple results and parallel assignment, which gots writes without
// destructuring: V8 runs a destructuring through the iterator protocol.
package main

import (
	"errors"
	"fmt"
	"strconv"
)

func divmod(a, b int) (int, int) { return a / b, a % b }

func pair() (string, error) { return "x", errors.New("e") }

type P struct{ X int }

func point() (P, bool) { return P{3}, true }

func main() {
	q, r := divmod(17, 5)
	fmt.Println(q, r)

	_, r2 := divmod(9, 4)
	q3, _ := divmod(9, 4)
	fmt.Println(r2, q3)

	var s string
	var err error
	s, err = pair()
	fmt.Println(s, err)

	var anyv any
	anyv, err = pair() // into an interface
	fmt.Println(anyv, err)

	p, ok := point()
	p.X++
	p2, _ := point()
	fmt.Println(p, ok, p2)

	n, err := strconv.Atoi("42")
	fmt.Println(n, err)

	// Parallel assignment in a loop's post statement.
	for i, j := 0, 5; i < j; i, j = i+1, j-1 {
		fmt.Print(i, j, " ")
	}
	fmt.Println()
	a, b := 0, 1
	for k := 0; k < 10; a, b = b, a+b {
		k++
	}
	fmt.Println(a, b)

	x, y := 1, 2
	x, y = y, x
	fmt.Println(x, y)

	f := func() (int, int) { return 7, 8 }
	c, d := f()
	fmt.Println(c, d)
}
