package town

// margin is how many cells beyond a town's houses and street At still
// finds it from, where the wild keeps clear of it.
const margin = 3

// At is the town cell col, row belongs to: the one of the region it lies
// in or, where a town planned on a fixed middle reaches over a region's
// edge, the neighbouring region's, or nil.
func (l *Layout) At(col, row int) *Town {
	rx, ry := floorDiv(col, l.Region), floorDiv(row, l.Region)
	own := l.Of(rx, ry)
	if own != nil && own.Grounds(col, row, margin) {
		return own
	}
	dx, dy := l.edge(col, rx), l.edge(row, ry)
	for _, y := range []int{0, dy} {
		for _, x := range []int{0, dx} {
			if t := l.Of(rx+x, ry+y); t != nil && (x != 0 || y != 0) && t.Grounds(col, row, margin) {
				return t
			}
		}
	}
	return own
}

// edge is -1 when cell c, in region r, is within a town's reach of the
// region's low edge, 1 near its high edge and 0 elsewhere.
func (l *Layout) edge(c, r int) int {
	reach := max(l.Reach(), l.Depth()) + margin
	off := c - r*l.Region
	if off < reach {
		return -1
	}
	if off >= l.Region-reach {
		return 1
	}
	return 0
}

// Near is every town whose plaza's middle is within dist cells of cell
// col, row, across and down, nearest region first.
func (l *Layout) Near(col, row, dist int) []*Town {
	var near []*Town
	rx0, rx1 := floorDiv(col-dist, l.Region), floorDiv(col+dist, l.Region)
	ry0, ry1 := floorDiv(row-dist, l.Region), floorDiv(row+dist, l.Region)
	for ry := ry0; ry <= ry1; ry++ {
		for rx := rx0; rx <= rx1; rx++ {
			t := l.Of(rx, ry)
			if t != nil && abs(t.Center.Col-col) <= dist && abs(t.Center.Row-row) <= dist {
				near = append(near, t)
			}
		}
	}
	return near
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
