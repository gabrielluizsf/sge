package human

// The can's hold while it is filled: in the right hand by its handle, as
// in the pour, at the well under the pail with the water falling into its
// mouth, and at the river sunk in the water.
var (
	FillGrip = Grip{
		Hands:  CanGrip.Hands,
		During: []string{"well", "dip"},
		Offset: CanGrip.Offset,
		Turn:   CanGrip.Turn,
	}
	// WellFall is the water poured into the can from the pail over it,
	// while the left hand tips the pail.
	WellFall = Fall{Can: FillGrip, Mouth: [3]float64{0, 0.3, 0}, Height: 0.35, From: 0.62, To: 0.82, During: []string{"well"}}
)
