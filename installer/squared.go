package installer

import (
	"image"
	"image/color"
	"image/draw"
)

// squared is img centred on a clear square of side pixels: each pixel
// the average, weighted by how much of it it covers and by alpha, of the
// source pixels under it.
func squared(img image.Image, side int) *image.NRGBA {
	b := img.Bounds()
	own := max(b.Dx(), b.Dy())
	src := image.NewNRGBA(image.Rect(0, 0, own, own))
	at := image.Pt((own-b.Dx())/2, (own-b.Dy())/2)
	draw.Draw(src, image.Rectangle{at, at.Add(b.Size())}, img, b.Min, draw.Src)
	out := image.NewNRGBA(image.Rect(0, 0, side, side))
	scale := float64(own) / float64(side)
	for y := range side {
		for x := range side {
			out.SetNRGBA(x, y, average(src, float64(x)*scale, float64(y)*scale, scale))
		}
	}
	return out
}

// average is the colour of the square of src at (x0, y0), size s.
func average(src *image.NRGBA, x0, y0, s float64) color.NRGBA {
	var r, g, b, a, w float64
	for y := int(y0); float64(y) < y0+s; y++ {
		wy := min(float64(y+1), y0+s) - max(float64(y), y0)
		for x := int(x0); float64(x) < x0+s; x++ {
			wx := min(float64(x+1), x0+s) - max(float64(x), x0)
			c := src.NRGBAAt(x, y)
			k := wx * wy * float64(c.A)
			r, g, b, a, w = r+k*float64(c.R), g+k*float64(c.G), b+k*float64(c.B), a+k, w+wx*wy
		}
	}
	if a == 0 {
		return color.NRGBA{}
	}
	return color.NRGBA{R: uint8(r/a + 0.5), G: uint8(g/a + 0.5), B: uint8(b/a + 0.5), A: uint8(a/w + 0.5)}
}
