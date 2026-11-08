package water

import "math"

// linear is the light of an sRGB component, 0 to 1.
func linear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

// srgb is the sRGB component of linear light c, clamped to 0 to 1.
func srgb(c float64) float64 {
	c = math.Max(0, math.Min(c, 1))
	if c <= 0.0031308 {
		return c * 12.92
	}
	return 1.055*math.Pow(c, 1/2.4) - 0.055
}

// mix is a at t = 0 and b at t = 1.
func mix(a, b, t float64) float64 { return a + (b-a)*t }

// mix3 is mix of each component of a and b.
func mix3(a, b [3]float64, t float64) [3]float64 {
	return [3]float64{mix(a[0], b[0], t), mix(a[1], b[1], t), mix(a[2], b[2], t)}
}

// smooth is 0 at or below lo, 1 at or above hi, and eased between.
func smooth(lo, hi, x float64) float64 {
	t := math.Max(0, math.Min((x-lo)/(hi-lo), 1))
	return t * t * (3 - 2*t)
}
