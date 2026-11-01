package water

import "math"

// Surface is how water looks, in linear light: its colours, how far one
// sees into it, and the sky and sun it mirrors.
type Surface struct {
	// Shallow and Deep are the water's own colour over the bed and far
	// from it.
	Shallow, Deep [3]float64
	// Clear is how deep, in metres, the bed shows to 1/e of itself.
	Clear float64
	// MaxDepth is the depth water reaches away from the shore, in metres.
	MaxDepth float64
	// Horizon and Zenith are the sky's colours at its edge and overhead.
	Horizon, Zenith [3]float64
	// Sun is the way toward the sun, as View is the way toward the camera;
	// SunLight its colour, and Shine how small its glints are.
	Sun, SunLight [3]float64
	Shine         float64
	// F0 is the mirrored share seen straight on, Foam the foam's colour.
	F0   float64
	Foam [3]float64
}

// Look is a pixel of water to shade: the surface's normal there, its
// depth in metres, the bed's colour under it in linear light, and how
// much foam is on it, 0 to 1.
type Look struct {
	Normal [3]float64
	Depth  float64
	Bed    [3]float64
	Foam   float64
}

// Shade is the colour, in linear light, of water l as the camera sees it:
// the bed seen through the water, fading into its colour with depth,
// under the sky mirrored by Fresnel's law, the sun's glints, and the foam.
func (s *Surface) Shade(l Look) [3]float64 {
	n := l.Normal
	seen := math.Max(dot(n, View), 0.02)
	body := s.body(l)
	f := Fresnel(seen, s.F0)
	r := unit([3]float64{2*seen*n[0] - View[0], 2*seen*n[1] - View[1], 2*seen*n[2] - View[2]})
	sky := mix3(s.Horizon, s.Zenith, math.Sqrt(math.Max(r[2], 0)))
	h := unit([3]float64{s.Sun[0] + View[0], s.Sun[1] + View[1], s.Sun[2] + View[2]})
	glint := math.Pow(math.Max(dot(n, h), 0), s.Shine) * smooth(0.02, 0.25, l.Depth)
	var out [3]float64
	for i := range out {
		c := body[i]*(1-f) + sky[i]*f + s.SunLight[i]*glint
		out[i] = mix(c, s.Foam[i]*(0.55+0.45*n[2]), l.Foam)
	}
	return out
}

// body is the colour of the water itself: the bed seen through it, fading
// to the water's own colour with depth.
func (s *Surface) body(l Look) [3]float64 {
	own := mix3(s.Shallow, s.Deep, smooth(0, s.MaxDepth, l.Depth))
	through := math.Exp(-l.Depth / math.Max(s.Clear, 1e-6))
	return mix3(own, l.Bed, through)
}
