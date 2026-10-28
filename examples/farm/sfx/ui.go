package sfx

import "github.com/gabrielluizsf/sge/sound"

// knock is a light tap on wood, at pitch hertz.
func knock(s []float32, pitch float64, when float64) {
	sound.MixInto(s, note{sound.Tone{From: pitch, To: pitch * 0.88, Seconds: 0.06}, sound.Envelope{Attack: 0.001, Decay: 0.045}}.make(), at(when, 1))
	sound.MixInto(s, burst{uint32(pitch), 0.012, 1000, 4000, sound.Envelope{Attack: 0.0005, Decay: 0.01}}.make(), at(when, 0.3))
}

// knocks is a clip of a knock at each pitch, 70ms apart.
func knocks(pitches ...float64) *sound.Clip {
	s := make([]float32, sound.Seconds(0.07*float64(len(pitches))+0.08))
	for i, p := range pitches {
		knock(s, p, 0.07*float64(i))
	}
	return sound.NewClip(s, 0.5)
}

// bells is a clip of a small bell at each pitch, step seconds apart.
func bells(step float64, pitches ...float64) *sound.Clip {
	s := make([]float32, sound.Seconds(step*float64(len(pitches))+0.4))
	for i, p := range pitches {
		for h, g := range []float64{1, 0.3} {
			b := note{sound.Tone{From: p * float64(h+1), Seconds: 0.4}, sound.Envelope{Attack: 0.002, Decay: 0.35}}.make()
			sound.MixInto(s, b, at(step*float64(i), g))
		}
	}
	return sound.NewClip(s, 0.5)
}

// deny is a muffled double thud: that cannot be done.
func deny() *sound.Clip {
	s := make([]float32, sound.Seconds(0.3))
	for _, t := range []float64{0, 0.13} {
		sound.MixInto(s, note{sound.Tone{From: 300, To: 220, Seconds: 0.12}, sound.Envelope{Attack: 0.005, Decay: 0.1}}.make(), at(t, 1))
	}
	sound.Lowpass(s, 1200)
	return sound.NewClip(s, 0.5)
}

// plant is a seed pressed into soil.
func plant() *sound.Clip {
	s := make([]float32, sound.Seconds(0.12))
	sound.MixInto(s, burst{31, 0.1, 60, 600, sound.Envelope{Attack: 0.003, Decay: 0.08}}.make(), at(0, 1))
	sound.MixInto(s, note{sound.Tone{From: 180, To: 120, Seconds: 0.08}, sound.Envelope{Attack: 0.002, Decay: 0.07}}.make(), at(0, 0.6))
	return sound.NewClip(s, 0.55)
}

// pour is the can's water falling on a crop, wobbling as it pours.
func pour() *sound.Clip {
	s := burst{41, 0.45, 600, 2500, sound.Envelope{Attack: 0.04, Decay: 0.4}}.make()
	wob := sound.Tone{From: 18, Seconds: 0.45}.Make()
	for i := range s {
		s[i] *= 0.7 + 0.3*wob[i]
	}
	return sound.NewClip(s, 0.45)
}

// harvest is a crop popping out of the ground, and a sparkle.
func harvest() *sound.Clip {
	s := make([]float32, sound.Seconds(0.45))
	sound.MixInto(s, note{sound.Tone{From: 350, To: 800, Seconds: 0.06}, sound.Envelope{Attack: 0.001, Decay: 0.05}}.make(), at(0, 1))
	for i, p := range []float64{2637, 3136} {
		sound.MixInto(s, note{sound.Tone{From: p, Seconds: 0.3}, sound.Envelope{Attack: 0.002, Decay: 0.25}}.make(), at(0.05*float64(i+1), 0.45))
	}
	return sound.NewClip(s, 0.5)
}
