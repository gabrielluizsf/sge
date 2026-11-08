package human

// ChopGrip is how the axe lies in the hands that chop: the left at the
// handle's end, the right up it, the blade's edge facing along the
// person's facing, at the tree. It is out of sight in every other motion.
var ChopGrip = Grip{
	Hands:  [2]string{"wrist.L", "wrist.R"},
	During: []string{"chop"},
	End:    -0.42,
	Facing: zAxis,
	Flat:   0.85,
}
