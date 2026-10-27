package sfx

import "github.com/gabrielluizsf/sge/sound"

// grassStep is a foot in grass: a swish of blades, and their spring back.
func grassStep(seed uint32) *sound.Clip {
	s := make([]float32, sound.Seconds(0.16))
	sound.MixInto(s, burst{seed, 0.12, 1800, 7000, sound.Envelope{Attack: 0.012, Decay: 0.09}}.make(), at(0, 1))
	sound.MixInto(s, burst{seed + 1, 0.08, 2500, 8000, sound.Envelope{Attack: 0.01, Decay: 0.06}}.make(), at(0.03, 0.5))
	return sound.NewClip(s, 0.5)
}

// soilStep is a foot on tilled soil: a soft thump, and clods crunching.
func soilStep(seed uint32) *sound.Clip {
	s := make([]float32, sound.Seconds(0.14))
	sound.MixInto(s, burst{seed, 0.1, 80, 700, sound.Envelope{Attack: 0.002, Decay: 0.06}}.make(), at(0, 1))
	r := rng(seed)
	for i := range 5 {
		grain := burst{seed + uint32(i) + 1, 0.014, 1500, 5000, sound.Envelope{Attack: 0.001, Decay: 0.011}}.make()
		sound.MixInto(s, grain, at(r.next()*0.05, 0.25+0.2*r.next()))
	}
	return sound.NewClip(s, 0.6)
}

// pathStep is a foot on the sandy path: grit shifting under it.
func pathStep(seed uint32) *sound.Clip {
	s := make([]float32, sound.Seconds(0.12))
	sound.MixInto(s, burst{seed, 0.08, 60, 400, sound.Envelope{Attack: 0.002, Decay: 0.05}}.make(), at(0, 0.6))
	r := rng(seed)
	for i := range 10 {
		grain := burst{seed + uint32(i) + 1, 0.008, 2500, 9000, sound.Envelope{Attack: 0.0005, Decay: 0.007}}.make()
		sound.MixInto(s, grain, at(r.next()*0.08, 0.2+0.4*r.next()))
	}
	return sound.NewClip(s, 0.45)
}
