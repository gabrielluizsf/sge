package flora

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/random"
)

// imageOf is a texture a plant is painted with.
type imageOf = image.Image

// Leaf is the shape and colour of a kind of leaf as its texture paints
// it: its half width at its widest as a share of the card it is painted
// on, near 0.5 filling it (the leaf's own width is the card's, set where
// it grows), where along it that widest is, how serrated
// its edge is, how many pairs of veins it has, and its colours, from the
// base to the tip and its veins, in sRGB.
type Leaf struct {
	Width, Widest float64
	Teeth         float64
	Veins         int
	Base, Tip     color.NRGBA
	Vein          color.NRGBA
}

// texture paints a leaf on a size by size texture, its stalk at the
// bottom middle and its tip at the top, transparent around it, its edge
// soft over a pixel; seed varies its mottling.
func (l Leaf) texture(size int, seed uint32) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	src := random.New(seed)
	blotch := [8]float64{}
	for i := range blotch {
		blotch[i] = src.Float()
	}
	for y := range size {
		t := 1 - (float64(y)+0.5)/float64(size) // 0 at the stalk, 1 at the tip
		half := l.halfWidth(t)
		if l.Teeth > 0 {
			half *= 1 - l.Teeth*0.5*(1+math.Sin(t*float64(max(l.Veins, 1))*4*math.Pi))*0.3
		}
		for x := range size {
			u := ((float64(x)+0.5)/float64(size) - 0.5) * 2 // -1 to 1 across
			edge := (half - math.Abs(u)) * float64(size) / 2
			if edge <= -0.5 {
				continue
			}
			alpha := math.Min(1, edge+0.5)
			c := mix(l.Base, l.Tip, t)
			// Lighter along the midrib and the veins, darker at the rim.
			vein := veinAt(u, t, l.Veins)
			c = mix(c, l.Vein, vein*0.55)
			rim := 1 - math.Min(1, math.Abs(u)/math.Max(half, 1e-6))
			shade := 0.82 + 0.18*rim + 0.06*(blotch[int(t*7.99)]-0.5)
			c = scaleColour(c, shade)
			img.SetNRGBA(x, y, color.NRGBA{c.R, c.G, c.B, uint8(alpha * 255)})
		}
	}
	return img
}

// halfWidth is the leaf's half width at t of its length, 0 to 1 across:
// widest at Widest, narrowing to a point at the tip and to the stalk.
func (l Leaf) halfWidth(t float64) float64 {
	if t < l.Widest {
		return l.Width * math.Sin(t/l.Widest*math.Pi/2)
	}
	return l.Width * math.Cos((t-l.Widest)/(1-l.Widest)*math.Pi/2)
}

// veinAt is how much of a vein the point u across and t along lies on:
// the midrib, and pairs of veins running out and up from it.
func veinAt(u, t float64, pairs int) float64 {
	v := math.Max(0, 1-math.Abs(u)*30)
	if pairs > 0 {
		// A vein leaves the midrib at each step and climbs as it goes out.
		along := t - math.Abs(u)*0.35
		k := along * float64(pairs)
		d := math.Abs(k - math.Round(k))
		v = math.Max(v, math.Max(0, 1-d*14)*0.6*(1-math.Abs(u)))
	}
	return v
}

func mix(a, b color.NRGBA, t float64) color.NRGBA {
	l := func(x, y uint8) uint8 { return uint8(float64(x) + (float64(y)-float64(x))*t) }
	return color.NRGBA{l(a.R, b.R), l(a.G, b.G), l(a.B, b.B), 255}
}

func scaleColour(c color.NRGBA, k float64) color.NRGBA {
	s := func(x uint8) uint8 { return uint8(math.Max(0, math.Min(255, float64(x)*k))) }
	return color.NRGBA{s(c.R), s(c.G), s(c.B), 255}
}
