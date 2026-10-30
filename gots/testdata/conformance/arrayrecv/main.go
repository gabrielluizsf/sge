// Command arrayrecv calls value-receiver methods of a named array type
// through pointers to it: to a variable, a struct's field and a slice's
// element, as render3d calls pose.Mat's methods.
package main

import "fmt"

// Mat is a named array with methods on its value.
type Mat [4]float64

// Sum is the sum of m's elements.
func (m Mat) Sum() float64 { return m[0] + m[1] + m[2] + m[3] }

// Scaled is m times k.
func (m Mat) Scaled(k float64) Mat {
	for i := range m {
		m[i] *= k
	}
	return m
}

// Reset is m made all zeros, by assigning the receiver a whole new array.
func (m Mat) Reset() Mat {
	m = Mat{}
	return m
}

// Set sets m's first element, through its pointer.
func (m *Mat) Set(v float64) { m[0] = v }

type holder struct{ M Mat }

func main() {
	v := Mat{1, 2, 3, 4}
	p := &v
	fmt.Println(p.Sum(), p.Scaled(2))
	h := &holder{M: Mat{5, 6, 7, 8}}
	q := &h.M
	fmt.Println(q.Sum(), q.Scaled(0.5).Sum())
	s := []Mat{{1, 1, 1, 1}, {2, 2, 2, 2}}
	r := &s[1]
	r.Set(10)
	fmt.Println(r.Sum(), s[1])
	// A value receiver is a copy: Scaled through the pointer, or on the
	// value, leaves the array as it was.
	p.Scaled(3)
	v.Scaled(3)
	fmt.Println(v, p.Sum(), p.Reset(), v)
	// Through an interface holding the pointer, and one holding the value.
	var sp, sv interface{ Sum() float64 } = p, v
	fmt.Println(sp.Sum(), sv.Sum())
	// A method value bound through the pointer copies the array then.
	f := p.Sum
	v[0] = 100
	fmt.Println(f(), p.Sum())
}
