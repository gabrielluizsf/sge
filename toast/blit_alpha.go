package toast

import "github.com/gabrielluizsf/antui/canvas"

// BlitAlpha draws src on cv at x, y as faint as opacity says, 0 to 1,
// without making a picture of its own: every pixel's alpha is scaled as it
// is blended.
func BlitAlpha(cv *canvas.Canvas, x, y int, src *canvas.Canvas, opacity float64) {
	if src == nil || opacity <= 0 {
		return
	}
	if opacity >= 1 {
		cv.Blit(x, y, src)
		return
	}
	for sy := 0; sy < src.Height; sy++ {
		for sx := 0; sx < src.Width; sx++ {
			if c := src.At(sx, sy); c.A() > 0 {
				cv.Pixel(x+sx, y+sy, canvas.Fade(c, int(float64(c.A())*opacity+0.5)))
			}
		}
	}
}
