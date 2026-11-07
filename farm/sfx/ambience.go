package sfx

import (
	"math"

	"github.com/gabrielluizsf/sge/sound"
)

// swell multiplies s by a slow rise and fall of loudness, from low to 1,
// of whole waves over its length, so a loop of it does not jump. The
// swell is slow, so it is worked out every 64 samples and drawn straight
// between.
func swell(s []float32, waves []float64, low float64) {
	n := float64(len(s))
	gain := func(i int) float64 {
		t := float64(i) / n
		v := 0.0
		for k, w := range waves {
			v += math.Sin(2*math.Pi*w*t+float64(k)*1.7) / float64(len(waves))
		}
		return low + (1-low)*(0.5+0.5*v)*(0.5+0.5*v)
	}
	const step = 64
	for i0 := 0; i0 < len(s); i0 += step {
		g0, g1 := gain(i0), gain(i0+step)
		for i := i0; i < min(i0+step, len(s)); i++ {
			s[i] *= float32(g0 + (g1-g0)*float64(i-i0)/step)
		}
	}
}

// Rustle is leaves in the wind, to loop near trees: a breath of noise
// rising and falling in gusts, with leaves ticking in it.
func rustle() *sound.Clip {
	s := sound.Noise(11, sound.Seconds(6.3))
	sound.Highpass(s, 900)
	sound.Lowpass(s, 5000)
	swell(s, []float64{1, 2, 5}, 0.25)
	r := rng(12)
	for i := range 140 {
		leaf := burst{uint32(100 + i), 0.01, 3000, 9000, sound.Envelope{Attack: 0.001, Decay: 0.008}}.make()
		sound.MixInto(s, leaf, at(r.next()*6.2, 0.6*r.next()))
	}
	return sound.NewClip(sound.Seamless(s, sound.Seconds(0.3)), 0.5)
}

// lap is water lapping at a shore, to loop near water.
func lap() *sound.Clip {
	s := sound.Noise(21, sound.Seconds(5.3))
	sound.Highpass(s, 60)
	sound.Lowpass(s, 450)
	swell(s, []float64{3, 4, 7}, 0.1)
	r := rng(22)
	for range 9 {
		f := 700 + 500*r.next()
		plip := note{sound.Tone{From: f, To: f * 1.6, Seconds: 0.04}, sound.Envelope{Attack: 0.002, Decay: 0.035}}.make()
		sound.MixInto(s, plip, at(r.next()*5.1, 0.25))
	}
	return sound.NewClip(sound.Seamless(s, sound.Seconds(0.3)), 0.5)
}

// chirp is a small bird in a tree: two to four quick notes.
func chirp(seed uint32) *sound.Clip {
	r := rng(seed)
	n := 2 + int(r.next()*3)
	s := make([]float32, sound.Seconds(0.12*float64(n)+0.1))
	base := 2600 + 900*r.next()
	for i := range n {
		from, to := base*(0.95+0.1*r.next()), base*(1.15+0.2*r.next())
		if i%2 == 1 {
			from, to = to, from
		}
		tweet := note{sound.Tone{From: from, To: to, Seconds: 0.06 + 0.04*r.next()}, sound.Envelope{Attack: 0.005, Decay: 0.06}}.make()
		sound.MixInto(s, tweet, at(0.12*float64(i), 1))
	}
	return sound.NewClip(s, 0.35)
}
