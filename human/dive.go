package human

import "math"

// Dive is going into the water from standing, once, in 0.8 seconds: a
// crouch with the arms swung back, a spring up and forward, the arms
// thrown overhead, the body tipping over to enter head first, and the
// swimmer's own pose as it ends. The body travels about 1.5 m ahead.
var Dive = loop{"dive", 0.8, 8}.once(func(t float64) Pose {
	stand, swim := Stand.Poses[0], Swim.Poses[0]
	key := []Pose{
		stand,
		{Lean: 0.5, Nod: 0.1, Hip: both(0.9), Knee: both(1.2), Ankle: both(-0.3), Abduct: both(0.1), Shoulder: both(-0.8), Raise: both(0.15), Elbow: both(0.3), Bob: -0.2},
		{Pitch: 0.25, Lean: 0.15, Hip: both(0.1), Knee: both(0.2), Ankle: both(0.85), Abduct: both(0.05), Raise: both(2.0), Plane: both(0.5), Elbow: both(0.2), Bob: 0.12, Forward: 0.3},
		{Pitch: 0.7, Lean: 0.1, Nod: -0.3, Hip: both(0.1), Knee: both(0.1), Ankle: both(0.85), Abduct: both(0.03), Raise: both(2.8), Plane: both(0.1), Elbow: both(0.1), Bob: 0.35, Forward: 0.8},
		{Pitch: 1.0, Lean: 0.1, Nod: -0.45, Hip: both(0.1), Knee: both(0.1), Ankle: both(0.85), Abduct: both(0.03), Raise: both(2.8), Plane: both(0.1), Elbow: both(0.1), Bob: 0.25, Forward: 1.1},
		{Pitch: 1.3, Nod: -0.7, Hip: both(0.1), Knee: both(0.15), Ankle: both(0.85), Abduct: both(0.05), Raise: both(2.75), Plane: both(0.15), Elbow: both(0.2), Bob: -0.3, Forward: 1.4},
	}
	key = append(key, mix(key[5], swim, 0.5), swim)
	key[6].Bob = -0.8
	key[6].Forward = 1.0
	key[7].Forward = 0 // the game sets the swimmer down where the leap ended
	x := t * float64(len(key)-1)
	i := min(int(x), len(key)-2)
	return mix(key[i], key[i+1], x-float64(i))
})

// Rise is coming up out of the water to stand, once, in a second: the
// swimmer's pose first, then the legs drawn under to find the bottom, the
// body coming upright, the arms down, and the standing pose as it ends.
var Rise = loop{"rise", 1, 8}.once(func(t float64) Pose {
	f := ease(t)
	p := mix(Swim.Poses[0], Stand.Poses[0], f)
	dip := math.Sin(math.Pi * f)
	p.Hip = both(p.Hip[0] + 0.7*dip)
	p.Knee = both(p.Knee[0] + 1.0*dip)
	p.Lean += 0.3 * dip
	p.Ankle = both(p.Ankle[0] * (1 - dip))
	return p
})
