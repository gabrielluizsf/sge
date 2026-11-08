package human

// HoeGrip is how the hoe lies in the hands that work it: the left at the
// end of the handle, the right 22 cm down it toward the blade, whose wide
// face is across the body. It is out of sight in every other motion.
var HoeGrip = Grip{
	Hands:  [2]string{"wrist.L", "wrist.R"},
	During: []string{"hoe"},
	End:    -0.52,
	Facing: xAxis,
}
