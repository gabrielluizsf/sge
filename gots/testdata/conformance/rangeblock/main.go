// Ranges over functions whose bodies block: receiving from a channel,
// sleeping, reading a pipe, inside iterators of every kind.
package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"sort"
	"time"
)

func count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				return
			}
		}
	}
}

func pairs(yield func(string, int) bool) {
	for i, s := range []string{"a", "b", "c"} {
		if !yield(s, i) {
			return
		}
	}
}

func main() {
	ch := make(chan int, 10)
	go func() {
		for i := range 10 {
			ch <- i * i
		}
	}()
	sum := 0
	for i := range count(5) {
		sum += <-ch + i
	}
	fmt.Println("sum:", sum)

	for s, i := range pairs {
		time.Sleep(time.Millisecond)
		fmt.Println(s, i)
		if i == 1 {
			break
		}
	}

	m := map[string]int{"x": 1, "y": 2, "z": 3}
	keys := slices.Collect(maps.Keys(m))
	sort.Strings(keys)
	got := 0
	for k := range maps.Keys(m) {
		got += m[k] * <-ch
	}
	fmt.Println(keys, got > 0)

	found := func() string {
		for i := range count(100) {
			if v := <-ch; v > 40 {
				return fmt.Sprint("found ", v, " at ", i)
			}
		}
		return "none"
	}()
	fmt.Println(found)
	// A body that does not block stays synchronous.
	n := 0
	for i := range count(4) {
		n += i
	}
	fmt.Println(n)
}
