package tree

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/random"
)

// blossoms paints a cluster of flowers on a 128 pixel card: many small
// five-petalled blooms packed in a round bunch over a few bare twigs,
// clear around it. The petals are in the colours petal and heart say.
func blossoms(seed uint32, petal, heart color.NRGBA) *image.NRGBA {
	const size = 128
	src := random.New(seed)
	p := newPen(size, size)
	twig := color.NRGBA{82, 62, 42, 255}
	for range 4 {
		p.line(stroke{[2]float64{64, size}, [2]float64{30 + 68*src.Float(), 26 + 40*src.Float()}, 2.2, twig})
	}
	for range 46 {
		a, d := 2*math.Pi*src.Float(), 44*math.Sqrt(src.Float())
		at := [2]float64{64 + math.Cos(a)*d, 62 + math.Sin(a)*d*0.9}
		p.bloom(flower{at, 5 + 4*src.Float(), blend(petal, color.NRGBA{petal.R, petal.G / 2, petal.B, 255}, 0.4*src.Float()), heart})
	}
	return p.img
}

// flower is a bloom to paint: its middle, its radius across, and the
// colours of its petals and its heart.
type flower struct {
	at           [2]float64
	r            float64
	petal, heart color.NRGBA
}

// bloom paints flower f: five round petals about a heart.
func (p pen) bloom(f flower) {
	for k := range 5 {
		a := 2 * math.Pi * float64(k) / 5
		p.dot(disc{[2]float64{f.at[0] + math.Cos(a)*f.r*0.55, f.at[1] + math.Sin(a)*f.r*0.55}, f.r * 0.55, f.petal})
	}
	p.dot(disc{f.at, f.r * 0.25, f.heart})
}
