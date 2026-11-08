package water

import "image"

// bedAt is the colour, in linear light, of the bed texture at the world
// point (wu, wv) of the repeat, each 0 to 1. The texture is laid out
// mirrored, so the world's repeat meets itself at its edges without a
// seam, whatever the picture.
func bedAt(bed *image.NRGBA, wu, wv float64) [3]float64 {
	b := bed.Bounds()
	x := int(pingPong(wu) * float64(b.Dx()-1))
	y := int(pingPong(wv) * float64(b.Dy()-1))
	c := bed.NRGBAAt(b.Min.X+x, b.Min.Y+y)
	return [3]float64{linear(float64(c.R) / 255), linear(float64(c.G) / 255), linear(float64(c.B) / 255)}
}

// pingPong is t, 0 to 1, run out to the far end of a picture and back,
// the same at 0 and at 1.
func pingPong(t float64) float64 {
	if t > 0.5 {
		t = 1 - t
	}
	return 2 * t
}
