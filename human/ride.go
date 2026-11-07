package human

// Where one sits on a horse, in metres from the spot on the ground under its
// middle: the height of the hips on its back, and how far behind its middle
// they are, as measured on the farm's draft horse. Its body is near enough a
// barrel of BarrelRadius, its middle BarrelHeight up and BarrelLength long.
const (
	RideHeight   = 1.33
	RideBack     = 0.16
	BarrelHeight = 0.97
	BarrelRadius = 0.31
	BarrelLength = 1.6
)

// rideBob lifts a body from standing to its seat on the horse; rideAnkle is
// the foot level in its stirrup, the heel a little down.
const (
	rideBob   = RideHeight - (thighLen + shankLen)
	rideAnkle = -0.42
)

// ridePose is sitting astride, the thighs spread round the barrel and the
// shanks hanging by it, the hands forward on the reins, the body going with
// the horse's walk by sway, from -1 to 1.
func ridePose(sway float64) Pose {
	return Pose{Lean: 0.08 + 0.03*sway, Nod: 0.1, Hip: both(0.92), Knee: both(1.42), Ankle: both(rideAnkle), Abduct: both(0.34),
		Shoulder: both(0.1), Raise: both(0.1), Plane: both(0.25), Elbow: both(1.32 + 0.08*sway), Wrist: both(0.05),
		Bob: rideBob + 0.022*sway, Forward: -RideBack}
}

// Ride is sitting on a horse, moving with it.
var Ride = loop{"ride", 1.3333, 16}.of(func(t float64) Pose { return ridePose(wave(t, 2, 0)) })

// Riding are the motions of those who ride, besides Motions.
var Riding = []Motion{Mount, Ride, Dismount}

// RideVersion changes whenever the riding does, or the barrel it is baked
// round, so what was baked of it is baked again, and no other motion is.
const RideVersion = 4
