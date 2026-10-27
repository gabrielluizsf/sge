package art

// box is a block standing on the cell: from u0, v0 to u1, v1 across it
// and from h0 up to h1, its top and the faces it turns to the screen's
// left and right drawn in those colours, outlined.
type box struct {
	u0, v0, u1, v1, h0, h1 float64
	top, left, right       rune
}

// block draws b: its faces, then its edges in the outline's colour.
func (r *raster) block(b box) {
	near, side := b.v1, b.u1
	r.poly(b.left, point{b.u0, near, b.h0}, point{b.u1, near, b.h0}, point{b.u1, near, b.h1}, point{b.u0, near, b.h1})
	r.poly(b.right, point{side, b.v0, b.h0}, point{side, b.v1, b.h0}, point{side, b.v1, b.h1}, point{side, b.v0, b.h1})
	r.poly(b.top, point{b.u0, b.v0, b.h1}, point{b.u1, b.v0, b.h1}, point{b.u1, b.v1, b.h1}, point{b.u0, b.v1, b.h1})
	for _, e := range [][2]point{
		{{b.u0, b.v0, b.h1}, {b.u1, b.v0, b.h1}}, {{b.u1, b.v0, b.h1}, {b.u1, b.v1, b.h1}},
		{{b.u1, b.v1, b.h1}, {b.u0, b.v1, b.h1}}, {{b.u0, b.v1, b.h1}, {b.u0, b.v0, b.h1}},
		{{b.u0, near, b.h0}, {b.u1, near, b.h0}}, {{side, b.v0, b.h0}, {side, b.v1, b.h0}},
		{{b.u0, near, b.h0}, {b.u0, near, b.h1}}, {{b.u1, b.v1, b.h0}, {b.u1, b.v1, b.h1}},
		{{side, b.v0, b.h0}, {side, b.v0, b.h1}},
	} {
		r.line('k', e[0], e[1])
	}
}
