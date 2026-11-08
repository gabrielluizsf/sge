package tree

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/random"
)

// skin paints a cactus's skin on a 128 pixel tile, laid to wrap once
// round a rib: green, darker in the groove at its middle, a column of
// pale areoles with their spines along the ridge at its edge.
func skin(seed uint32) *image.NRGBA {
	const size = 128
	src := random.New(seed)
	p := newPen(size, size)
	green := color.NRGBA{70, 120, 84, 255}
	for y := range size {
		for x := range size {
			ridge := 0.78 + 0.22*math.Abs(math.Cos(math.Pi*float64(x)/size))
			streak := 0.94 + 0.06*math.Sin(float64(x)*0.9+float64(y)*0.05)
			c := scaleBy(green, ridge*streak)
			p.img.SetNRGBA(x, y, color.NRGBA{c.R, c.G, c.B, 255})
		}
	}
	for _, x := range []float64{0, size} {
		for y := 16.0; y < size; y += 32 {
			areole(p, [2]float64{x, y}, src)
		}
	}
	return p.img
}

// areole paints a cactus's areole at at: a dark patch and a fan of pale
// spines.
func areole(p pen, at [2]float64, src *random.Source) {
	p.dot(disc{at, 4.5, color.NRGBA{48, 78, 52, 255}})
	for range 9 {
		a := 2 * math.Pi * src.Float()
		l := 7 + 7*src.Float()
		p.line(stroke{at, [2]float64{at[0] + math.Cos(a)*l, at[1] + math.Sin(a)*l}, 1.2, color.NRGBA{236, 228, 200, 255}})
	}
	p.dot(disc{at, 2.2, color.NRGBA{214, 202, 168, 255}})
}

// scaleBy is c times k.
func scaleBy(c color.NRGBA, k float64) color.NRGBA {
	s := func(v uint8) uint8 { return uint8(math.Max(0, math.Min(255, float64(v)*k))) }
	return color.NRGBA{s(c.R), s(c.G), s(c.B), 255}
}
