package sfx

import "github.com/gabrielluizsf/sge/sound"

// burst is a puff of noise, Seconds long, between Low and High hertz,
// shaped by its envelope.
type burst struct {
	seed      uint32
	seconds   float64
	low, high float64
	env       sound.Envelope
}

func (b burst) make() []float32 {
	s := sound.Noise(b.seed, sound.Seconds(b.seconds))
	sound.Highpass(s, b.low)
	sound.Lowpass(s, b.high)
	b.env.Apply(s)
	return s
}

// note is a tone shaped by an envelope.
type note struct {
	tone sound.Tone
	env  sound.Envelope
}

func (n note) make() []float32 {
	s := n.tone.Make()
	n.env.Apply(s)
	return s
}

// rng is a small deterministic random sequence, the same sounds every run.
type rng uint32

func (r *rng) next() float64 {
	*r = *r*1664525 + 1013904223
	return float64(*r>>8) / float64(1<<24)
}

// at is the place of sound at t seconds, gain loud.
func at(t, gain float64) sound.Place { return sound.Place{Sample: sound.Seconds(t), Gain: gain} }
