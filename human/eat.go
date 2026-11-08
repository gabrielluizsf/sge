package human

// The beats of eating a bite: standing relaxed with the food held at the
// chest in the right hand, the forearm folding up to bring it to the
// mouth, the head bowing the last bit to meet it and taking the bite,
// the hand coming back down to the chest while the jaw works with a
// slight nod three times, and a lift of the chin to swallow. The left arm
// hangs loose, a little out from the side.
var eatBeats = []beat{
	{at: 0, ease: ease, pose: bite{nod: 0.05, shoulder: 0.55, elbow: 1.55, raise: 0.10, wrist: 0.12, twist: 0.0}.pose()},
	{at: 0.18, ease: ease, pose: bite{nod: 0.12, shoulder: 1.05, elbow: 2.25, raise: 0.12, wrist: 0.20, twist: 0.0}.pose()},
	{at: 0.28, ease: ease, pose: bite{nod: 0.20, shoulder: 1.10, elbow: 2.35, raise: 0.12, wrist: 0.30, twist: 0.02}.pose()},
	{at: 0.40, ease: ease, pose: bite{nod: 0.10, shoulder: 0.75, elbow: 1.85, raise: 0.11, wrist: 0.14, twist: 0.0}.pose()},
	{at: 0.48, ease: ease, pose: bite{nod: 0.14, shoulder: 0.58, elbow: 1.60, raise: 0.10, wrist: 0.17, twist: 0.02}.pose()},
	{at: 0.56, ease: ease, pose: bite{nod: 0.06, shoulder: 0.55, elbow: 1.55, raise: 0.10, wrist: 0.12, twist: -0.02}.pose()},
	{at: 0.64, ease: ease, pose: bite{nod: 0.14, shoulder: 0.55, elbow: 1.55, raise: 0.10, wrist: 0.17, twist: 0.02}.pose()},
	{at: 0.72, ease: ease, pose: bite{nod: 0.06, shoulder: 0.55, elbow: 1.55, raise: 0.10, wrist: 0.12, twist: -0.02}.pose()},
	{at: 0.80, ease: ease, pose: bite{nod: 0.14, shoulder: 0.55, elbow: 1.55, raise: 0.10, wrist: 0.17, twist: 0.02}.pose()},
	{at: 0.90, ease: ease, pose: bite{nod: -0.06, shoulder: 0.55, elbow: 1.55, raise: 0.10, wrist: 0.10, twist: 0.0}.pose()},
}

// bite is a moment of eating: the head bowed nod, the right arm raised
// shoulder forward with its elbow bent elbow and its upper arm out raise,
// the wrist bent wrist, and the trunk turned twist.
type bite struct{ nod, shoulder, elbow, raise, wrist, twist float64 }

// pose is the whole pose of the bite, the left arm hanging loose.
func (b bite) pose() Pose {
	return Pose{
		Lean: 0.04 + b.nod*0.2, Nod: b.nod, Twist: b.twist, Hip: both(0.02), Knee: both(0.03), Abduct: both(0.1),
		Shoulder: [2]float64{0.04, b.shoulder}, Raise: [2]float64{0.12, b.raise}, Plane: [2]float64{0, 0.7},
		Elbow: [2]float64{0.25, b.elbow}, Wrist: [2]float64{0, b.wrist},
	}
}

// Eat is eating a bite: the food to the mouth, chewed, and swallowed.
var Eat = loop{"eat", 2.2, 40}.of(func(t float64) Pose { return poseAt(eatBeats, t) })

// Eating are the motions of those who eat, besides Motions.
var Eating = []Motion{Eat}

// EatVersion changes whenever the eating does, or the food in the hand,
// so what was baked of it is baked again, and no other motion is.
const EatVersion = 1
