package main

import (
	"fmt"
	"sort"
)

type Key struct {
	A string
	B int
}

func main() {
	m := map[string]int{"one": 1, "two": 2}
	m["three"] = 3
	m["one"]++
	v, ok := m["two"]
	_, missing := m["zero"]
	fmt.Println(len(m), m["one"], v, ok, missing, m["nope"])
	delete(m, "two")
	fmt.Println(m)

	var nilMap map[string]bool
	fmt.Println(nilMap["x"], len(nilMap), nilMap == nil)

	sk := map[Key]string{}
	sk[Key{"a", 1}] = "first"
	sk[Key{"a", 1}] = "again"
	sk[Key{"b", 2}] = "second"
	fmt.Println(len(sk), sk[Key{"a", 1}], sk)

	counts := map[rune]int{}
	for _, r := range "banana" {
		counts[r]++
	}
	keys := make([]string, 0)
	for r, c := range counts {
		keys = append(keys, fmt.Sprintf("%c=%d", r, c))
	}
	sort.Strings(keys)
	fmt.Println(keys)

	type pt struct{ X, Y int }
	pm := map[string]pt{"a": {1, 2}}
	got := pm["a"]
	got.X = 50
	fmt.Println(pm["a"], got)

	ifm := map[any]int{1: 1, "1": 2, Key{"k", 0}: 3}
	fmt.Println(ifm[1], ifm["1"], ifm[Key{"k", 0}], ifm[2])

	nested := map[string][]int{}
	nested["x"] = append(nested["x"], 1, 2)
	fmt.Println(nested)
}

func init() {
	// Entries deleted and written again while ranging are not seen twice,
	// and entries deleted before they are reached are not seen at all.
	attrs := map[string]string{"xml:space": "preserve", "id": "g"}
	seen := 0
	for k, v := range attrs {
		seen++
		delete(attrs, k)
		attrs[k] = v + "!"
	}
	fmt.Println(seen, len(attrs), attrs["id"])
	m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
	visits := 0
	for k := range m {
		visits++
		for j := range m {
			if j != k {
				delete(m, j)
			}
		}
	}
	fmt.Println(visits, len(m))
}
