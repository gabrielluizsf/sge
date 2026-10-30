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

// wave is a sine of a loop's t, as many times round as times, from phase.
func wave(t, times, phase float64) float64 { return math.Sin(2 * math.Pi * (t*times + phase)) }

// both is the same for each side.
func both(v float64) [2]float64 { return [2]float64{v, v} }

// Stand is standing at ease, breathing: the arms hang a little out and
// bent, the weight shifting.
var Stand = loop{"stand", 4, 8}.of(func(t float64) Pose {
	b := wave(t, 1, 0)
	return Pose{Lean: 0.02 + 0.01*b, Nod: 0.05, Hip: both(0.02), Knee: both(0.04), Raise: both(0.08 + 0.01*b),
		Elbow: both(0.2), Shoulder: both(0.03), Bob: 0.003 * b}
})

// Walk is walking, a stride a second: thighs swing and knees bend as the
// leg comes forward, arms swing against the legs, the body rising over
// each step.
var Walk = loop{"walk", 1, 16}.of(func(t float64) Pose {
	s, c := wave(t, 1, 0), wave(t, 1, 0.25)
	return Pose{Lean: 0.06, Nod: 0.04, Twist: 0.06 * s,
		Hip:      [2]float64{0.42 * s, -0.42 * s},
		Knee:     [2]float64{0.12 + 0.75*max(0, c), 0.12 + 0.75*max(0, -c)},
		Ankle:    [2]float64{-0.15 * s, 0.15 * s},
		Shoulder: [2]float64{-0.38 * s, 0.38 * s}, Raise: both(0.08),
		Elbow: [2]float64{0.3 + 0.2*max(0, -s), 0.3 + 0.2*max(0, s)},
		Bob:   -0.025 + 0.025*math.Abs(wave(t, 1, 0.25))}
})

// Chop is chopping wood with an axe held in both hands: raised over the
// right shoulder, brought down hard, and back up.
var Chop = loop{"chop", 1.2, 12}.of(func(t float64) Pose {
	up := 0.5 + 0.5*math.Cos(2*math.Pi*math.Pow(t, 0.7)) // 1 raised, 0 struck
	return Pose{Lean: 0.1 + 0.35*(1-up), Nod: 0.25, Twist: -0.25 * up,
		Hip: both(0.1), Knee: both(0.25),
		Shoulder: both(0.4 + 2.2*up), Raise: [2]float64{0.1, 0.25 + 0.2*up}, Elbow: both(0.5 + 0.9*up),
		Bob: -0.03}
})

// Hoe is working the ground with a hoe: reaching forward, pulling back.
var Hoe = loop{"hoe", 1.4, 12}.of(func(t float64) Pose {
	r := 0.5 + 0.5*wave(t, 1, 0)
	return Pose{Lean: 0.35 + 0.25*r, Nod: 0.3, Hip: [2]float64{0.35, -0.1}, Knee: [2]float64{0.35, 0.15},
		Shoulder: both(0.5 + 0.7*r), Raise: both(0.1), Elbow: both(1.1 - 0.7*r), Bob: -0.04}
})

// Swim is the crawl: lying on the water, each arm pulled through in turn,
// the legs kicking fast, the head turning to breathe.
var Swim = loop{"swim", 1.6, 16}.of(func(t float64) Pose {
	arm := func(phase float64) float64 { return math.Pi * (t + phase) * 2 } // full turns
	kick := wave(t, 3, 0)
	return Pose{Pitch: 1.45, Nod: -0.9, Twist: 0.3 * wave(t, 1, 0),
		Hip: [2]float64{0.25 * kick, -0.25 * kick}, Knee: [2]float64{0.2 + 0.2*max(0, kick), 0.2 + 0.2*max(0, -kick)},
		Ankle:    both(0.9),
		Shoulder: [2]float64{arm(0), arm(0.5)}, Raise: both(0.15),
		Elbow: [2]float64{0.4 + 0.5*max(0, wave(t, 1, 0)), 0.4 + 0.5*max(0, wave(t, 1, 0.5))},
		Bob:   -1.02}
})

// Tread is treading water: upright, head above it, the legs pedalling and
// the hands sculling out and in.
var Tread = loop{"tread", 1.5, 12}.of(func(t float64) Pose {
	p, q := wave(t, 1, 0), wave(t, 1, 0.25)
	return Pose{Lean: 0.08, Nod: -0.05,
		Hip: [2]float64{0.5 + 0.3*p, 0.5 - 0.3*p}, Knee: [2]float64{0.9 + 0.4*q, 0.9 - 0.4*q}, Ankle: both(0.3),
		Shoulder: both(0.6), Raise: both(0.7 + 0.25*p), Elbow: both(0.8), Bob: -1.35}
})

// Motions are every motion a person has.
var Motions = []Motion{Stand, Walk, Chop, Hoe, Swim, Tread}
