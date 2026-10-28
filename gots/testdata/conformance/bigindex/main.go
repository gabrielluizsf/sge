// Indexes and slice bounds of every integer type, uint64 among them, as
// sha256 slices its padding.
package main

import "fmt"

func main() {
	buf := make([]byte, 16)
	arr := [8]int{}
	var i64 uint64 = 3
	var i32 uint32 = 5
	var i8 int8 = 1
	buf[i64] = 7
	arr[i64] = 9
	buf[i32]++
	p := &buf[i64]
	*p += 1
	q := &arr[i8]
	*q = 4
	fmt.Println(buf[i64], buf[i32], arr[i64], arr[i8], buf[:i64+2], buf[i8:i64], len(buf[i64:]), arr[i8:i64])
	s := "uint64 index"
	fmt.Println(string(s[i64]), s[i64:i64+4], s[:i32])
	tmp := make([]byte, 72)
	var t uint64 = 56
	fmt.Println(len(tmp[0 : t+8]))
}
