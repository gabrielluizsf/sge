package human

// The beats of drinking from a bottle: held at the chest in the right hand,
// its neck up; brought to the lips and tipped, the head going back with it for
// a long pull, the throat working; and brought down again with a breath. The
// left hand rests on the hip all the while.
var drinkBeats = []beat{
	{at: 0, ease: ease, pose: sip{nod: 0.08, shoulder: 0.5, elbow: 1.7, wrist: -0.1}.pose()},
	{at: 0.18, ease: ease, pose: sip{nod: 0, shoulder: 0.8, elbow: 2.15, wrist: 0.1}.pose()},
	{at: 0.3, ease: ease, pose: sip{nod: -0.28, lean: -0.05, shoulder: 0.98, elbow: 2.2, wrist: 0.4}.pose()},
	{at: 0.5, ease: ease, pose: sip{nod: -0.36, lean: -0.07, shoulder: 1.06, elbow: 2.15, wrist: 0.55}.pose()},
	{at: 0.7, ease: ease, pose: sip{nod: -0.42, lean: -0.08, shoulder: 1.12, elbow: 2.1, wrist: 0.68}.pose()},
	{at: 0.84, ease: ease, pose: sip{nod: -0.02, shoulder: 0.7, elbow: 2.0, wrist: 0.05}.pose()},
}

// sip is a moment of drinking: the head bowed nod, or tipped back when it is
// less than none, the body leaning lean, the right arm raised shoulder with
// its elbow bent elbow, and the wrist bent wrist.
type sip struct{ nod, lean, shoulder, elbow, wrist float64 }

// pose is the whole pose of the sip, the left hand on the hip.
func (s sip) pose() Pose {
	return Pose{
		Lean: 0.03 + s.lean, Nod: s.nod, Hip: both(0.02), Knee: both(0.04), Abduct: both(0.09),
		Shoulder: [2]float64{-0.2, s.shoulder}, Raise: [2]float64{0.55, 0.22}, Plane: [2]float64{0.1, 0.45},
		Elbow: [2]float64{1.75, s.elbow}, Wrist: [2]float64{0.2, s.wrist},
	}
}

// Drink is drinking from a bottle in the right hand.
var Drink = loop{"drink", 2.6, 40}.of(func(t float64) Pose { return poseAt(drinkBeats, t) })

// Drinking are the motions of those who drink, besides Motions.
var Drinking = []Motion{Drink}

// DrinkVersion changes whenever the drinking does, or the bottle in the hand,
// so what was baked of it is baked again, and no other motion is.
const DrinkVersion = 2

// BottleGrip is how the bottle lies in the right hand while it is drunk
// from: gripped about its middle, its neck toward the mouth, and out of sight
// in the other motions.
var BottleGrip = Grip{
	Hands:  [2]string{"wrist.R", "wrist.R"},
	During: []string{"drink"},
	Offset: [3]float64{0.01, -0.08, 0.1},
	Turn:   [4]float64{-0.7071068, 0, 0, 0.7071068},
}
