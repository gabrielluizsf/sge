// Program appendcopy checks that append and copy copy struct values, and
// the arrays in them, as Go does: changing the copy leaves the original.
package main

import "fmt"

type quat [4]float64

type node struct {
	name string
	rot  [4]float64
	kids []int
}

func main() {
	orig := []node{{name: "a", rot: [4]float64{0, 0, 0, 1}}, {name: "b", rot: [4]float64{0, 0, 0, 1}}}
	c := append([]node{}, orig...)
	c[1].rot = [4]float64(quat{1, 0, 0, 0})
	c[0].rot[3] = 5
	c[0].name = "z"
	fmt.Println(orig[0].name, orig[0].rot, orig[1].rot)
	fmt.Println(c[0].name, c[0].rot, c[1].rot)

	d := make([]node, len(orig))
	copy(d, orig)
	d[0].rot[0] = 9
	fmt.Println(orig[0].rot, d[0].rot)

	q := quat{1, 2, 3, 4}
	r := [4]float64(q)
	r[0] = 7
	fmt.Println(q, r)

	arr := [2]node{orig[0], orig[1]}
	cp := arr
	cp[0].rot[1] = 3
	fmt.Println(arr[0].rot, cp[0].rot)
}
