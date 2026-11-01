package water

// Preset is a kind of water: how it looks and how choppy its waves are.
type Preset struct {
	Name      string
	Surface   Surface
	Steepness float64
}

// The sky's colours and the sun's, in linear light, and the way toward
// the sun in the game's world, as package bake lights its pictures.
var (
	sunward  = unit([3]float64{-0.35, 0.6, 1})
	horizon  = [3]float64{0.62, 0.72, 0.82}
	zenith   = [3]float64{0.16, 0.34, 0.68}
	sunlight = [3]float64{0.8, 0.75, 0.65}
)

// Lake is a lake's water: green-brown, a clear metre, small ripples.
func Lake() Preset {
	return Preset{Name: "lake", Steepness: 0.7, Surface: Surface{
		Shallow: [3]float64{0.07, 0.17, 0.12}, Deep: [3]float64{0.012, 0.07, 0.08},
		Clear: 0.7, MaxDepth: 1.5, Horizon: horizon, Zenith: zenith,
		Sun: sunward, SunLight: sunlight, Shine: 90, F0: 0.03, Foam: [3]float64{0.8, 0.85, 0.85}}}
}

// Sea is a tropical sea's water: turquoise over pale sand, clearer, and
// choppier.
func Sea() Preset {
	return Preset{Name: "sea", Steepness: 1.1, Surface: Surface{
		Shallow: [3]float64{0.02, 0.38, 0.40}, Deep: [3]float64{0.003, 0.07, 0.17},
		Clear: 1.4, MaxDepth: 3, Horizon: horizon, Zenith: zenith,
		Sun: sunward, SunLight: sunlight, Shine: 70, F0: 0.03, Foam: [3]float64{0.92, 0.95, 0.95}}}
}
