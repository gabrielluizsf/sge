package main

import (
	"fmt"
	"unsafe"
)

func f64bits(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }
func f64from(b uint64) float64 { return *(*float64)(unsafe.Pointer(&b)) }
func f32bits(f float32) uint32 { return *(*uint32)(unsafe.Pointer(&f)) }
func f32from(b uint32) float32 { return *(*float32)(unsafe.Pointer(&b)) }
func i64bits(f float64) int64  { return *(*int64)(unsafe.Pointer(&f)) }

func main() {
	fmt.Println(f64bits(1.5), f64bits(-0.0), f64from(0x400921fb54442d18))
	fmt.Println(f32bits(1.5), f32from(0x3fc00000), i64bits(-2))
	s := []byte("bytes")
	fmt.Println(unsafe.String(unsafe.SliceData(s), len(s)), unsafe.Sizeof(int32(0)), unsafe.Sizeof(s))
	p := unsafe.StringData("hey")
	fmt.Println(string(unsafe.Slice(p, 3)))
}
