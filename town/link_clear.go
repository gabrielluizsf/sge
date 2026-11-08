package town

// keepOff is how many cells a road keeps from a town it does not go to
// (the houses, yards and what the wild keeps clear of them), and from the
// places it is told to avoid.
const keepOff = 3

// clear reports whether the road p, from town t to n, keeps off every
// other town and every Keep area.
func (l *Layout) clear(p []Rect, t, n *Town) bool {
	for _, r := range p {
		for _, k := range l.Keep {
			if r.Within(k, 0) {
				return false
			}
		}
		rx0, rx1 := floorDiv(r.C0-keepOff-l.Reach(), l.Region), floorDiv(r.C1+keepOff+l.Reach(), l.Region)
		ry0, ry1 := floorDiv(r.R0-keepOff-l.Depth(), l.Region), floorDiv(r.R1+keepOff+l.Depth(), l.Region)
		for ry := ry0; ry <= ry1; ry++ {
			for rx := rx0; rx <= rx1; rx++ {
				m := l.Of(rx, ry)
				if m != nil && m != t && m != n && r.Within(m.box(), keepOff) {
					return false
				}
			}
		}
	}
	return true
}

// box is the cells of the town: its street and its houses' reach.
func (t *Town) box() Rect {
	return Rect{t.Center.Col - t.Reach, t.Center.Row - t.Depth, t.Center.Col + t.Reach, t.Center.Row + t.Depth}
}
