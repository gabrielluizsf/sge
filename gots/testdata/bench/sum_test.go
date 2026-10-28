package bench

import "testing"

var sink int

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Fatalf("Sum = %d", got)
	}
}

// BenchmarkSumLoop times the loop only, through b.Loop.
func BenchmarkSumLoop(b *testing.B) {
	xs := make([]int, 1000)
	for b.Loop() {
		sink = Sum(xs)
	}
}

// BenchmarkSumN times the loop only, through ResetTimer and b.N.
func BenchmarkSumN(b *testing.B) {
	xs := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = Sum(xs)
	}
}

// BenchmarkSkipped is skipped, and must say so rather than fail.
func BenchmarkSkipped(b *testing.B) {
	b.Skip("not this one")
}
