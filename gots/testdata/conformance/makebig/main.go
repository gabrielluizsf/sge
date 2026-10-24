// make takes its length and capacity in any integer type, uint64 too,
// which gots keeps in a BigInt.
package main

import "fmt"

func main() {
	var n uint64 = 3
	b := make([]byte, n, n+2)
	m := make(map[string]int, n)
	c := make(chan int, n)
	var i8 int8 = 2
	s := make([]string, i8)
	fmt.Println(len(b), cap(b), len(m), cap(c), len(s))
}
