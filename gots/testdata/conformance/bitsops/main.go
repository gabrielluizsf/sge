// bitsops: math/bits's 32-bit functions, which are the runtime's under
// gots (rt/bits.ts), on edge values. uint's, 32 bits there too, are the
// same functions, but 64 bits in the Go this is checked against.
package main

import (
	"fmt"
	"math/bits"
)

var words = []uint32{0, 1, 2, 3, 0x7fffffff, 0x80000000, 0x80000001, 0xfffffffe, 0xffffffff, 0xdeadbeef, 0x00010000, 0x0000ffff, 0x12345678}

func main() {
	for _, x := range words {
		fmt.Println(x, bits.LeadingZeros32(x), bits.Len32(x), bits.TrailingZeros32(x), bits.OnesCount32(x),
			bits.RotateLeft32(x, 7), bits.RotateLeft32(x, -13), bits.RotateLeft32(x, 32), bits.RotateLeft32(x, 0))
		for _, y := range words {
			for _, c := range []uint32{0, 1} {
				s, co := bits.Add32(x, y, c)
				d, bo := bits.Sub32(x, y, c)
				fmt.Println(s, co, d, bo)
			}
			hi, lo := bits.Mul32(x, y)
			fmt.Println(hi, lo)
		}
	}
}
