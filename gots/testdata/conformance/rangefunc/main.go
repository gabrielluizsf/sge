// Range over functions whose bodies leave the loop: break, continue,
// return, and jumps to outer loops, through nested iterators too.
package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"sort"
)

func count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func pairs(n int) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i := 0; i < n; i++ {
			if !yield(i, fmt.Sprint("v", i)) {
				return
			}
		}
	}
}

func firstOver(limit int) (int, bool) {
	for x := range count(100) {
		if x*x > limit {
			return x, true
		}
	}
	return 0, false
}

func find(grid [][]int, want int) (r, c int) {
	for i := range count(len(grid)) {
		for j, v := range slices.All(grid[i]) {
			if v == want {
				return i, j
			}
		}
	}
	return -1, -1
}

func named() (out []int, err error) {
	defer func() { out = append(out, 99) }()
	for x := range count(10) {
		if x == 3 {
			continue
		}
		if x == 6 {
			return out, nil
		}
		out = append(out, x)
	}
	return nil, nil
}

// stubborn is an iterator that ignores yield saying stop.
func stubborn(yield func(int) bool) {
	yield(1)
	yield(2)
}

func breaksStubborn() (panicked bool) {
	defer func() { panicked = recover() != nil }()
	for x := range stubborn {
		if x == 1 {
			break
		}
		fmt.Println("never", x)
	}
	return false
}

func main() {
	fmt.Println("stubborn iterator panics:", breaksStubborn())
	fmt.Println(firstOver(50))
	fmt.Println(find([][]int{{1, 2}, {3, 4}}, 4))
	fmt.Println(named())

	sum := 0
	for x := range count(10) {
		if x%2 == 0 {
			continue
		}
		if x > 7 {
			break
		}
		sum += x
	}
	fmt.Println("sum", sum)

outer:
	for i := range 3 {
		for k, v := range pairs(5) {
			if k == 2 {
				continue outer
			}
			if i == 2 {
				break outer
			}
			fmt.Print(i, k, v, " ")
		}
	}
	fmt.Println()

	for i := range count(3) {
		switch i {
		case 1:
			fmt.Print("one ")
			break
		default:
			fmt.Print(i, " ")
		}
	}
	fmt.Println()

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := slices.Collect(maps.Keys(m))
	sort.Strings(keys)
	fmt.Println(keys, slices.Sorted(maps.Values(m)))
	c := maps.Clone(m)
	c["z"] = 9
	fmt.Println(len(m), len(c), maps.Clone(map[int]int(nil)) == nil)
	for k := range maps.Keys(m) {
		if k != "" {
			fmt.Println("some key")
			break
		}
	}
}
