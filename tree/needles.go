package tree

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/random"
)

// needles paints a spray of pine needles on a 128 pixel card: a twig up
// its middle, from the bottom, and needles along it, leaning forward,
// shorter toward the tip. Clear between them.
func needles(seed uint32) *image.NRGBA {
	const size = 128
	src := random.New(seed)
	p := newPen(size, size)
	dark, light := color.NRGBA{30, 72, 42, 255}, color.NRGBA{92, 142, 72, 255}
	p.line(stroke{[2]float64{64, size}, [2]float64{64, 6}, 2.4, color.NRGBA{84, 62, 40, 255}})
	const pairs = 58
	for i := range pairs {
		t := float64(i) / pairs
		y := size - 4 - t*118
		for _, side := range []float64{-1, 1} {
			a := 0.4 + 0.45*src.Float()
			l := 48*(1-0.62*t) + 8*src.Float()
			p.line(stroke{[2]float64{64, y}, [2]float64{64 + side*math.Sin(a)*l, y - math.Cos(a)*l}, 1.8, blend(dark, light, src.Float())})
		}
	}
	return p.img
}
