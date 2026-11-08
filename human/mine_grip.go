package human

// MineGrip is how the pickaxe lies in the hands that swing it: the left
// at the handle's end, the right up it, the points facing along the
// person's facing, at the rock. It is out of sight in every other motion.
var MineGrip = Grip{
	Hands:  [2]string{"wrist.L", "wrist.R"},
	During: []string{"mine"},
	End:    -0.42,
	Facing: zAxis,
	Flat:   0.85,
}
