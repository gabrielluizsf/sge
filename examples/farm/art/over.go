package art

// over is top drawn over base: base shows where top is '.'.
func over(base, top Grid) Grid {
	out := base.copyOf()
	for y, row := range top {
		for x, c := range []rune(row) {
			if c != '.' {
				out = out.set(x, y, c)
			}
		}
	}
	return out
}

// fill is a grid of one letter.
func fill(c rune) Grid {
	row := ""
	for range Size {
		row += string(c)
	}
	g := make(Grid, Size)
	for i := range g {
		g[i] = row
	}
	return g
}
