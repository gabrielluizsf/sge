package sound

import "math"

// Tone is a sine Seconds long, gliding from From hertz to To; To 0 holds
// From.
type Tone struct {
	From, To, Seconds float64
}

// Make is the tone's samples at Rate.
func (t Tone) Make() []float32 {
	to := t.To
	if to == 0 {
		to = t.From
	}
	s := make([]float32, int(t.Seconds*Rate))
	phase := 0.0
	for i := range s {
		f := t.From + (to-t.From)*float64(i)/float64(len(s))
		phase += 2 * math.Pi * f / Rate
		s[i] = float32(math.Sin(phase))
	}
	return s
}
