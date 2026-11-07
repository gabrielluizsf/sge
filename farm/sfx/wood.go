package sfx

import "github.com/gabrielluizsf/sge/sound"

// chop is an axe biting into a trunk: a hollow thunk and a splintering
// crack.
func chop(seed uint32) *sound.Clip {
	s := make([]float32, sound.Seconds(0.25))
	sound.MixInto(s, note{sound.Tone{From: 190 + float64(seed%5)*12, To: 120, Seconds: 0.14}, sound.Envelope{Attack: 0.002, Decay: 0.12}}.make(), at(0, 1))
	sound.MixInto(s, burst{seed, 0.08, 1500, 6000, sound.Envelope{Attack: 0.001, Decay: 0.06}}.make(), at(0.004, 0.6))
	return sound.NewClip(s, 0.5)
}

// fell is a tree coming down: a creak, a rush of leaves and a thump.
func fell() *sound.Clip {
	s := make([]float32, sound.Seconds(1.1))
	sound.MixInto(s, note{sound.Tone{From: 260, To: 140, Seconds: 0.4}, sound.Envelope{Attack: 0.05, Decay: 0.3}}.make(), at(0, 0.5))
	sound.MixInto(s, burst{9, 0.5, 2000, 7000, sound.Envelope{Attack: 0.2, Decay: 0.3}}.make(), at(0.25, 0.5))
	sound.MixInto(s, note{sound.Tone{From: 90, To: 50, Seconds: 0.3}, sound.Envelope{Attack: 0.003, Decay: 0.28}}.make(), at(0.7, 1))
	sound.Lowpass(s, 5000)
	return sound.NewClip(s, 0.5)
}
