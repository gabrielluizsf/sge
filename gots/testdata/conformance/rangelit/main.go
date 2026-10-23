// Ranging over array literals, which gots writes as plain JavaScript
// arrays: every element kind has to read back as the literal held it.
package main

import "fmt"

type pt struct{ x, y int }

func main() {
	fy := 0.25
	sum := 0.0
	for j, w := range [2]float64{1 - fy, fy} {
		sum += float64(j+1) * w
	}
	fmt.Println("float64", sum)

	var f32 float32 = 0.1
	for i, v := range [2]float32{f32 * 3, f32} {
		fmt.Println("float32", i, v)
	}

	a, b := int8(100), int8(100)
	for i, v := range [3]int8{a + b, -a - b - b, 7} {
		fmt.Println("int8", i, v)
	}

	x := uint8(250)
	for _, v := range [2]uint8{x + 10, x} {
		fmt.Println("uint8", v)
	}

	big := uint64(1) << 63
	for _, v := range [2]uint64{big, big + big - 1} {
		fmt.Println("uint64", v)
	}

	for i, p := range [2]pt{{1, 2}, {3, 4}} {
		p.x *= 10 // a copy: the literal is not changed
		fmt.Println("struct", i, p)
	}

	for i, s := range [3]string{"a", 2: "c"} {
		fmt.Printf("sparse %d %q\n", i, s)
	}

	for i := range [4]int{} {
		fmt.Println("empty", i)
	}

	for i, v := range &[2]int{5, 6} {
		fmt.Println("pointer", i, v)
	}
}
