package human

import "math"

// The two sides of a person on a horse's near side: the leg that goes in the
// stirrup and the one that swings over, as they are seen on the screen.
const (
	swingLeg   = 0
	stirrupLeg = 1
)

// MountAside is how far from a horse's middle one stands to get on it, and where
// one is when down from it, in metres.
const MountAside = 0.7

// mountKey is a moment of getting on or down: the root's place as Aside, Bob
// and Forward, the body's Turn and Lean, each ankle's spot on the horse's
// frame and how far its toes hang from level, and the arms and head.
type mountKey struct {
	at         float64
	root       [3]float64
	turn, lean float64
	feet       [2][3]float64
	droop      [2]float64
	upper      Pose
}

// mountPath is a list of mountKeys as a motion: smooth curves through each
// key's root, turn, lean and feet, the arms mixed from key to key.
type mountPath []mountKey

// curveOf is the curve of a value of each key.
func (ks mountPath) curveOf(value func(mountKey) float64) curve {
	c := curve{at: make([]float64, len(ks)), keys: make([]float64, len(ks))}
	for i, k := range ks {
		c.at[i], c.keys[i] = k.at, value(k)
	}
	return c
}

// upperAt is the arms and head at t.
func (ks mountPath) upperAt(t float64) Pose {
	k := 0
	for k < len(ks)-2 && t >= ks[k+1].at {
		k++
	}
	return mix(ks[k].upper, ks[k+1].upper, ease((t-ks[k].at)/(ks[k+1].at-ks[k].at)))
}

// pose is the path at t.
func (ks mountPath) pose(t float64) Pose {
	value := func(v func(mountKey) float64) float64 { return ks.curveOf(v).of(t) }
	p := ks.upperAt(t)
	p.Aside = value(func(k mountKey) float64 { return k.root[0] })
	p.Bob = value(func(k mountKey) float64 { return k.root[1] })
	p.Forward = value(func(k mountKey) float64 { return k.root[2] })
	p.Turn = value(func(k mountKey) float64 { return k.turn })
	p.Lean = value(func(k mountKey) float64 { return k.lean })
	for s := range 2 {
		var spot [3]float64
		for i := range spot {
			spot[i] = value(func(k mountKey) float64 { return k.feet[s][i] })
		}
		legTo(&p, s, spot)
		p.Ankle[s] = p.Hip[s] - p.Knee[s] + value(func(k mountKey) float64 { return k.droop[s] })
	}
	return p
}

// motion is the path as a motion of a name, its ends exactly first and last.
func (ks mountPath) motion(name string, first, last Pose) Motion {
	m := loop{name, MountLength, mountFrames}.once(ks.pose)
	m.Poses[0], m.Poses[len(m.Poses)-1] = first, last
	return m
}

// standBeside is standing, off to the side of the horse.
func standBeside() Pose {
	p := Stand.Poses[0]
	p.Aside = -MountAside
	return p
}

// Mount is getting on a horse from its near side, as a rider is taught: turned
// to its side, a hand on its mane and one on the saddle, a foot in the stirrup;
// a push off the other foot, up straight on that leg; the other leg back over
// its croup, clear of it; sinking gently into the seat.
var Mount = mountKeys.motion("mount", standBeside(), ridePose(0))

// rootBehind is how far the root bone is behind a body's middle, which, turned
// a quarter, comes out to its side.
const rootBehind = -rootAhead

// standOut is standing beside the horse with one's back to it: as standing on
// that spot, turned a quarter from the horse's heading.
func standOut() Pose {
	p := Stand.Poses[0]
	p.Turn, p.Aside, p.Forward = -math.Pi/2, -MountAside+rootBehind, rootBehind
	return p
}

// Dismount is getting down, by the side one sits on: sitting sideways, then a
// jump out from the saddle.
var Dismount = lasting(dismountKeys.motion("dismount", ridePose(0), standOut()), DismountLength)

// lasting is m taking the seconds given.
func lasting(m Motion, seconds float64) Motion {
	m.Length = seconds
	return m
}
