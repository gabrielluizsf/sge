package human

// The beats of watering a plot: the can carried upright at the right hip
// by its handle, lifted and tipped forward until its spout points down at
// the soil a stride ahead, held there with a slight sway while the water
// runs, and tipped back, lowered and carried again. The left arm rests,
// a little out for balance while the can is tipped.
var waterBeats = []beat{
	{at: 0, ease: ease, pose: Pose{
		Lean: 0.12, Nod: 0.25, Twist: 0.0, Hip: [2]float64{0.08, 0}, Knee: [2]float64{0.1, 0.08}, Abduct: both(0.1), Bob: -0.01,
		Shoulder: [2]float64{0.05, -0.50}, Raise: [2]float64{0.10, 0.24}, Plane: [2]float64{0.00, 0.42}, Elbow: [2]float64{0.30, 0.98}, Wrist: [2]float64{0.00, 0.90}}},
	{at: 0.2, ease: ease, pose: Pose{
		Lean: 0.2, Nod: 0.3, Twist: 0.03, Hip: [2]float64{0.1, 0}, Knee: [2]float64{0.12, 0.08}, Abduct: both(0.1), Bob: -0.015,
		Shoulder: [2]float64{0.05, -0.50}, Raise: [2]float64{0.18, 0.18}, Plane: [2]float64{0.00, 0.45}, Elbow: [2]float64{0.30, 1.08}, Wrist: [2]float64{0.00, 0.90}}},
	{at: 0.35, ease: ease, pose: Pose{
		Lean: 0.28, Nod: 0.4, Twist: 0.06, Hip: [2]float64{0.12, 0}, Knee: [2]float64{0.14, 0.1}, Abduct: both(0.1), Bob: -0.02,
		Shoulder: [2]float64{0.05, -0.34}, Raise: [2]float64{0.28, -0.04}, Plane: [2]float64{0.00, 0.43}, Elbow: [2]float64{0.35, 1.28}, Wrist: [2]float64{0.00, 0.10}}},
	{at: 0.55, ease: ease, pose: Pose{
		Lean: 0.32, Nod: 0.45, Twist: 0.04, Hip: [2]float64{0.12, 0}, Knee: [2]float64{0.14, 0.1}, Abduct: both(0.1), Bob: -0.02,
		Shoulder: [2]float64{0.05, -0.27}, Raise: [2]float64{0.30, -0.17}, Plane: [2]float64{0.00, 0.45}, Elbow: [2]float64{0.35, 1.55}, Wrist: [2]float64{0.00, -0.37}}},
	{at: 0.72, ease: ease, pose: Pose{
		Lean: 0.3, Nod: 0.45, Twist: 0.08, Hip: [2]float64{0.12, 0}, Knee: [2]float64{0.14, 0.1}, Abduct: both(0.1), Bob: -0.02,
		Shoulder: [2]float64{0.05, -0.25}, Raise: [2]float64{0.28, -0.17}, Plane: [2]float64{0.00, 0.42}, Elbow: [2]float64{0.35, 1.50}, Wrist: [2]float64{0.00, -0.27}}},
	{at: 0.82, ease: ease, pose: Pose{
		Lean: 0.22, Nod: 0.32, Twist: 0.03, Hip: [2]float64{0.1, 0}, Knee: [2]float64{0.12, 0.08}, Abduct: both(0.1), Bob: -0.015,
		Shoulder: [2]float64{0.05, -0.50}, Raise: [2]float64{0.18, 0.20}, Plane: [2]float64{0.00, 0.45}, Elbow: [2]float64{0.30, 1.12}, Wrist: [2]float64{0.00, 0.79}}},
	{at: 0.92, ease: ease, pose: Pose{
		Lean: 0.14, Nod: 0.26, Twist: 0.0, Hip: [2]float64{0.08, 0}, Knee: [2]float64{0.1, 0.08}, Abduct: both(0.1), Bob: -0.01,
		Shoulder: [2]float64{0.05, -0.50}, Raise: [2]float64{0.12, 0.24}, Plane: [2]float64{0.00, 0.41}, Elbow: [2]float64{0.30, 1.00}, Wrist: [2]float64{0.00, 0.90}}},
}

// Water is watering with a can: lifting it, tipping it to pour on the
// soil ahead, and carrying it again.
var Water = loop{"water", 2.0, 40}.of(func(t float64) Pose { return poseAt(waterBeats, t) })

// Farming are the motions of those who work a plot, besides Motions.
var Farming = []Motion{Water, Pick}
