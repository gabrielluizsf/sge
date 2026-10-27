package art

// A tile's picture: a cell's diamond, IsoW wide and half that tall, at
// the bottom of a picture IsoH tall, so what stands on the cell reaches
// up out of it.
const (
	IsoW = 2 * Size
	IsoH = 3 * Size
	// isoTop is the row of the diamond's top corner.
	isoTop = IsoH - Size
)

// blank is a transparent grid w by h.
func blank(w, h int) Grid {
	row := make([]rune, w)
	for i := range row {
		row[i] = '.'
	}
	g := make(Grid, h)
	for i := range g {
		g[i] = string(row)
	}
	return g
}

// isoGround is the square tile g laid flat on its cell's diamond: each
// pixel of the diamond shows the pixel of g under it.
func isoGround(g Grid) Grid {
	out := blank(IsoW, IsoH)
	for y := isoTop; y < IsoH; y++ {
		for x := range IsoW {
			// The point of the cell, u across and v down, drawn here.
			sx, sy := float64(x)+0.5-Size, float64(y-isoTop)+0.5
			u, v := sy+sx/2, sy-sx/2
			if u >= 0 && v >= 0 && u < Size && v < Size {
				out = out.set(x, y, []rune(g[int(v)])[int(u)])
			}
		}
	}
	return out
}
