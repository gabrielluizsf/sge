package water

import "math"

// View is the way toward the camera in the game's world, X along a row,
// Y down a column, Z up: the isometric corner, as package bake renders.
var View = unit([3]float64{1, 1, 1})

// Fresnel is how much light a surface mirrors, by Schlick's
// approximation, when seen at an angle whose cosine to the normal is
// cos: f0 straight on, all of it at the grazing angle. Water's f0 is
// about 0.02.
func Fresnel(cos, f0 float64) float64 {
	return f0 + (1-f0)*math.Pow(1-math.Max(0, math.Min(cos, 1)), 5)
}
