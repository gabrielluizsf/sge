package human

// MorselGrip is how the food lies in the right hand while it is eaten,
// held in the palm and out of sight in the other motions.
var MorselGrip = Grip{
	Hands:  [2]string{"wrist.R", "wrist.R"},
	During: []string{"eat"},
	Offset: [3]float64{0, -0.09, 0.01},
	Turn:   [4]float64{0, 0, 0, 1},
}
