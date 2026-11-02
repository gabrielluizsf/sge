package tilemap

// cellRange is a block of cells, from column c0, row r0 to column c1, row
// r1, and whether it holds any.
type cellRange struct {
	c0, r0, c1, r1 int
	valid          bool
}

// checked is the block, valid when it holds a cell.
func (r cellRange) checked() cellRange {
	r.valid = r.c0 <= r.c1 && r.r0 <= r.r1
	return r
}

// whole is the block of all the cells of m.
func (m *Map) whole() cellRange {
	return cellRange{c1: m.Cols - 1, r1: m.Rows - 1}.checked()
}

// has reports whether cell col, row is in the block.
func (r cellRange) has(col, row int) bool {
	return r.valid && col >= r.c0 && col <= r.c1 && row >= r.r0 && row <= r.r1
}

// covers reports whether every cell of o is in the block.
func (r cellRange) covers(o cellRange) bool {
	return !o.valid || r.valid && o.c0 >= r.c0 && o.c1 <= r.c1 && o.r0 >= r.r0 && o.r1 <= r.r1
}
