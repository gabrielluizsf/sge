package terrain

import "math"

// Noise is value noise: a number from 0 to 1 at x, y that changes
// smoothly, from random values at whole points eased between.
func Noise(seed uint32, x, y float64) float64 {
	x0, y0 := math.Floor(x), math.Floor(y)
	c, r := int(x0), int(y0)
	sx, sy := ease(x-x0), ease(y-y0)
	top := lerp(Chance(seed, c, r), Chance(seed, c+1, r), sx)
	bottom := lerp(Chance(seed, c, r+1), Chance(seed, c+1, r+1), sx)
	return lerp(top, bottom, sy)
}

// ease is smoothstep: flat at 0 and 1, so the noise has no creases.
func ease(t float64) float64 { return t * t * (3 - 2*t) }

func lerp(a, b, t float64) float64 { return a + (b-a)*t }
