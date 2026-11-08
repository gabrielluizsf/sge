package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// place is a point of a canvas.
type place struct{ x, y int }

// copySpans copies t's runs onto cv at p, clipped to cv's clip: no pixel
// of them needs blending.
func copySpans(cv *canvas.Canvas, t trimmed, p place) {
	clip := cv.Clip
	for r, s := range t.spans {
		cy := p.y + r
		if cy < clip.Y || cy >= clip.Y+clip.Height || s[0] >= s[1] {
			continue
		}
		a, b := max(s[0], clip.X-p.x), min(s[1], clip.X+clip.Width-p.x)
		if a >= b {
			continue
		}
		src := t.img.Pixels[r*t.img.Stride+a : r*t.img.Stride+b]
		copy(cv.Pixels[cy*cv.Stride+p.x+a:], src)
	}
}
