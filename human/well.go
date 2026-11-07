package human

// The beats of filling the can at a well: standing at its edge, the left
// hand reaching up the rope and pulling it down to the chest, hand over
// hand twice, the body leaning back on each pull; then the pail is
// tipped by the left hand over the can the right hand holds out level
// under it, and the water runs in; and the can, full, is brought in to
// the hip and carried. The right arm holds the can all the while.
var wellBeats = []beat{
	{at: 0, ease: ease, pose: Pose{
		Lean: 0.10, Nod: 0.30, Hip: both(0.05), Knee: both(0.08), Abduct: both(0.1), Bob: -0.01,
		Shoulder: [2]float64{0.05, -0.50}, Raise: [2]float64{0.10, 0.24}, Plane: [2]float64{0, 0.42}, Elbow: [2]float64{0.30, 0.98}, Wrist: [2]float64{0, 0.90}}},
	{at: 0.15, ease: ease, pose: Pose{
		Lean: 0.28, Nod: 0.35, Hip: both(0.10), Knee: both(0.14), Abduct: both(0.1), Bob: -0.02,
		Shoulder: [2]float64{1.20, 0.60}, Raise: [2]float64{0.10, 0.12}, Plane: [2]float64{0.30, 0.40}, Elbow: [2]float64{0.40, 1.00}, Wrist: [2]float64{0, 0.10}}},
	{at: 0.28, ease: ease, pose: Pose{
		Lean: 0.10, Nod: 0.30, Hip: both(0.08), Knee: both(0.12), Abduct: both(0.1), Bob: -0.015,
		Shoulder: [2]float64{0.50, 0.60}, Raise: [2]float64{0.12, 0.12}, Plane: [2]float64{0.40, 0.40}, Elbow: [2]float64{1.70, 1.00}, Wrist: [2]float64{0, 0.10}}},
	{at: 0.40, ease: ease, pose: Pose{
		Lean: 0.30, Nod: 0.35, Hip: both(0.12), Knee: both(0.16), Abduct: both(0.1), Bob: -0.025,
		Shoulder: [2]float64{1.20, 0.60}, Raise: [2]float64{0.10, 0.12}, Plane: [2]float64{0.30, 0.40}, Elbow: [2]float64{0.50, 1.00}, Wrist: [2]float64{0, 0.10}}},
	{at: 0.52, ease: ease, pose: Pose{
		Lean: 0.08, Nod: 0.30, Hip: both(0.08), Knee: both(0.12), Abduct: both(0.1), Bob: -0.015,
		Shoulder: [2]float64{0.40, 0.60}, Raise: [2]float64{0.12, 0.12}, Plane: [2]float64{0.40, 0.40}, Elbow: [2]float64{1.80, 1.00}, Wrist: [2]float64{0, 0.10}}},
	{at: 0.62, ease: ease, pose: Pose{
		Lean: 0.30, Nod: 0.45, Hip: both(0.12), Knee: both(0.16), Abduct: both(0.1), Bob: -0.025,
		Shoulder: [2]float64{0.95, 0.70}, Raise: [2]float64{0.15, 0.12}, Plane: [2]float64{0.35, 0.40}, Elbow: [2]float64{1.00, 1.10}, Wrist: [2]float64{0.5, 0.15}}},
	{at: 0.80, ease: ease, pose: Pose{
		Lean: 0.32, Nod: 0.45, Hip: both(0.12), Knee: both(0.16), Abduct: both(0.1), Bob: -0.025,
		Shoulder: [2]float64{0.98, 0.72}, Raise: [2]float64{0.15, 0.12}, Plane: [2]float64{0.35, 0.40}, Elbow: [2]float64{1.00, 1.10}, Wrist: [2]float64{0.55, 0.15}}},
	{at: 0.90, ease: ease, pose: Pose{
		Lean: 0.15, Nod: 0.30, Hip: both(0.08), Knee: both(0.10), Abduct: both(0.1), Bob: -0.01,
		Shoulder: [2]float64{0.05, -0.50}, Raise: [2]float64{0.10, 0.24}, Plane: [2]float64{0, 0.42}, Elbow: [2]float64{0.30, 0.98}, Wrist: [2]float64{0, 0.90}}},
}

// Well is filling the can at a well: hauling the rope, tipping the pail
// into the can, and taking the can away.
var Well = loop{"well", 2.6, 40}.of(func(t float64) Pose { return poseAt(wellBeats, t) })

// WellVersion changes whenever the well's motion does, or the can in the
// hand that holds it, so what was baked of it is baked again, and no
// other motion is.
const WellVersion = 1
