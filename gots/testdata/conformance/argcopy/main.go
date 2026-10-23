// Struct and array arguments are copies in Go. gots leaves the copy out
// where nothing could tell; every case here could tell.
package main

import "fmt"

type V struct{ X, Y float64 }

// reads only: the copy may be left out.
func sum(o V) float64 { return o.X + o.Y }

func (v V) Add(o V) V { return V{v.X + o.X, v.Y + o.Y} }

// mutates its parameter.
func bump(o V) float64 {
	o.X += 100
	return o.X
}

// keeps it in a closure.
func later(o V) func() float64 { return func() float64 { return o.X } }

// returns it, and stores it.
func same(o V) V { return o }

var kept V

func keep(o V) { kept = o }

// takes its address.
func addr(o V) *V { return &o }

// swaps through a pointer while reading a copy of the same variable.
func (v *V) SetSwapped(o V) {
	v.X = o.Y
	v.Y = o.X
}

// arrays alias through a slice of them.
func sliceOf(a [3]int) []int { return a[:] }

func first(a [3]int) int { return a[0] }

func show(o V) { fmt.Println("deferred", o) }

func main() {
	p := V{1, 2}
	fmt.Println(sum(p), p.Add(p))

	fmt.Println(bump(p), p)

	f := later(p)
	p.X = 7
	fmt.Println(f(), p)

	q := same(p)
	q.X = 9
	keep(p)
	p.Y = 8
	fmt.Println(p, q, kept)

	r := addr(p)
	r.X = 11
	fmt.Println(p, *r)

	s := V{1, 2}
	s.SetSwapped(s)
	fmt.Println(s)

	a := [3]int{1, 2, 3}
	sl := sliceOf(a)
	sl[0] = 100
	fmt.Println(a, sl, first(a))

	d := V{3, 4}
	defer show(d)
	d.X = 30

	// A variable whose address is taken may change during a call.
	w := V{5, 6}
	pw := &w
	g := func() { pw.X = 50 }
	fmt.Println(sumThen(w, g), w)

	// Storing into a map copies, as any assignment does.
	mm := map[string]V{}
	y := V{1, 2}
	mm["y"] = y
	y.X = 5
	fmt.Println(mm, y)

	z := V{1, 1}
	b, sl2, ar, iface := stash(z)
	z.X = 99
	fmt.Println(z, b, sl2, ar, iface, list, byName)
}

type box struct{ v V }

var (
	list []V
	byName = map[string]V{}
)

// stores it in every way that keeps it.
func stash(o V) (box, []V, [1]V, any) {
	list = append(list, o)
	byName["o"] = o
	return box{o}, []V{o}, [1]V{o}, o
}

// sumThen reads o after f runs; o must be the value at the call.
func sumThen(o V, f func()) float64 {
	f()
	return o.X + o.Y
}
