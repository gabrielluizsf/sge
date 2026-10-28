package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// spansOf is, for a picture whose every row is one run of wholly opaque
// pixels with nothing else showing, as a ground tile's diamond is, each
// row's run from x0 up to x1; nil for any other picture, which is
// blended.
func spansOf(img *canvas.Canvas) [][2]int {
	spans := make([][2]int, img.Height)
	for y := range img.Height {
		x0, x1 := -1, -1
		for x := range img.Width {
			switch a := img.At(x, y).A(); {
			case a == 255 && x1 < 0:
				if x0 < 0 {
					x0 = x
				}
			case a == 255:
				return nil // a second run
			case a != 0:
				return nil // a pixel to blend
			case x0 >= 0 && x1 < 0:
				x1 = x
			}
		}
		if x0 >= 0 && x1 < 0 {
			x1 = img.Width
		}
		spans[y] = [2]int{max(x0, 0), max(x1, 0)}
	}
	return spans
}

// copySpans copies t's runs onto cv at x, y, clipped to it: no pixel of
// them needs blending.
func copySpans(cv *canvas.Canvas, t trimmed, x, y int) {
	for r, s := range t.spans {
		cy := y + r
		if cy < 0 || cy >= cv.Height || s[0] >= s[1] {
			continue
		}
		a, b := max(s[0], -x), min(s[1], cv.Width-x)
		if a >= b {
			continue
		}
		src := t.img.Pixels[r*t.img.Stride+a : r*t.img.Stride+b]
		copy(cv.Pixels[cy*cv.Stride+x+a:], src)
	}
}
