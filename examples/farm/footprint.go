package farm

import (
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// rect is the cells a building stands on: cols by rows of them from col,
// row. They are solid, so nobody walks into the building or over it.
type rect struct{ col, row, cols, rows int }

// has reports whether cell c, r is one of the rect's.
func (r rect) has(c, row int) bool {
	return c >= r.col && c < r.col+r.cols && row >= r.row && row < r.row+r.rows
}

// middle is the world point in the middle of the rect: where a building's
// picture stands.
func (r rect) middle() vector.Vector3 {
	return vector.Vector3{X: (float64(r.col) + float64(r.cols)/2) * tile, Y: (float64(r.row) + float64(r.rows)/2) * tile}
}

// houseRect is the cells of a town's house.
func houseRect(h town.House) rect { return rect{h.Col, h.Row, h.Width, h.Depth} }

// stallRect is the cells of the stall at the cell at, a's footprint for it:
// across the middle of the cell, its counter's row the last of them, where
// one stands to buy.
func stallRect(a Art, at town.Cell) rect {
	cols, rows := a.Footprint("stall")
	return rect{at.Col - cols/2, at.Row - rows + 1, cols, rows}
}

// near reports whether cell c, row is one of the rect's or beside it: no
// water laps at a house.
func (r rect) near(c, row int) bool {
	return c >= r.col-1 && c <= r.col+r.cols && row >= r.row-1 && row <= r.row+r.rows
}

// wellRect is the cells of a town's well, a's footprint for it, none when
// the look has no well.
func wellRect(a Art, t *town.Town) rect {
	cols, rows := a.Footprint("well")
	return rect{t.Well.Col, t.Well.Row, cols, rows}
}
