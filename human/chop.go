package human

// The beats of a felling cut, a right-handed swing at a trunk ahead, the
// left foot forward: the ready stance, the long backswing over the right
// shoulder with the trunk turned and the weight on the back foot, the
// strike, driven by hips, trunk and then arms, the follow-through, and the
// recovery to the stance. The arms' joints were fitted so that the hands
// are on the axe's handle, the left at its end and the right up it, close
// together as the strike nears.
var chopBeats = []beat{
	{at: 0, ease: ease, pose: Pose{
		Lean: 0.2, Nod: 0.3, Twist: -0.1, Hip: [2]float64{0.3, -0.05}, Knee: [2]float64{0.3, 0.25}, Abduct: both(0.15), Bob: -0.04,
		Shoulder: [2]float64{0.50, 0.64}, Raise: [2]float64{0.10, 0.51}, Plane: [2]float64{2.21, 2.30}, Elbow: [2]float64{0.41, 0.13}, Wrist: [2]float64{0, 0}}},
	{at: 0.3, ease: ease, pose: Pose{
		Lean: 0.1, Nod: 0.25, Twist: -0.35, Hip: [2]float64{0.3, 0.05}, Knee: [2]float64{0.2, 0.32}, Abduct: both(0.15), Bob: -0.045, Forward: -0.02,
		Shoulder: [2]float64{-0.30, 0.86}, Raise: [2]float64{0.57, 0.60}, Plane: [2]float64{0.59, 1.85}, Elbow: [2]float64{1.66, 1.06}, Wrist: [2]float64{0, 0}}},
	{at: 0.55, ease: ease, pose: Pose{
		Lean: 0.0, Nod: 0.2, Twist: -0.6, Hip: [2]float64{0.3, 0.1}, Knee: [2]float64{0.15, 0.4}, Abduct: both(0.15), Bob: -0.05, Forward: -0.04,
		Shoulder: [2]float64{0.16, 1.52}, Raise: [2]float64{0.55, 0.75}, Plane: [2]float64{1.20, 0.48}, Elbow: [2]float64{2.50, 2.50}, Wrist: [2]float64{0, 0}}},
	{at: 0.65, ease: hasten, ordered: true, pose: Pose{
		Lean: 0.3, Nod: 0.3, Twist: 0.0, Hip: [2]float64{0.45, -0.05}, Knee: [2]float64{0.35, 0.35}, Abduct: both(0.15), Bob: -0.06,
		Shoulder: [2]float64{0.93, 1.70}, Raise: [2]float64{0.64, 0.70}, Plane: [2]float64{1.13, 1.06}, Elbow: [2]float64{2.49, 1.92}, Wrist: [2]float64{0, 0}}},
	{at: 0.75, ease: straight, ordered: true, pose: Pose{
		Lean: 0.5, Nod: 0.35, Twist: 0.3, Hip: [2]float64{0.55, -0.1}, Knee: [2]float64{0.5, 0.3}, Ankle: [2]float64{0, 0.25}, Abduct: both(0.15), Bob: -0.08, Forward: 0.05,
		Shoulder: [2]float64{0.39, 0.86}, Raise: [2]float64{0.59, -0.12}, Plane: [2]float64{2.00, 1.21}, Elbow: [2]float64{0.91, 0.50}, Wrist: [2]float64{0, 0}}},
	{at: 0.825, ease: settle, pose: Pose{
		Lean: 0.6, Nod: 0.35, Twist: 0.4, Hip: [2]float64{0.55, -0.1}, Knee: [2]float64{0.5, 0.3}, Ankle: [2]float64{0, 0.25}, Abduct: both(0.15), Bob: -0.08, Forward: 0.05,
		Shoulder: [2]float64{0.29, 0.93}, Raise: [2]float64{0.74, -0.17}, Plane: [2]float64{2.30, 1.19}, Elbow: [2]float64{0.40, 0.38}, Wrist: [2]float64{0, 0}}},
	{at: 0.875, ease: ease, pose: Pose{
		Lean: 0.4, Nod: 0.33, Twist: 0.2, Hip: [2]float64{0.4, -0.05}, Knee: [2]float64{0.4, 0.27}, Ankle: [2]float64{0, 0.1}, Abduct: both(0.15), Bob: -0.06, Forward: 0.02,
		Shoulder: [2]float64{0.28, 0.62}, Raise: [2]float64{0.60, 0.30}, Plane: [2]float64{2.30, 2.07}, Elbow: [2]float64{0.30, 0.27}, Wrist: [2]float64{0, 0}}},
}

// Chop is chopping wood with an axe held in both hands, in a felling cut.
var Chop = loop{"chop", 1.4, 40}.of(func(t float64) Pose { return poseAt(chopBeats, t) })
