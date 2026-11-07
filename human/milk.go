package human

// A seat as low as a milking stool: the thighs rising to the knees, the
// shanks folded back under them, the body bent forward over the knees.
const (
	milkHip  = 1.75
	milkKnee = 2.05
	milkLean = 0.3
)

// milkPose is sitting on the stool, the arms reaching forward and down,
// each hand as far along its pull as pull says: 0 up at the teat, 1 drawn
// down.
func milkPose(pull [2]float64, breath float64) Pose {
	p := Pose{Lean: milkLean + 0.01*breath, Nod: 0.2, Hip: both(milkHip), Knee: both(milkKnee), Ankle: both(-milkHip + milkKnee - 0.05),
		Abduct: both(0.3), Plane: both(0.1), Raise: both(0.06), Wrist: both(0.1)}
	for i, u := range pull {
		p.Shoulder[i] = 0.62 - 0.1*u
		p.Elbow[i] = 1.5 - 0.22*u
	}
	p.Bob = legReach(p.Hip[0], p.Knee[0]) - (thighLen + shankLen)
	return p
}

// Milk is milking, seated on the stool: the hands draw down in turn.
var Milk = loop{"milk", 1.2, 16}.of(func(t float64) Pose {
	return milkPose([2]float64{(1 + wave(t, 1, 0)) / 2, (1 + wave(t, 1, 0.5)) / 2}, wave(t, 1, 0.25))
})

// MilkDown is sitting down to milk, from standing to the stool with the
// hands at the teats; played from its end to its start, it is getting up.
var MilkDown = loop{"milkdown", 0.8, 8}.once(func(t float64) Pose {
	return mix(Stand.Poses[0], Milk.Poses[0], ease(t))
})

// Milking are the motions of those who milk, besides Motions.
var Milking = []Motion{MilkDown, Milk}

// MilkVersion changes whenever the milking does, or the stool and pail
// under it, so what was baked of it is baked again, and no other motion is.
const MilkVersion = 1
