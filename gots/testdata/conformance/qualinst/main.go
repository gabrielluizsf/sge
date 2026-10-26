// qualinst: generic functions of other packages called with their type
// arguments written out, one or several.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
)

func main() {
	_, err := os.Open("/no/such/file")
	pe, ok := errors.AsType[*fs.PathError](err)
	fmt.Println(ok, pe.Op)
	_, ok = errors.AsType[*os.LinkError](err)
	fmt.Println(ok)
	s := []int{3, 1, 2}
	fmt.Println(slices.Index[[]int, int](s, 2))
	slices.Sort[[]int](s)
	fmt.Println(s)
}

func init() {
	err := os.Rename("/no/such/old", "/no/such/new")
	le, ok := errors.AsType[*os.LinkError](err)
	fmt.Println(ok, le.Op, errors.Is(err, fs.ErrNotExist))
}
