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

// rideBob lifts a body from standing to its seat on the horse.
const rideBob = RideHeight - (thighLen + shankLen)

// ridePose is sitting astride, the thighs spread round the barrel and the
// shanks hanging by it, the hands forward on the reins, the body going with
// the horse's walk by sway, from -1 to 1.
func ridePose(sway float64) Pose {
	return Pose{Lean: 0.08 + 0.03*sway, Nod: 0.1, Hip: both(0.92), Knee: both(1.42), Ankle: both(0.28), Abduct: both(0.34),
		Shoulder: both(0.1), Raise: both(0.1), Plane: both(0.25), Elbow: both(1.32 + 0.08*sway), Wrist: both(0.05),
		Bob: rideBob + 0.022*sway, Forward: -RideBack}
}

// Ride is sitting on a horse, moving with it.
var Ride = loop{"ride", 1.3333, 16}.of(func(t float64) Pose { return ridePose(wave(t, 2, 0)) })

// mountKeys are the moments of getting on, from beside the horse: a foot
// raised to the stirrup, hands up on its back; pulled up on that leg, the
// other swinging behind, knee bent; that knee coming over its back; settling
// down; and in the seat.
var mountKeys = []Pose{
	Stand.Poses[0],
	{Lean: 0.16, Nod: 0.05, Hip: [2]float64{1.2, 0.06}, Knee: [2]float64{1.75, 0.12}, Ankle: [2]float64{0.3, 0}, Abduct: [2]float64{0.2, 0.06},
		Shoulder: both(1.2), Raise: both(0.1), Elbow: both(0.75), Bob: 0.02},
	{Lean: 0.38, Nod: 0.1, Hip: [2]float64{1.0, -0.25}, Knee: [2]float64{1.5, 0.6}, Ankle: [2]float64{0.25, 0.3}, Abduct: [2]float64{0.25, 0.2},
		Shoulder: both(0.95), Raise: both(0.12), Elbow: both(0.9), Bob: 0.22, Forward: -RideBack * 0.4},
	{Lean: 0.36, Nod: 0.1, Hip: [2]float64{1.0, 0.2}, Knee: [2]float64{1.45, 1.1}, Ankle: both(0.25), Abduct: [2]float64{0.32, 0.75},
		Shoulder: both(0.6), Raise: both(0.12), Elbow: both(1.05), Bob: rideBob - 0.03, Forward: -RideBack * 0.8},
	{Lean: 0.2, Nod: 0.1, Hip: [2]float64{0.95, 0.8}, Knee: [2]float64{1.42, 1.35}, Ankle: both(0.28), Abduct: [2]float64{0.34, 0.45},
		Shoulder: both(0.25), Raise: both(0.1), Elbow: both(1.25), Bob: rideBob + 0.01, Forward: -RideBack},
	ridePose(0),
}

// Mount is getting on a horse, from beside it; run from its end to its start,
// it is getting down.
var Mount = loop{"mount", 1.2, 16}.once(func(t float64) Pose {
	u := t * float64(len(mountKeys)-1)
	k := min(int(u), len(mountKeys)-2)
	return mix(mountKeys[k], mountKeys[k+1], ease(u-float64(k)))
})

// Riding are the motions of those who ride, besides Motions.
var Riding = []Motion{Mount, Ride}

// RideVersion changes whenever the riding does, or the barrel it is baked
// round, so what was baked of it is baked again, and no other motion is.
const RideVersion = 1
