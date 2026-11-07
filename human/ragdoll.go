package human

// Limb is a part of a ragdoll: the bone it follows, the part it hangs from,
// empty for the first, how thick it is, in metres, and how far it may turn
// from its rest about its joint, in radians, as a body's joints go.
type Limb struct {
	Bone, Parent string
	Radius       float64
	Swing, Twist float64
}

// Ragdoll are the bones that make a body fall as a body does, parents before
// their children, each with the range of its joint: the few a physics engine
// needs of the skeleton's many.
var Ragdoll = []Limb{
	{Bone: "spine05", Radius: 0.14},
	{Bone: "spine03", Parent: "spine05", Radius: 0.14, Swing: 0.5, Twist: 0.35},
	{Bone: "spine01", Parent: "spine03", Radius: 0.15, Swing: 0.4, Twist: 0.35},
	{Bone: "head", Parent: "spine01", Radius: 0.11, Swing: 0.7, Twist: 0.9},
	{Bone: "upperarm01.L", Parent: "spine01", Radius: 0.05, Swing: 2.2, Twist: 1.2},
	{Bone: "lowerarm01.L", Parent: "upperarm01.L", Radius: 0.04, Swing: 2.4, Twist: 0.1},
	{Bone: "wrist.L", Parent: "lowerarm01.L", Radius: 0.04, Swing: 1.1, Twist: 1.2},
	{Bone: "upperarm01.R", Parent: "spine01", Radius: 0.05, Swing: 2.2, Twist: 1.2},
	{Bone: "lowerarm01.R", Parent: "upperarm01.R", Radius: 0.04, Swing: 2.4, Twist: 0.1},
	{Bone: "wrist.R", Parent: "lowerarm01.R", Radius: 0.04, Swing: 1.1, Twist: 1.2},
	{Bone: "upperleg01.L", Parent: "spine05", Radius: 0.08, Swing: 1.6, Twist: 0.6},
	{Bone: "lowerleg01.L", Parent: "upperleg01.L", Radius: 0.06, Swing: 2.3, Twist: 0.1},
	{Bone: "foot.L", Parent: "lowerleg01.L", Radius: 0.05, Swing: 0.8, Twist: 0.3},
	{Bone: "upperleg01.R", Parent: "spine05", Radius: 0.08, Swing: 1.6, Twist: 0.6},
	{Bone: "lowerleg01.R", Parent: "upperleg01.R", Radius: 0.06, Swing: 2.3, Twist: 0.1},
	{Bone: "foot.R", Parent: "lowerleg01.R", Radius: 0.05, Swing: 0.8, Twist: 0.3},
}
