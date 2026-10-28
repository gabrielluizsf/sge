// Values in interfaces: gots hands a type switch case the boxed value
// itself when it only reads it, and builds boxes with their class named.
package main

import "fmt"

type V struct{ X, Y int }

func (v V) String() string { return fmt.Sprintf("V(%d,%d)", v.X, v.Y) }

type plain struct{ A int }

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1fC", float64(c)) }

type Pair[T any] struct{ L, R T }

func (p Pair[T]) String() string { return fmt.Sprintf("<%v|%v>", p.L, p.R) }

type onlyPtr struct{ N int }

func (o *onlyPtr) String() string { return "ptr" }

func describe(x any) string {
	switch v := x.(type) {
	case V:
		return "read " + v.String()
	case plain:
		v.A = 100 // a copy
		return fmt.Sprint("changed ", v)
	case Celsius:
		f := func() Celsius { return v + 1 } // kept by a closure
		v = 0
		return fmt.Sprint("closure ", f())
	case Pair[int]:
		p := &v // address taken
		p.L = 7
		return fmt.Sprint("addr ", v)
	}
	return "other"
}

func main() {
	vals := []any{V{1, 2}, plain{3}, Celsius(21.5), Pair[int]{1, 2}, Pair[string]{"a", "b"}, struct{ Z int }{4}, onlyPtr{5}, &onlyPtr{6}, 7, "s"}
	for _, x := range vals {
		fmt.Println(describe(x), x)
	}
	// Nothing describe did changed what the interfaces hold.
	fmt.Println(vals...)

	commaOk(V{1, 1})
	commaOk(Celsius(2))

	var s fmt.Stringer = V{8, 9}
	fmt.Println(s.String(), s)
	s = Celsius(-3)
	fmt.Println(s.String())
	s = Pair[float64]{1.5, 2.5}
	fmt.Println(s)
}

var store = map[string]V{}

// commaOk takes values out of x every way v, ok := x.(T) can be written.
func commaOk(x any) {
	if v, ok := x.(V); ok {
		fmt.Println("if-init", v, ok)
	}
	v, ok := x.(V)
	v.X = 50 // a copy
	fmt.Println("changed", v, ok, x)

	var r V
	r, ok = x.(V) // into a variable that exists
	r.Y = 60
	fmt.Println("existing", r, ok, x)

	var anyv any
	anyv, ok = x.(V) // into an interface
	fmt.Println("iface", anyv, ok)

	_, ok = x.(Celsius)
	fmt.Println("blank", ok)

	store["k"], ok = x.(V) // into a map
	if got, _ := x.(V); ok {
		got.X = 70
	}
	fmt.Println("map", store, x)

	var c, okc = x.(Celsius)
	fmt.Println("var", c, okc)

	p := &v
	w, okw := x.(V)
	p = &w // address taken: w must be a copy
	p.X = 80
	fmt.Println("pinned", w, okw, x)
}
