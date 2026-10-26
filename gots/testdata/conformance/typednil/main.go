// A nil pointer put in an interface is not a nil interface: it has a type,
// and its methods may be called with a nil receiver.
package main

import (
	"errors"
	"fmt"
)

type myErr struct{ msg string }

func (e *myErr) Error() string {
	if e == nil {
		return "nil *myErr"
	}
	return e.msg
}

// find returns a typed nil, the classic mistake.
func find(fail bool) error {
	var e *myErr
	if fail {
		e = &myErr{"failed"}
	}
	return e
}

type node struct{ next *node }

func (n *node) Len() int {
	if n == nil {
		return 0
	}
	return 1 + n.next.Len()
}

type lener interface{ Len() int }

func main() {
	err := find(false)
	fmt.Println(err == nil, err != nil, err.Error())
	fmt.Printf("%v %T\n", err, err)

	var e *myErr
	var any1, any2 any = e, (*myErr)(nil)
	fmt.Println(any1 == any2, any1 == nil, any1 == any(e))

	p, ok := err.(*myErr)
	fmt.Println(p == nil, ok)
	switch v := err.(type) {
	case nil:
		fmt.Println("nil interface")
	case *myErr:
		fmt.Println("typed nil:", v == nil)
	}
	var target *myErr
	fmt.Println(errors.As(err, &target), target == nil)

	var l lener = (*node)(nil)
	fmt.Println(l.Len(), lener(&node{&node{}}).Len())
	var nothing error
	fmt.Println(nothing == nil)
	m := map[any]int{any1: 1}
	fmt.Println(m[any2], len(m))
}
