package human

// The watering can's own measures, in its frame: where its handle is held,
// where its spout ends and which way the spout points there, as
// building.Can has them.
var (
	canHandle = [3]float64{-0.2, 0.22, 0}
	canSpout  = [3]float64{0.55, 0.39, 0}
	canOut    = [3]float64{0.88, 0.48, 0}
)

// CanGrip is how the watering can lies in the right hand, by its handle,
// and CanSpill the stream it pours, shown while its spout points down;
// both only while watering. The offset and turn were fitted to the poses
// of Water, so the can is carried upright and tipped over the soil.
var (
	CanGrip = Grip{
		Hands:  [2]string{"wrist.R", "wrist.R"},
		During: []string{"water"},
		Offset: [3]float64{-0.2666, -0.2098, 0.0156},
		Turn:   [4]float64{0.08307322, 0.84687064, 0.39535492, -0.34583729},
	}
	CanSpill = Spill{Can: CanGrip, Spout: canSpout, Out: canOut, Tilt: 0.3, During: []string{"water"}}
)
