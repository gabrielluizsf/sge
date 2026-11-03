package human

import "math"

// loop is a motion's name, how long it takes, and how many poses it is
// drawn from.
type loop struct {
	name   string
	length float64
	keys   int
}

// of is the motion of the poses f gives through a loop, t from 0 to 1.
func (l loop) of(f func(t float64) Pose) Motion {
	m := Motion{Name: l.name, Length: l.length}
	for k := range l.keys {
		m.Poses = append(m.Poses, f(float64(k)/float64(l.keys)))
	}
	return m
}

// once is the motion of the poses f gives from t 0, where it starts, to
// t 1, where it ends, which is not run on into the start.
func (l loop) once(f func(t float64) Pose) Motion {
	m := Motion{Name: l.name, Length: l.length, Once: true}
	for k := range l.keys {
		m.Poses = append(m.Poses, f(float64(k)/float64(l.keys-1)))
	}
	return m
}

// wave is a sine of a loop's t, as many times round as times, from phase.
func wave(t, times, phase float64) float64 { return math.Sin(2 * math.Pi * (t*times + phase)) }

// both is the same for each side.
func both(v float64) [2]float64 { return [2]float64{v, v} }

// Stand is standing at ease, breathing: the feet a hip's width apart, the
// arms hang a little out and bent, the weight shifting.
var Stand = loop{"stand", 4, 8}.of(func(t float64) Pose {
	b := wave(t, 1, 0)
	return Pose{Lean: 0.02 + 0.01*b, Nod: 0.05, Hip: both(0.02), Knee: both(0.04), Abduct: both(0.03), Raise: both(0.08 + 0.01*b),
		Elbow: both(0.2), Shoulder: both(0.03), Bob: 0.003 * b}
})

// Tread is treading water: upright, head above it, the legs pedalling and
// the hands sculling out and in.
var Tread = loop{"tread", 1.5, 12}.of(func(t float64) Pose {
	p, q := wave(t, 1, 0), wave(t, 1, 0.25)
	return Pose{Lean: 0.08, Nod: -0.05,
		Hip: [2]float64{0.5 + 0.3*p, 0.5 - 0.3*p}, Knee: [2]float64{0.9 + 0.4*q, 0.9 - 0.4*q}, Ankle: both(0.3), Abduct: both(0.1),
		Shoulder: both(0.6), Raise: both(0.7 + 0.25*p), Elbow: both(0.8), Bob: -1.35}
})

// Motions are every motion a person has: the loops, and the Dive into
// the water and the Rise out of it, which are run through once.
var Motions = []Motion{Stand, Walk, Chop, Hoe, Swim, Tread, Dive, Rise, Run}

// MotionsVersion changes whenever a motion does, so what was baked of the
// old ones is baked again.
const MotionsVersion = 3

// ChopVersion changes whenever the chop does, or the axe in the hands
// that swing it, so what was baked of it is baked again, and no other
// motion is.
const ChopVersion = 2

// ForgeVersion is the same for the smith's stroke at the anvil, and the
// hammer and tongs in their hands.
const ForgeVersion = 1

// HoeVersion is the same for the hoe's stroke and the hoe in the hands
// that swing it, WaterVersion for the pour and its can, PickVersion for picking a crop and RunVersion for running.
const (
	HoeVersion   = 1
	WaterVersion = 1
	PickVersion  = 1
	RunVersion   = 2
)
