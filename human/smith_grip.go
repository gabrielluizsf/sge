package human

// The smith's tools in their hands, shown only while they forge. The
// hammer is held in the right fist by its handle, its flat face down at
// the stroke; the tongs are held in the left by their reins, the iron
// they hold laid on the anvil's face, 52 cm ahead and 83 cm high. The
// offsets and turns were fitted to the poses of Forge.
var (
	HammerGrip = Grip{
		Hands:  [2]string{"wrist.R", "wrist.R"},
		During: []string{"forge"},
		Offset: [3]float64{-0.1196, -0.0020, -0.0125},
		Turn:   [4]float64{-0.17134199, -0.80701186, -0.54641428, -0.14424029},
	}
	TongsGrip = Grip{
		Hands:  [2]string{"wrist.L", "wrist.L"},
		During: []string{"forge"},
		Offset: [3]float64{0.1695, 0.0591, 0.1235},
		Turn:   [4]float64{0.05038785, -0.26385989, 0.29956836, 0.91547683},
	}
)
