package main

import "fmt"

func main() {
	var u8 uint8 = 250
	u8 += 10
	fmt.Println(u8)
	var i8 int8 = 127
	i8++
	fmt.Println(i8)
	var u16 uint16 = 0
	u16--
	fmt.Println(u16)
	var i32 int32 = 2147483647
	i32 += 1
	fmt.Println(i32)
	var u32 uint32 = 1 << 31
	u32 *= 2
	fmt.Println(u32)
	u32 = 0xdeadbeef
	fmt.Println(u32>>4, u32<<4, u32&0xff, u32|1, u32^0xffffffff, ^u32)
	h := uint32(2166136261)
	for _, c := range []byte("hash me") {
		h ^= uint32(c)
		h *= 16777619
	}
	fmt.Println(h)

	var u64 uint64 = 1<<64 - 1
	fmt.Println(u64, u64+1, u64>>60)
	x := uint64(0x9E3779B97F4A7C15)
	x ^= x >> 30
	x *= 0xBF58476D1CE4E5B9
	fmt.Println(x)

	n := 1 << 40
	fmt.Println(n, n>>35, n|7, n&(n-1), -7/2, -7%2, 7/-2)
	m128, u200, m32768 := int8(-128), uint8(200), int16(-32768)
	fmt.Println(m128/-1, u200+uint8(100), m32768-1, -m128)

	f := 3.99
	big33 := int64(1)<<33 + 5
	fmt.Println(int(f), int(-f), uint8(int(f)+254), int32(big33), uint16(big33))
	fmt.Println(float32(0.1), float64(float32(0.1)), 1/3.0)

	var r rune = 'é'
	fmt.Println(r, string(r), byte('A'+2))

	a, b := 17, 5
	fmt.Println(a/b, a%b, a&^b, a<<2, a>>1, ^a, -a)
	var shift uint = 70
	fmt.Println(1<<shift == 0, int32(1)<<shift, uint8(1)<<7)

	const big = 1 << 62
	fmt.Println(big/(1<<60), uint64(big)*2)
	fmt.Println(max(3, 9, 1), min(2.5, -1.0), max("a", "b"))
}
