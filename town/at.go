package town

// At is the town of the region cell col, row lies in, or nil.
func (l *Layout) At(col, row int) *Town {
	return l.Of(floorDiv(col, l.Region), floorDiv(row, l.Region))
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
