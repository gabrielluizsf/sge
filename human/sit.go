package human

import "math"

// The seat of a chair: thighs level, and knees bent a right angle so the
// shanks hang straight down to the floor.
const (
	sitHip  = math.Pi / 2 * 0.95
	sitKnee = math.Pi / 2 * 0.95
	// sitLean is the back's lean, against a chair's back that slopes away.
	sitLean = -0.09
)

// Sit is sitting on a chair, the hands in the lap, breathing: the thighs
// level, the knees bent, the feet flat on the floor, the body as low as the
// legs leave it.
var Sit = loop{"sit", 4, 8}.of(func(t float64) Pose {
	b := wave(t, 1, 0)
	p := Pose{Lean: sitLean + 0.01*b, Nod: 0.02, Hip: both(sitHip), Knee: both(sitKnee), Ankle: both(-sitHip + sitKnee - 0.05),
		Abduct: both(0.06), Shoulder: both(0.15), Raise: both(0.02), Elbow: both(1.05), Wrist: both(0.1), Plane: both(0)}
	// The feet stay on the floor: the body is as high as the legs reach.
	p.Bob = legReach(p.Hip[0], p.Knee[0]) - (thighLen + shankLen) + 0.003*b
	return p
})

// Sitting are the motions of one who sits on a chair, besides Motions.
var Sitting = []Motion{Sit}

// SitVersion changes whenever the sitting does, so what was baked of it is
// baked again, and no other motion is.
const SitVersion = 2
