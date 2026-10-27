package art

// isoFence is a run of fence along a cell's row, seen from the corner:
// two rails from the cell's one side to the other, so a row of them
// joins up, and a post in the middle.
func isoFence() Grid {
	r := &raster{g: blank(IsoW, IsoH), ox: Size, oy: isoTop}
	for _, h := range []float64{3, 7} {
		r.line('o', point{0, Size / 2, h}, point{Size, Size / 2, h})
		r.line('q', point{0, Size / 2, h + 1}, point{Size, Size / 2, h + 1})
	}
	r.block(box{u0: 7, v0: 7, u1: 9, v1: 9, h0: 0, h1: 10, top: 'q', left: 'q', right: 'o'})
	return r.g
}
