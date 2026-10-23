package main

import "fmt"

func counter() (func() int, func()) {
	n := 0
	return func() int { n++; return n }, func() { n = 0 }
}

func apply(xs []int, f func(int) int) []int {
	out := make([]int, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

func sum(nums ...int) int {
	t := 0
	for _, n := range nums {
		t += n
	}
	return t
}

func swap(a, b string) (string, string) { return b, a }

func fib() func() int {
	a, b := 0, 1
	return func() int {
		r := a
		a, b = b, a+b
		return r
	}
}

func inc(p *int) { *p++ }

type op func(float64, float64) float64

func main() {
	next, reset := counter()
	next()
	next()
	fmt.Println(next())
	reset()
	fmt.Println(next())

	sq := apply([]int{1, 2, 3}, func(x int) int { return x * x })
	fmt.Println(sq, sum(), sum(1, 2), sum(sq...))

	a, b := swap("x", "y")
	fmt.Println(a, b)
	fmt.Println(swap(swap("1", "2")))

	f := fib()
	for i := 0; i < 8; i++ {
		fmt.Print(f(), " ")
	}
	fmt.Println()

	x := 5
	inc(&x)
	inc(&x)
	px := &x
	*px += 10
	fmt.Println(x, *px == x)

	ops := map[string]op{
		"add": func(a, b float64) float64 { return a + b },
		"mul": func(a, b float64) float64 { return a * b },
	}
	fmt.Println(ops["add"](2, 3), ops["mul"](2, 3))

	var nilFn func()
	fmt.Println(nilFn == nil)

	result := func(n int) int {
		if n <= 1 {
			return 1
		}
		return n * 2
	}(5)
	fmt.Println(result)

	var fact func(int) int
	fact = func(n int) int {
		if n == 0 {
			return 1
		}
		return n * fact(n-1)
	}
	fmt.Println(fact(10))

	vals := []int{3, 1, 2}
	ptrs := []*int{}
	for i := range vals {
		ptrs = append(ptrs, &vals[i])
	}
	*ptrs[0] = 30
	fmt.Println(vals)

	fmt.Println(fmt.Sprintf("%d-%s", 7, "z"), len(fmt.Sprint(1, 2)))
}
