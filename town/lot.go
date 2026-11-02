package town

// size is the footprint of a house of the kind: Width in Col, Depth in
// Row.
func (l *Layout) size(kind int) Cell {
	if l.Footprint == nil {
		return Cell{Col: defaultW, Row: defaultW}
	}
	c := l.Footprint(kind)
	return Cell{Col: max(c.Col, 1), Row: max(c.Row, 1)}
}

// lot is the largest footprint of any kind of house.
func (l *Layout) lot() Cell {
	var lot Cell
	for k := range max(l.Kinds, 1) {
		c := l.size(k)
		lot = Cell{Col: max(lot.Col, c.Col), Row: max(lot.Row, c.Row)}
	}
	return lot
}

// pitch is how far apart, along the street, the houses of a side start.
func (l *Layout) pitch() int { return l.lot().Col + max(l.Gap, 1) }

// reachOf is how far the street runs each way when n houses line each
// side of it.
func (l *Layout) reachOf(n int) int { return l.lead() + (n-1)*l.pitch() + l.lot().Col }

// lead is how far from the plaza's middle the first houses stand: the
// west ones end that many cells from it, the east ones start one cell
// nearer, the stall, anvil and forge being on the west side of it.
func (l *Layout) lead() int {
	return max(l.Lead, defaultLead)
}

// Lots is how many houses line a side of the street each way from the
// plaza, one or two, as many as fit in a region with the town's ends and
// the margin the wild keeps clear of them.
func (l *Layout) Lots() int {
	if 2*(l.reachOf(2)+3) < l.Region {
		return 2
	}
	return 1
}

// Reach is how far the main streets run each way from a plaza's middle.
func (l *Layout) Reach() int { return l.reachOf(l.Lots()) }

// Depth is how many cells the houses and yards of a town reach from the
// street's middle.
func (l *Layout) Depth() int { return l.lot().Row + 1 }

// slot is where a house may stand along a side of the street, -1 north
// and 1 south: its lot's first column from the plaza's middle, and whether
// it is west of it.
type slot struct {
	start, side int
	west        bool
}

// slots are the lots of a side of the street, west to east.
func (l *Layout) slots() []slot {
	var out []slot
	n, w, pitch := l.Lots(), l.lot().Col, l.pitch()
	for k := n - 1; k >= 0; k-- {
		out = append(out, slot{start: -l.lead() - w + 1 - k*pitch, west: true})
	}
	for k := range n {
		out = append(out, slot{start: l.lead() - 1 + k*pitch})
	}
	return out
}
