package sfx

import (
	"math"
	"testing"
)

// Every sound is made, is heard, peaks no louder than it should, and is
// a number all through; loops are seconds long.
func TestEverySoundIsMade(t *testing.T) {
	check := func(name string, s []float32) {
		var peak float64
		for _, x := range s {
			if x != x || math.IsInf(float64(x), 0) {
				t.Fatalf("%s has a sample of %v", name, x)
			}
			peak = math.Max(peak, math.Abs(float64(x)))
		}
		if peak < 0.1 || peak > 0.71 {
			t.Errorf("%s peaks at %v", name, peak)
		}
	}
	for name, cs := range Clips() {
		for _, c := range cs {
			check(name, c.Samples)
		}
	}
	for name, c := range Loops() {
		check(name, c.Samples)
		if len(c.Samples) < 4*c.Rate {
			t.Errorf("the %s loop is %d samples", name, len(c.Samples))
		}
	}
}

func BenchmarkMakeEverySound(b *testing.B) {
	for b.Loop() {
		Clips()
		Loops()
	}
}
