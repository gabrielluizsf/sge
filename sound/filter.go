package sound

import "math"

// Lowpass softens s in place, at Rate, above cutoff hertz: one pole, as
// muffled as a sound through a wall.
func Lowpass(s []float32, cutoff float64) {
	a := float32(1 - math.Exp(-2*math.Pi*cutoff/Rate))
	var y float32
	for i, x := range s {
		y += (x - y) * a
		s[i] = y
	}
}

// Highpass thins s in place, at Rate, below cutoff hertz: what Lowpass
// would keep, taken away.
func Highpass(s []float32, cutoff float64) {
	a := float32(1 - math.Exp(-2*math.Pi*cutoff/Rate))
	var y float32
	for i, x := range s {
		y += (x - y) * a
		s[i] = x - y
	}
}
