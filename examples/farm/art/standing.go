package art

// standing is the upright sprite g standing in the middle of a cell's
// diamond, on a shadow: a tree, a crop. Its transparent bottom rows are
// cut first, so what it stands on is its lowest pixel.
func standing(g Grid, shadow bool) Grid {
	g = trimmed(g)
	out := blank(IsoW, IsoH)
	if shadow {
		r := &raster{g: out, ox: Size, oy: isoTop}
		r.ellipse('h', 7, 3.5)
		out = r.g
	}
	foot := isoTop + Size/2 + 1
	for y, row := range g {
		for x, c := range []rune(row) {
			if c != '.' {
				out = out.set(x+(IsoW-len([]rune(row)))/2, foot-len(g)+y, c)
			}
		}
	}
	return out
}

// trimmed is g without its transparent rows at the bottom.
func trimmed(g Grid) Grid {
	n := len(g)
	for n > 0 && isBlankRow(g[n-1]) {
		n--
	}
	return g[:n]
}

func isBlankRow(row string) bool {
	for _, c := range row {
		if c != '.' {
			return false
		}
	}
	return true
}

// ellipse fills an ellipse rx by ry pixels on the ground, around the
// cell's middle, with c.
func (r *raster) ellipse(c rune, rx, ry float64) {
	cx, cy := r.at(point{Size / 2, Size / 2, 0})
	for y := range r.g {
		for x := range len([]rune(r.g[y])) {
			dx, dy := (float64(x)+0.5-cx)/rx, (float64(y)+0.5-cy)/ry
			if dx*dx+dy*dy <= 1 {
				r.g = r.g.set(x, y, c)
			}
		}
	}
}
