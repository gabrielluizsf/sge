package human

// The beats of filling the can at a river's bank: crouching down on the
// heels, the left hand flat on the bank for balance, the right arm
// reaching out and down to sink the can into the water, holding it there
// while it fills, and lifting it out, heavy, as the legs straighten to
// stand with it at the side.
var dipBeats = []beat{
	{at: 0, ease: ease, pose: Pose{
		Lean: 0.12, Nod: 0.30, Hip: both(0.05), Knee: both(0.08), Abduct: both(0.1), Bob: -0.01,
		Shoulder: [2]float64{0.05, -0.50}, Raise: [2]float64{0.10, 0.24}, Plane: [2]float64{0, 0.42}, Elbow: [2]float64{0.30, 0.98}, Wrist: [2]float64{0, 0.90}}},
	{at: 0.22, ease: ease, pose: Pose{
		Lean: 0.50, Nod: 0.40, Hip: both(1.30), Knee: both(2.20), Abduct: both(0.18), Bob: -0.60, Forward: 0.25,
		Shoulder: [2]float64{0.75, 0.80}, Raise: [2]float64{0.10, 0.15}, Plane: [2]float64{0.30, 0.30}, Elbow: [2]float64{0.40, 0.50}, Wrist: [2]float64{0.10, 0.30}}},
	{at: 0.38, ease: ease, pose: Pose{
		Lean: 0.62, Nod: 0.40, Twist: -0.05, Hip: both(1.35), Knee: both(2.30), Abduct: both(0.18), Bob: -0.64, Forward: 0.30,
		Shoulder: [2]float64{0.85, 1.25}, Raise: [2]float64{0.10, 0.15}, Plane: [2]float64{0.30, 0.20}, Elbow: [2]float64{0.30, 0.25}, Wrist: [2]float64{0.10, 0.50}}},
	{at: 0.55, ease: ease, pose: Pose{
		Lean: 0.64, Nod: 0.40, Twist: -0.05, Hip: both(1.35), Knee: both(2.30), Abduct: both(0.18), Bob: -0.64, Forward: 0.30,
		Shoulder: [2]float64{0.85, 1.28}, Raise: [2]float64{0.10, 0.15}, Plane: [2]float64{0.30, 0.20}, Elbow: [2]float64{0.30, 0.25}, Wrist: [2]float64{0.10, 0.30}}},
	{at: 0.72, ease: ease, pose: Pose{
		Lean: 0.40, Nod: 0.35, Hip: both(0.85), Knee: both(1.50), Abduct: both(0.14), Bob: -0.38, Forward: 0.18,
		Shoulder: [2]float64{0.40, 0.60}, Raise: [2]float64{0.10, 0.15}, Plane: [2]float64{0.20, 0.20}, Elbow: [2]float64{0.30, 0.90}, Wrist: [2]float64{0.10, 0.10}}},
	{at: 0.88, ease: ease, pose: Pose{
		Lean: 0.15, Nod: 0.30, Hip: both(0.10), Knee: both(0.14), Abduct: both(0.1), Bob: -0.03, Forward: 0.04,
		Shoulder: [2]float64{0.05, -0.50}, Raise: [2]float64{0.10, 0.24}, Plane: [2]float64{0, 0.42}, Elbow: [2]float64{0.30, 0.98}, Wrist: [2]float64{0, 0.90}}},
}

// Dip is filling the can at a river: crouching at the bank, sinking the
// can in, and rising with it.
var Dip = loop{"dip", 2.4, 40}.of(func(t float64) Pose { return poseAt(dipBeats, t) })

// DipVersion changes whenever the dip does, or the can in the hand that
// holds it, so what was baked of it is baked again, and no other motion is.
const DipVersion = 1

// Drawing are the motions of those who fill a can, besides Motions.
var Drawing = []Motion{Well, Dip}
