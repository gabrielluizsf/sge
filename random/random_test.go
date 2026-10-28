package random

import "testing"

// The same seed plays the same numbers back.
func TestTheSameSeedIsTheSameSequence(t *testing.T) {
	a, b := New(7), New(7)
	for range 100 {
		if a.Float() != b.Float() {
			t.Fatal("the sequences parted")
		}
	}
}

// Numbers stay in their range and reach both ends of it.
func TestNumbersStayInTheirRange(t *testing.T) {
	s := New(1)
	seen := map[int]bool{}
	for range 1000 {
		f := s.Float()
		if f < 0 || f >= 1 {
			t.Fatalf("Float %v", f)
		}
		n := s.Between(2, 5)
		if n < 2 || n > 5 {
			t.Fatalf("Between %d", n)
		}
		seen[n] = true
	}
	if len(seen) != 4 {
		t.Errorf("Between(2, 5) gave only %v", seen)
	}
	if s.Intn(0) != 0 || s.Intn(-3) != 0 {
		t.Error("Intn of nothing is not 0")
	}
}
