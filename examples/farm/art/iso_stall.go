package art

// Stall is the shop's stall seen from the corner, two cells along a row:
// a counter with the produce on it, under a striped awning on four
// poles. It is an object, not a tile: StallPivot is the point of it that
// stands on the middle of its two cells.
func Stall() Grid {
	r := &raster{g: blank(3*Size, 3*Size+Size/2), ox: Size, oy: 2 * Size}
	pole := func(u, v float64) {
		r.block(box{u0: u, v0: v, u1: u + 1.5, v1: v + 1.5, h0: 0, h1: 26, top: 'o', left: 'q', right: 'o'})
	}
	pole(2, 3)
	pole(28.5, 3)
	r.block(box{u0: 2, v0: 5, u1: 30, v1: 14, h0: 0, h1: 10, top: 'q', left: 'q', right: 'o'})
	for i, c := range []rune{'n', 'r', 'y', 'n', 'r', 'y', 'n'} {
		u := 5 + float64(i)*3.5
		r.block(box{u0: u, v0: 8, u1: u + 2, v1: 10, h0: 10, h1: 12, top: c, left: c, right: c})
	}
	pole(2, 13)
	pole(28.5, 13)
	// The awning: stripes along the row, red and white.
	r.block(box{u0: 0, v0: 1, u1: 32, v1: 16, h0: 26, h1: 28, top: 'w', left: 'R', right: 'R'})
	for u := 0.0; u < 32; u += 4 {
		r.poly('r', point{u, 1, 28}, point{u + 2, 1, 28}, point{u + 2, 16, 28}, point{u, 16, 28})
	}
	return r.g
}

// StallPivot is the point of Stall's picture, as fractions of its width
// and height, that stands on the middle of its two cells.
var StallPivot = [2]float64{(Size + Size - Size/2) / (3.0 * Size), (2*Size + (Size+Size/2)/2) / (3.0*Size + Size/2)}
