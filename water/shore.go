package water

import (
	"image"
	"image/color"
	"math"
)

// Shore is a tile of water to make: the mask of which neighbours are
// water too, the water's Surface, a Normals map for the frame, Span of
// cells across, the Ground tile the land and the bed are, and how far
// into its loop the frame is, from 0 to 1, for the foam. Col and Row are
// the cell's place; only whether each is odd matters, as for ground
// patches: the waves and the shoreline's wiggle are those of the world,
// repeating every Span cells, so no two cells beside each other alike.
type Shore struct {
	Mask     int
	Col, Row int
	Surface  *Surface
	Normals  *image.NRGBA
	Ground   *image.NRGBA
	// Bed, when set, is the picture of the bed under the water, seen
	// through it, in place of the ground's colour.
	Bed   *image.NRGBA
	Phase float64
}

// Span is how many cells across the world's waves and shoreline repeat.
const Span = 2

// spot is a pixel of a tile: where it lies in the cell, u along a row and
// v down a column, each 0 to 1, where it lies in the world's repeat, wu
// and wv each 0 to 1, and the ground's colour there, in linear light.
type spot struct {
	u, v, wu, wv float64
	ground       [3]float64
}

// The shore's shape, in cells: the water's edge lies about edge from the
// land, wavering by reach, and depth is full rampLen from it.
const (
	edge    = 0.23
	reach   = 0.07
	rampLen = 0.65
)

// Tile is the tile of a shore: a cell of water in the shape of Ground,
// the diamond, with land where the mask says a neighbour is not water, its
// edge a soft, wavering shoreline of damp ground and foam, the water
// deepening away from it. Every tile of a mask meets every other, for
// the shore depends on where a pixel lies in the cell and on the land
// around, never on the tile.
func Tile(s Shore) *image.NRGBA {
	b := s.Ground.Bounds()
	out := image.NewNRGBA(b)
	land := landOf(s.Mask)
	w, h := float64(b.Dx()), float64(b.Dy())
	for y := range b.Dy() {
		for x := range b.Dx() {
			g := s.Ground.NRGBAAt(b.Min.X+x, b.Min.Y+y)
			if g.A == 0 {
				continue
			}
			a, c := (float64(x)+0.5-w/2)/(w/2), (float64(y)+0.5)/(h/2)
			p := spot{u: math.Max(0, math.Min((a+c)/2, 1)), v: math.Max(0, math.Min((c-a)/2, 1)),
				ground: [3]float64{linear(float64(g.R) / 255), linear(float64(g.G) / 255), linear(float64(g.B) / 255)}}
			p.wu, p.wv = s.world(p.u, p.v)
			rgb := s.pixel(land, p)
			out.SetNRGBA(b.Min.X+x, b.Min.Y+y, color.NRGBA{uint8(math.Round(srgb(rgb[0]) * 255)), uint8(math.Round(srgb(rgb[1]) * 255)), uint8(math.Round(srgb(rgb[2]) * 255)), g.A})
		}
	}
	return out
}

// pixel is the colour, in linear light, of the shore at p: ground on the
// land side of the water's edge, water on the other, blended across a
// sliver of it.
func (s *Shore) pixel(land [][2]float64, p spot) [3]float64 {
	wob := noise([2]float64{p.wu, p.wv}, 8, 11)*2 - 1 + 0.4*(noise([2]float64{p.wu, p.wv}, 16, 12)*2-1)
	x := distance(land, p.u, p.v) - (edge + reach*wob)
	cover := smooth(-0.012, 0.012, x)
	var dry, wet [3]float64
	if cover < 1 {
		damp := 0.6 + 0.4*smooth(0, 0.1, -x)
		dry = [3]float64{p.ground[0] * damp, p.ground[1] * damp, p.ground[2] * damp}
	}
	if cover > 0 {
		wet = s.water(p, x)
	}
	return mix3(dry, wet, cover)
}

// water is the water at p, x cells beyond the water's edge: deeper as x
// grows, calmer where shallow, and foam at the edge whose reach swells and
// ebbs with the loop.
func (s *Shore) water(p spot, x float64) [3]float64 {
	sf := s.Surface
	depth := sf.MaxDepth * smooth(0, rampLen, x)
	n := normalAt(s.Normals, p.wu, p.wv)
	calm := 0.35 + 0.65*smooth(0, 0.3, depth)
	n = unit([3]float64{n[0] * calm, n[1] * calm, 1 - calm + n[2]*calm})
	swell := 0.5 + 0.5*math.Sin(2*math.Pi*(s.Phase+noise([2]float64{p.wu, p.wv}, 12, 13)))
	reachNow := 0.04 + 0.07*swell
	lace := 0.45 + 0.55*noise([2]float64{p.wu, p.wv}, 28, 14)
	foam := math.Min(1, (1-smooth(0, reachNow, x))*lace*1.25)
	bed := [3]float64{p.ground[0] * 0.75, p.ground[1] * 0.75, p.ground[2] * 0.75}
	if s.Bed != nil {
		bed = bedAt(s.Bed, p.wu, p.wv)
	}
	return sf.Shade(Look{Normal: n, Depth: depth, Bed: bed, Foam: foam})
}

// parity is 0 for an even cell place and 1 for an odd, negative or not.
func parity(n int) int { return (n%Span + Span) % Span }

// world is where (u, v) of the cell lies in the world's repeat, each 0 to
// 1: neighbouring cells meet there, a cell's far edge the next one's near.
func (s *Shore) world(u, v float64) (float64, float64) {
	return (float64(parity(s.Col)) + u) / Span, (float64(parity(s.Row)) + v) / Span
}
