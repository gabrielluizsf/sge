package town

// Rect is a box of cells, both corners included.
type Rect struct{ C0, R0, C1, R1 int }

// span is the rectangle between two cells, whichever way round they are.
func span(c0, r0, c1, r1 int) Rect {
	return Rect{min(c0, c1), min(r0, r1), max(c0, c1), max(r0, r1)}
}

// Dist is how many cells cell col, row is from r, across or down, whichever
// is more: 0 on it.
func (r Rect) Dist(col, row int) int {
	return max(r.C0-col, 0, col-r.C1, r.R0-row, 0, row-r.R1)
}

// Within reports whether r and o are no more than d cells apart.
func (r Rect) Within(o Rect, d int) bool {
	return o.C0-d <= r.C1 && r.C0-d <= o.C1 && o.R0-d <= r.R1 && r.R0-d <= o.R1
}
