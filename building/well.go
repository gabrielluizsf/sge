package building

import (
	"image/color"
	"math"
)

// Well is a village well, 1.7 m across: a round wall of stone, coped,
// water dark down in it; two posts holding a windlass, its crank, rope
// and bucket; and a little tiled roof over it.
var Well = Building{Name: "well", Width: 1.7, Depth: 1.7, build: well}

// ring is a round wall about the origin, from radius in to out, between
// heights lo and hi, of sides sides.
type ring struct {
	in, out, lo, hi float64
	sides           int
	paint           paint
}

// wallRound adds ring r, a segment of it to each side.
func (b *builder) wallRound(r ring) {
	for i := range r.sides {
		a0 := 2 * math.Pi * float64(i) / float64(r.sides)
		a1 := 2 * math.Pi * float64(i+1) / float64(r.sides)
		var c [8]vec
		for k, y := range []float64{r.lo, r.hi} {
			c[4*k+0] = vec{r.in * math.Cos(a0), y, r.in * math.Sin(a0)}
			c[4*k+1] = vec{r.out * math.Cos(a0), y, r.out * math.Sin(a0)}
			c[4*k+2] = vec{r.out * math.Cos(a1), y, r.out * math.Sin(a1)}
			c[4*k+3] = vec{r.in * math.Cos(a1), y, r.in * math.Sin(a1)}
		}
		b.convex(hexa(c, r.paint))
	}
}

func well(b *builder) {
	stone, wood := paint{of: Stone}, paint{of: Timber}
	b.wallRound(ring{in: 0.58, out: 0.8, lo: 0, hi: 0.78, sides: 16, paint: stone})
	b.wallRound(ring{in: 0.55, out: 0.85, lo: 0.78, hi: 0.88, sides: 16, paint: stone})
	b.turn(column{a: vec{}, b: vec{0, 0.3, 0}, ra: 0.6, rb: 0.6, sides: 16, paint: paint{of: Water}})
	const post, axle = 0.72, 1.5
	for _, x := range []float64{-post, post} {
		b.block(box{lo: vec{x - 0.06, 0.88, -0.06}, hi: vec{x + 0.06, 2.15, 0.06}, paint: wood})
	}
	b.turn(column{a: vec{-post - 0.08, axle, 0}, b: vec{post + 0.08, axle, 0}, ra: 0.07, rb: 0.07, sides: 10, paint: wood})
	iron := paint{of: Iron}
	b.lay(beam{a: vec{post + 0.1, axle, 0}, b: vec{post + 0.1, axle - 0.25, 0.08}, w: 0.03, h: 0.03, paint: iron})
	b.lay(beam{a: vec{post + 0.1, axle - 0.25, 0.08}, b: vec{post + 0.24, axle - 0.25, 0.08}, w: 0.035, h: 0.035, paint: wood})
	rope := paint{of: Flat, tint: color.NRGBA{170, 150, 110, 255}}
	b.turn(column{a: vec{0, axle - 0.07, 0}, b: vec{0, 1.12, 0}, ra: 0.012, rb: 0.012, sides: 6, paint: rope})
	b.turn(column{a: vec{0, 0.86, 0}, b: vec{0, 1.12, 0}, ra: 0.12, rb: 0.15, sides: 12, paint: paint{of: Planks}})
	for _, y := range []float64{0.9, 1.06} {
		r := 0.12 + 0.03*(y-0.86)/0.26 + 0.006
		b.turn(column{a: vec{0, y, 0}, b: vec{0, y + 0.025, 0}, ra: r, rb: r + 0.003, sides: 12, paint: iron})
	}
	p := pitched{w: 2 * post, d: 1.1, eaves: 2.15, slope: math.Tan(35 * math.Pi / 180), over: 0.25}
	b.gable(cover{p, paint{of: Tiles}})
	// Its frame: a tie across each post's top, a king post up to the
	// ridge, and a plate along each eave.
	for _, x := range []float64{-post, post} {
		b.lay(beam{a: vec{x, 2.1, -0.62}, b: vec{x, 2.1, 0.62}, w: 0.1, h: 0.1, paint: wood})
		b.block(box{lo: vec{x - 0.05, 2.15, -0.05}, hi: vec{x + 0.05, p.ridge(), 0.05}, paint: wood})
	}
	for _, z := range []float64{-0.55, 0.55} {
		b.lay(beam{a: vec{-post - 0.1, 2.15, z}, b: vec{post + 0.1, 2.15, z}, w: 0.08, h: 0.08, paint: wood})
	}
}
