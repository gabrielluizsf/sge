// Package bench has benchmarks for gots to run under deno bench.
package bench

// Sum adds up xs.
func Sum(xs []int) int {
	s := 0
	for _, x := range xs {
		s += x
	}
	return s
}
