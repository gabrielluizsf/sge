package main

import "fmt"

func classify(n int) string {
	switch {
	case n < 0:
		return "neg"
	case n == 0:
		return "zero"
	case n%2 == 0:
		return "even"
	default:
		return "odd"
	}
}

func main() {
	for i := -1; i < 4; i++ {
		fmt.Print(classify(i), " ")
	}
	fmt.Println()

outer:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j == 3 {
				continue outer
			}
			if i == 3 {
				break outer
			}
			fmt.Print(i, j, ";")
		}
	}
	fmt.Println()

	for i := 0; i < 6; i++ {
		switch i {
		case 0, 1:
			fmt.Print("low ")
		case 2:
			fmt.Print("two ")
			fallthrough
		case 3:
			fmt.Print("three-ish ")
		case 4:
			if i == 4 {
				break
			}
			fmt.Print("never ")
		default:
			fmt.Print("high ")
		}
	}
	fmt.Println()

	n := 0
	for n < 5 {
		n += 2
	}
	fmt.Println(n)

	sum := 0
	for i := range 5 {
		sum += i
	}
	fmt.Println(sum)

	var funcs []func() int
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() int { return i * 10 })
	}
	for _, f := range funcs {
		fmt.Print(f(), " ")
	}
	fmt.Println()

	for i, j := 0, 10; i < j; i, j = i+1, j-2 {
		fmt.Print(i, ":", j, " ")
	}
	fmt.Println()

	x := 1
	{
		x := x + 1
		fmt.Println("inner", x)
	}
	fmt.Println("outer", x)

	if v := classify(7); v == "odd" {
		fmt.Println("seven is", v)
	} else if v == "even" {
		fmt.Println("no")
	} else {
		fmt.Println("never")
	}

	count := 0
	for {
		count++
		if count > 3 {
			break
		}
	}
	fmt.Println(count)

	switch x := 5; {
	case x > 3:
		fmt.Println("big")
	}
}
