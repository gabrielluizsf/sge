package sfx

import "github.com/gabrielluizsf/sge/sound"

// bubbles adds n bubbles to s, each a quick rising tone, within spread
// seconds of its start.
func bubbles(s []float32, seed uint32, n int) {
	r := rng(seed)
	for range n {
		f := 450 + 350*r.next()
		b := note{sound.Tone{From: f, To: f * 2.2, Seconds: 0.05}, sound.Envelope{Attack: 0.003, Decay: 0.04}}.make()
		sound.MixInto(s, b, at(0.05+r.next()*0.25, 0.18))
	}
}

// stroke is an arm swept through the water, swimming.
func stroke(seed uint32) *sound.Clip {
	s := make([]float32, sound.Seconds(0.42))
	sound.MixInto(s, burst{seed, 0.38, 300, 2200, sound.Envelope{Attack: 0.04, Decay: 0.3}}.make(), at(0, 1))
	bubbles(s, seed+7, 3)
	return sound.NewClip(s, 0.45)
}

// splash is the farmer going into the water.
func splash(seed uint32) *sound.Clip {
	s := make([]float32, sound.Seconds(0.7))
	sound.MixInto(s, burst{seed, 0.6, 150, 3500, sound.Envelope{Attack: 0.005, Decay: 0.5}}.make(), at(0, 1))
	bubbles(s, seed+3, 6)
	return sound.NewClip(s, 0.7)
}
