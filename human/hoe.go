package human

// The beats of a stroke with the hoe, held in both hands, the left near
// the end of the handle and the right lower down it, the left foot
// forward and the knees bent: the blade lying on the ground it has
// drawn, lifted forward and over the shoulder with the trunk coming up,
// thrown down by the hips, trunk and then arms, to bite the soil a stride
// ahead, and drawn back toward the feet, which turns the soil.
var hoeBeats = []beat{
	{at: 0, ease: ease, pose: Pose{
		Lean: 0.75, Nod: 0.3, Twist: 0.0, Hip: [2]float64{0.7, -0.1}, Knee: [2]float64{0.95, 0.4}, Abduct: both(0.12), Bob: -0.14, Forward: 0.05, Shoulder: [2]float64{-0.32, 0.26}, Raise: [2]float64{0.70, 0.65}, Plane: [2]float64{0.96, 2.05}, Elbow: [2]float64{1.74, 0.14}, Wrist: [2]float64{0, 0}}},
	{at: 0.25, ease: ease, pose: Pose{
		Lean: 0.3, Nod: 0.3, Twist: -0.1, Hip: [2]float64{0.4, -0.05}, Knee: [2]float64{0.55, 0.25}, Abduct: both(0.12), Bob: -0.05, Forward: 0.03, Shoulder: [2]float64{-0.36, 0.21}, Raise: [2]float64{0.48, 0.58}, Plane: [2]float64{1.17, 2.25}, Elbow: [2]float64{1.00, 0.11}, Wrist: [2]float64{0, 0}}},
	{at: 0.42, ease: ease, pose: Pose{
		Lean: 0.15, Nod: 0.25, Twist: -0.2, Hip: [2]float64{0.35, 0}, Knee: [2]float64{0.4, 0.2}, Abduct: both(0.12), Bob: -0.02, Forward: 0, Shoulder: [2]float64{-0.21, 0.60}, Raise: [2]float64{0.57, 0.12}, Plane: [2]float64{1.23, 2.30}, Elbow: [2]float64{0.69, 1.25}, Wrist: [2]float64{0, 0}}},
	{at: 0.52, ease: hasten, ordered: true, pose: Pose{
		Lean: 0.5, Nod: 0.3, Twist: -0.05, Hip: [2]float64{0.55, -0.05}, Knee: [2]float64{0.7, 0.3}, Abduct: both(0.12), Bob: -0.08, Forward: 0.04, Shoulder: [2]float64{-0.27, 0.39}, Raise: [2]float64{0.64, 0.37}, Plane: [2]float64{0.95, 1.22}, Elbow: [2]float64{1.44, 1.76}, Wrist: [2]float64{0, 0}}},
	{at: 0.6, ease: straight, ordered: true, pose: Pose{
		Lean: 0.85, Nod: 0.35, Twist: 0.05, Hip: [2]float64{0.8, -0.1}, Knee: [2]float64{1.1, 0.45}, Abduct: both(0.12), Bob: -0.18, Forward: 0.06, Shoulder: [2]float64{0.21, 0.57}, Raise: [2]float64{0.95, 0.74}, Plane: [2]float64{1.14, 1.70}, Elbow: [2]float64{1.64, 0.41}, Wrist: [2]float64{0, 0}}},
	{at: 0.78, ease: ease, pose: Pose{
		Lean: 0.8, Nod: 0.3, Twist: 0.0, Hip: [2]float64{0.75, -0.1}, Knee: [2]float64{1.0, 0.42}, Abduct: both(0.12), Bob: -0.16, Forward: 0.05, Shoulder: [2]float64{-0.02, 0.46}, Raise: [2]float64{0.80, 0.67}, Plane: [2]float64{1.01, 1.95}, Elbow: [2]float64{1.74, 0.23}, Wrist: [2]float64{0, 0}}},
	{at: 0.9, ease: ease, pose: Pose{
		Lean: 0.75, Nod: 0.3, Twist: 0.0, Hip: [2]float64{0.7, -0.1}, Knee: [2]float64{0.95, 0.4}, Abduct: both(0.12), Bob: -0.14, Forward: 0.05, Shoulder: [2]float64{-0.29, 0.36}, Raise: [2]float64{0.68, 0.56}, Plane: [2]float64{0.97, 2.08}, Elbow: [2]float64{1.72, 0.17}, Wrist: [2]float64{0, 0}}},
}

// Hoe is working the ground with a hoe: lifting it, bringing it down to
// bite the soil a stride ahead, and drawing it back.
var Hoe = loop{"hoe", 1.4, 40}.of(func(t float64) Pose { return poseAt(hoeBeats, t) })
