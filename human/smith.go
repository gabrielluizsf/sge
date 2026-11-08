package human

// anvilStep is how far the smith has stepped from where they stand, the
// middle of their cell, toward the anvil, which stands in the next cell, a
// metre on: the iron lies on its face 52 cm from where they then stand, a
// metre from the cell's middle.
const anvilStep = 0.48

// The beats of a stroke at the anvil, the left foot a little forward, the
// feet apart and the knees soft: the hammer lying on the iron, lifted
// elbow first, cocked back at the height of the head at the top, thrown
// down by the wrist as the trunk comes over it, struck on the iron at the
// anvil's face, and bounced up from it while the left hand, which holds
// the iron in tongs, lifts it a little and sets it down again.
var smithBeats = []beat{
	{at: 0, ease: ease, pose: Pose{
		Lean: 0.25, Nod: 0.45, Twist: -0.05, Hip: [2]float64{0.12, 0}, Knee: [2]float64{0.15, 0.12}, Abduct: both(0.2), Bob: -0.02, Forward: anvilStep, Shoulder: [2]float64{-0.40, -0.07}, Raise: [2]float64{0.48, 0.52}, Plane: [2]float64{1.12, 1.73}, Elbow: [2]float64{0.89, 0.65}, Wrist: [2]float64{-0.04, -0.24}}},
	{at: 0.25, ease: ease, pose: Pose{
		Lean: 0.18, Nod: 0.4, Twist: -0.15, Hip: [2]float64{0.12, 0}, Knee: [2]float64{0.13, 0.1}, Abduct: both(0.2), Bob: -0.01, Forward: anvilStep, Shoulder: [2]float64{-0.37, 0.95}, Raise: [2]float64{0.36, -0.17}, Plane: [2]float64{1.02, 1.86}, Elbow: [2]float64{0.89, 1.54}, Wrist: [2]float64{-0.05, -0.90}}},
	{at: 0.5, ease: ease, pose: Pose{
		Lean: 0.1, Nod: 0.12, Twist: -0.25, Hip: [2]float64{0.12, 0}, Knee: [2]float64{0.1, 0.08}, Abduct: both(0.2), Ankle: both(0.1), Bob: 0, Forward: anvilStep, Shoulder: [2]float64{-0.37, 0.12}, Raise: [2]float64{0.29, 1.75}, Plane: [2]float64{0.94, 0.6}, Elbow: [2]float64{0.88, 2.50}, Wrist: [2]float64{-0.07, -0.90}}},
	{at: 0.62, ease: hasten, ordered: true, withTrunk: true, pose: Pose{
		Lean: 0.25, Nod: 0.45, Twist: -0.1, Hip: [2]float64{0.14, 0}, Knee: [2]float64{0.18, 0.15}, Abduct: both(0.2), Bob: -0.03, Forward: anvilStep, Shoulder: [2]float64{-0.32, -0.50}, Raise: [2]float64{0.37, 1.08}, Plane: [2]float64{1.06, 1.49}, Elbow: [2]float64{0.89, 1.55}, Wrist: [2]float64{-0.03, -0.82}}},
	{at: 0.7, ease: straight, ordered: true, withTrunk: true, pose: Pose{
		Lean: 0.38, Nod: 0.5, Twist: 0.05, Hip: [2]float64{0.16, 0}, Knee: [2]float64{0.2, 0.18}, Abduct: both(0.2), Bob: -0.04, Forward: anvilStep, Shoulder: [2]float64{-0.34, -0.47}, Raise: [2]float64{0.52, 0.99}, Plane: [2]float64{1.20, 1.66}, Elbow: [2]float64{0.89, 0.53}, Wrist: [2]float64{0.00, -0.16}}},
	{at: 0.78, ease: settle, pose: Pose{
		Lean: 0.3, Nod: 0.45, Twist: 0, Hip: [2]float64{0.14, 0}, Knee: [2]float64{0.17, 0.14}, Abduct: both(0.2), Bob: -0.03, Forward: anvilStep, Shoulder: [2]float64{-0.38, -0.15}, Raise: [2]float64{0.50, 0.79}, Plane: [2]float64{1.16, 1.85}, Elbow: [2]float64{0.89, 0.45}, Wrist: [2]float64{-0.03, -0.32}}},
	{at: 0.9, ease: ease, pose: Pose{
		Lean: 0.27, Nod: 0.45, Twist: -0.03, Hip: [2]float64{0.12, 0}, Knee: [2]float64{0.15, 0.12}, Abduct: both(0.2), Bob: -0.02, Forward: anvilStep, Shoulder: [2]float64{-0.43, -0.13}, Raise: [2]float64{0.51, 0.62}, Plane: [2]float64{1.10, 1.77}, Elbow: [2]float64{0.97, 0.55}, Wrist: [2]float64{0.00, -0.22}}},
}

// Forge is a smith's stroke at the anvil, once a second.
var Forge = loop{"forge", 1.0, 40}.of(func(t float64) Pose { return poseAt(smithBeats, t) })

// Crafts are the motions of those who work at a trade, besides Motions.
var Crafts = []Motion{Forge}
