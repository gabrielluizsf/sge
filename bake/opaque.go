package bake

import "image/color"

// Opaque makes s's picture wholly opaque where it is mostly covered and
// clear elsewhere, so tiles laid side by side, as the ground's are, meet
// without a gap or a seam of half-clear edge.
func Opaque(s Sprite) {
	b := s.Image.Bounds()
	for y := range b.Dy() {
		for x := range b.Dx() {
			if c := s.Image.NRGBAAt(x, y); c.A >= 96 {
				s.Image.SetNRGBA(x, y, color.NRGBA{c.R, c.G, c.B, 255})
			} else {
				s.Image.SetNRGBA(x, y, color.NRGBA{})
			}
		}
	}
}
