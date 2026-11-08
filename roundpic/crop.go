// Package roundpic cuts a picture to a circle with a thin ring, for a
// portrait in a bar. It makes the picture once; the caller keeps it, so
// nothing is masked again every frame.
package roundpic

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// ringWidth is the ring's thickness, in pixels of the circle.
const ringWidth = 2.0

// Crop is the top square of src, as wide as src is, scaled to d by d
// pixels and cut to a circle with a ring of colour ring. A nil src is a
// neutral circle.
func Crop(src *canvas.Canvas, d int, ring canvas.Color) *canvas.Canvas {
	if src == nil || src.Width == 0 || src.Height == 0 {
		return Neutral(d, 0xFF9AA0A6, ring)
	}
	side := min(src.Width, src.Height)
	x0 := (src.Width - side) / 2
	return disc(d, ring, func(x, y int) canvas.Color {
		return src.At(x0+x*side/d, y*side/d)
	})
}

// Neutral is a plain circle of colour fill with a ring, for who has no
// picture yet.
func Neutral(d int, fill, ring canvas.Color) *canvas.Canvas {
	return disc(d, ring, func(int, int) canvas.Color { return fill })
}

// disc is a d by d picture of what pixel gives, inside a circle, its
// edge softened by a pixel and its rim the ring.
func disc(d int, ring canvas.Color, pixel func(x, y int) canvas.Color) *canvas.Canvas {
	out, err := canvas.NewCanvas(d, d)
	if err != nil {
		return nil
	}
	r := float64(d) / 2
	for y := 0; y < d; y++ {
		for x := 0; x < d; x++ {
			dist := math.Hypot(float64(x)+0.5-r, float64(y)+0.5-r)
			cover := min(1, max(0, r-dist)) // 1 inside, falling to 0 over the last pixel
			if cover == 0 {
				continue
			}
			c := pixel(x, y)
			if dist > r-ringWidth {
				c = ring
			}
			out.Put(x, y, canvas.Fade(c, int(float64(c.A())*cover+0.5)))
		}
	}
	return out
}
