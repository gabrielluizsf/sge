package human

// The beats of picking a crop from the ground a stride ahead: standing and
// looking at the plant, stepping in and hinging at the hips into a squat
// with the knees bent and the back fairly straight, reaching down with
// both hands to the plant, gripping it and pulling, which lifts it, and
// rising to stand with the crop held at the chest.
var pickBeats = []beat{
	{at: 0, ease: ease, pose: Pose{
		Lean: 0.1, Nod: 0.35, Hip: [2]float64{0.05, 0}, Knee: [2]float64{0.08, 0.06}, Abduct: both(0.1), Bob: -0.01, Shoulder: [2]float64{-0.13, -0.14}, Raise: [2]float64{-0.00, 0.08}, Plane: [2]float64{0.34, 0.96}, Elbow: [2]float64{0.42, 0.33}, Wrist: [2]float64{0, 0}}},
	{at: 0.2, ease: ease, pose: Pose{
		Lean: 0.5, Nod: 0.4, Hip: both(0.75), Knee: both(1.3), Abduct: both(0.14), Bob: -0.28, Forward: 0.15, Shoulder: [2]float64{0.42, 0.11}, Raise: [2]float64{-0.11, 0.39}, Plane: [2]float64{0.43, 1.60}, Elbow: [2]float64{0.43, 0.29}, Wrist: [2]float64{0, 0}}},
	{at: 0.38, ease: ease, pose: Pose{
		Lean: 0.75, Nod: 0.4, Hip: both(1.3), Knee: both(2.2), Abduct: both(0.16), Bob: -0.62, Forward: 0.3, Shoulder: [2]float64{0.67, 0.23}, Raise: [2]float64{-0.14, 0.58}, Plane: [2]float64{0.53, 1.74}, Elbow: [2]float64{0.55, 0.33}, Wrist: [2]float64{0, 0}}},
	{at: 0.52, ease: ease, pose: Pose{
		Lean: 0.8, Nod: 0.4, Hip: both(1.3), Knee: both(2.2), Abduct: both(0.16), Bob: -0.64, Forward: 0.3, Shoulder: [2]float64{0.82, 0.32}, Raise: [2]float64{-0.17, 0.59}, Plane: [2]float64{0.62, 1.81}, Elbow: [2]float64{0.51, 0.25}, Wrist: [2]float64{0, 0}}},
	{at: 0.68, ease: ease, ordered: true, pose: Pose{
		Lean: 0.5, Nod: 0.35, Hip: both(0.9), Knee: both(1.5), Abduct: both(0.14), Bob: -0.37, Forward: 0.2, Shoulder: [2]float64{0.61, 0.20}, Raise: [2]float64{-0.15, 0.52}, Plane: [2]float64{0.61, 1.56}, Elbow: [2]float64{0.58, 0.50}, Wrist: [2]float64{0, 0}}},
	{at: 0.85, ease: ease, ordered: true, pose: Pose{
		Lean: 0.12, Nod: 0.3, Hip: [2]float64{0.1, 0.02}, Knee: [2]float64{0.2, 0.1}, Abduct: both(0.1), Bob: -0.03, Forward: 0.05, Shoulder: [2]float64{0.15, -0.14}, Raise: [2]float64{-0.01, 0.61}, Plane: [2]float64{0.49, 0.97}, Elbow: [2]float64{1.29, 1.32}, Wrist: [2]float64{0, 0}}},
	{at: 0.95, ease: ease, pose: Pose{
		Lean: 0.1, Nod: 0.3, Hip: [2]float64{0.06, 0}, Knee: [2]float64{0.1, 0.07}, Abduct: both(0.1), Bob: -0.015, Forward: 0.01, Shoulder: [2]float64{0.12, -0.16}, Raise: [2]float64{0.02, 0.62}, Plane: [2]float64{0.47, 0.94}, Elbow: [2]float64{1.26, 1.26}, Wrist: [2]float64{0, 0}}},
}

// Pick is picking a crop: bending into a squat to the plant, pulling it
// up and rising with it.
var Pick = loop{"pick", 1.6, 40}.of(func(t float64) Pose { return poseAt(pickBeats, t) })
