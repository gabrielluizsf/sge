package sound

import "math"

// Envelope shapes a sound's loudness: rising over Attack seconds, then
// dying away, to a thousandth in Decay seconds more.
type Envelope struct {
	Attack, Decay float64
}

// Apply shapes s, at Rate, in place.
func (e Envelope) Apply(s []float32) {
	a := int(e.Attack * Rate)
	for i := range min(a, len(s)) {
		s[i] *= float32(i) / float32(a)
	}
	// Dying by the same factor every sample is dying exponentially.
	fall := math.Exp(-math.Log(1000) / math.Max(e.Decay*Rate, 1))
	g := 1.0
	for i := max(a, 0); i < len(s); i++ {
		s[i] *= float32(g)
		g *= fall
	}
}
