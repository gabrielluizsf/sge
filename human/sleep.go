package human

// Sleep is lying on the back in bed, asleep: the legs straight and a little
// apart, the arms by the sides, the head on the pillow, and the chest
// rising and falling slowly. Four poses to the breath are enough for it.
var Sleep = loop{"sleep", 4, 4}.of(func(t float64) Pose {
	b := wave(t, 1, 0)
	return Pose{Pitch: -sleepPitch, Lean: 0.02 + 0.02*b, Nod: 0.12,
		Hip: both(0.02), Knee: both(0.05), Ankle: both(0.2), Abduct: both(0.08),
		Shoulder: both(0.05), Raise: both(0.1), Elbow: both(0.25), Wrist: both(0.05),
		Bob: sleepBob + 0.004*b}
})

// The body lies flat: turned back by sleepPitch about the root, and let
// down by sleepBob so the back rests where the mattress is.
const (
	sleepPitch = 1.5708
	sleepBob   = -0.3
)

// Sleeping are the motions of one asleep in a bed, besides Motions.
var Sleeping = []Motion{Sleep}

// SleepVersion changes whenever the sleeping does, so what was baked of it
// is baked again, and no other motion is.
const SleepVersion = 1
