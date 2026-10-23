package main

import "fmt"

type Point struct{ X, Y int }

type Rect struct {
	Min, Max Point
	Tags     [2]string
}

func (p Point) Moved(dx int) Point {
	p.X += dx
	return p
}

func (p *Point) Move(dx int) { p.X += dx }

func grow(r Rect) Rect {
	r.Max.X++
	r.Tags[0] = "grown"
	return r
}

func main() {
	a := Point{1, 2}
	b := a
	b.X = 10
	fmt.Println(a, b)

	c := a.Moved(5)
	fmt.Println(a, c)

	p := &a
	p.Move(3)
	fmt.Println(a, *p)
	a.Move(1)
	fmt.Println(a.X, p.X)

	r := Rect{Min: a, Max: Point{9, 9}, Tags: [2]string{"x", "y"}}
	g := grow(r)
	fmt.Println(r, g)

	arr := [3]int{1, 2, 3}
	arr2 := arr
	arr2[0] = 99
	fmt.Println(arr, arr2, arr == [3]int{1, 2, 3}, arr2 != arr)

	q := &r.Min
	q.X = -1
	fmt.Println(r.Min)

	r.Min = Point{7, 7}
	fmt.Println(*q, q == &r.Min)

	pts := []Point{{1, 1}, {2, 2}}
	first := pts[0]
	first.X = 100
	pts[1].Y = 50
	fmt.Println(pts, first)

	for _, pt := range pts {
		pt.X = 0
	}
	fmt.Println(pts)

	var zero Point
	fmt.Println(zero == Point{}, zero)

	anon := struct {
		Name string
		Age  int
	}{"ana", 30}
	other := anon
	other.Age++
	fmt.Println(anon, other)
	fmt.Printf("%+v %v\n", anon, &a)
}
