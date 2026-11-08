package tree

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/random"
)

// frond paints a palm's pinnate frond on a 96 by 256 card: a rachis up
// its middle, bare at its foot, and narrow leaflets in pairs along it,
// leaning forward, long in the middle and short at the tip.
func frond(seed uint32) *image.NRGBA {
	const w, h = 96, 256
	src := random.New(seed)
	p := newPen(w, h)
	dark, light := color.NRGBA{44, 94, 38, 255}, color.NRGBA{112, 158, 62, 255}
	p.line(stroke{[2]float64{48, h}, [2]float64{48, 3}, 3, color.NRGBA{120, 130, 70, 255}})
	for y := h - 26.0; y > 6; y -= 2.6 {
		t := (h - y) / h
		l := 8 + 38*math.Sin(math.Pi*math.Min(1, 0.15+0.85*math.Pow(1-t, 0.7)))
		for _, side := range []float64{-1, 1} {
			a := 0.95 + 0.2*src.Float()
			p.line(stroke{[2]float64{48, y}, [2]float64{48 + side*math.Sin(a)*l, y - math.Cos(a)*l}, 2.4, blend(dark, light, src.Float())})
		}
	}
	return p.img
}
