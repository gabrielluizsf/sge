// Values reached through an unexported field stay read-only whatever is
// done to them: converted, iterated over, appended to, copied into, or
// used as a map key or element.
package main

import (
	"fmt"
	"reflect"
)

type T struct{ N int }

type box struct {
	p     *T
	m     map[string]int
	s     []int
	Pub   map[string]int
	Slice []int
}

// panics reports whether f panics.
func panics(f func()) (did bool) {
	defer func() { did = recover() != nil }()
	f()
	return false
}

func main() {
	b := box{p: &T{}, m: map[string]int{"k": 1}, s: []int{1}, Pub: map[string]int{}, Slice: []int{0}}
	v := reflect.ValueOf(&b).Elem()
	p, m, s := v.Field(0), v.Field(1), v.Field(2)

	fmt.Println("convert", panics(func() { p.Convert(p.Type()).Elem().Field(0).SetInt(5) }))
	fmt.Println("iter", panics(func() {
		it := m.MapRange()
		it.Next()
		_ = it.Value().Interface()
	}))
	fmt.Println("append", panics(func() { reflect.Append(s, reflect.ValueOf(2)) }))
	fmt.Println("copy into", panics(func() { reflect.Copy(s, reflect.ValueOf([]int{9})) }))
	fmt.Println("copy from", panics(func() { reflect.Copy(v.Field(4), s) }))
	fmt.Println("map key", panics(func() { v.Field(3).SetMapIndex(m.MapKeys()[0], reflect.ValueOf(1)) }))
	fmt.Println("map elem", panics(func() { v.Field(3).SetMapIndex(reflect.ValueOf("x"), m.MapIndex(reflect.ValueOf("k"))) }))
	fmt.Println(b.p.N, b.s, b.Pub, b.Slice)
}
