package main

import (
	"errors"
	"fmt"
)

func order() {
	for i := 0; i < 3; i++ {
		defer fmt.Print(i, " ")
	}
	fmt.Print("body ")
}

func named() (n int, err error) {
	defer func() {
		n *= 2
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	n = 21
	panic("oops")
}

func safeDiv(a, b int) (q int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("division failed")
		}
	}()
	return a / b, nil
}

func nilDeref() (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = "caught nil"
		}
	}()
	var p *struct{ X int }
	return fmt.Sprint(p.X)
}

func argsNow() {
	x := 1
	defer fmt.Println("deferred x =", x)
	x = 2
	fmt.Println("x =", x)
}

func repanic() (s string) {
	defer func() {
		s = fmt.Sprint(recover())
	}()
	defer func() {
		panic("second")
	}()
	panic("first")
}

type closer struct{ name string }

func (c *closer) Close() { fmt.Println("close", c.name) }

func resources() {
	a := &closer{"a"}
	defer a.Close()
	b := &closer{"b"}
	defer b.Close()
	fmt.Println("using")
}

func plain() int {
	defer func() {}()
	return 5
}

func main() {
	order()
	fmt.Println()
	fmt.Println(named())
	fmt.Println(safeDiv(7, 2))
	fmt.Println(safeDiv(1, 0))
	fmt.Println(nilDeref())
	argsNow()
	fmt.Println(repanic())
	resources()
	fmt.Println(plain())
	func() {
		defer func() {
			r := recover()
			err, ok := r.(error)
			fmt.Println("custom panic:", r, ok, err)
		}()
		panic(errors.New("as error"))
	}()
}
