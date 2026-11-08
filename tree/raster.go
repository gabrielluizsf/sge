package tree

import (
	"image"
	"image/color"
	"math"
)

// pen paints strokes and discs onto a texture, edges soft over a pixel,
// clipped to the image.
type pen struct{ img *image.NRGBA }

// newPen is a pen on a clear w by h texture.
func newPen(w, h int) pen { return pen{image.NewNRGBA(image.Rect(0, 0, w, h))} }

// stroke is a line from a to b, w wide, of colour c, in pixels.
type stroke struct {
	a, b [2]float64
	w    float64
	c    color.NRGBA
}

// disc is a round dot at at, r in radius, of colour c, in pixels.
type disc struct {
	at [2]float64
	r  float64
	c  color.NRGBA
}

// line paints stroke s.
func (p pen) line(s stroke) {
	half := s.w / 2
	p.cover(bounds{s.a, s.b, half + 1}, func(x, y float64) float64 {
		return half - segment(s, x, y) + 0.5
	}, s.c)
}

// dot paints disc d.
func (p pen) dot(d disc) {
	p.cover(bounds{d.at, d.at, d.r + 1}, func(x, y float64) float64 {
		return d.r - math.Hypot(x-d.at[0], y-d.at[1]) + 0.5
	}, d.c)
}

// bounds are the box round two points, grown by pad.
type bounds struct {
	a, b [2]float64
	pad  float64
}

// cover paints colour c over the pixels in box, as much as inside says
// they are: 1 or more wholly, 0 or less not at all.
func (p pen) cover(box bounds, inside func(x, y float64) float64, c color.NRGBA) {
	r := p.img.Bounds()
	x0 := max(r.Min.X, int(math.Floor(min(box.a[0], box.b[0])-box.pad)))
	x1 := min(r.Max.X-1, int(math.Ceil(max(box.a[0], box.b[0])+box.pad)))
	y0 := max(r.Min.Y, int(math.Floor(min(box.a[1], box.b[1])-box.pad)))
	y1 := min(r.Max.Y-1, int(math.Ceil(max(box.a[1], box.b[1])+box.pad)))
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if k := math.Min(1, inside(float64(x)+0.5, float64(y)+0.5)); k > 0 {
				p.img.SetNRGBA(x, y, over(p.img.NRGBAAt(x, y), color.NRGBA{c.R, c.G, c.B, uint8(k * float64(c.A))}))
			}
		}
	}
}

// segment is how far the point x, y is from stroke s's line.
func segment(s stroke, x, y float64) float64 {
	dx, dy := s.b[0]-s.a[0], s.b[1]-s.a[1]
	t := 0.0
	if l := dx*dx + dy*dy; l > 0 {
		t = math.Max(0, math.Min(1, ((x-s.a[0])*dx+(y-s.a[1])*dy)/l))
	}
	return math.Hypot(x-(s.a[0]+t*dx), y-(s.a[1]+t*dy))
}

// over is top laid over under, straight alpha.
func over(under, top color.NRGBA) color.NRGBA {
	ta, ua := float64(top.A)/255, float64(under.A)/255
	a := ta + ua*(1-ta)
	if a == 0 {
		return color.NRGBA{}
	}
	mix := func(t, u uint8) uint8 { return uint8((float64(t)*ta + float64(u)*ua*(1-ta)) / a) }
	return color.NRGBA{mix(top.R, under.R), mix(top.G, under.G), mix(top.B, under.B), uint8(a * 255)}
}

// blend is a mixed with b, t of the way to b.
func blend(a, b color.NRGBA, t float64) color.NRGBA {
	m := func(x, y uint8) uint8 { return uint8(float64(x) + (float64(y)-float64(x))*t) }
	return color.NRGBA{m(a.R, b.R), m(a.G, b.G), m(a.B, b.B), 255}
}
