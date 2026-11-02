package farm

import "github.com/gabrielluizsf/sge/pathfind"

// spring is a place to fill a can: the dry cell to stand on, and the way
// to face the water from it.
type spring struct {
	stand pathfind.Cell
	face  string
}

// springRange is how many cells from where they start villagers look for
// water.
const springRange = 40

// bank is a way to look from a cell: a step, and the way it faces. The
// farmer who looks up faces row - 1.
type bank struct {
	dc, dr int
	face   string
}

var banks = []bank{{1, 0, "right"}, {-1, 0, "left"}, {0, 1, "down"}, {0, -1, "up"}}

// springNear is the nearest water to cell from that can be reached from
// a dry cell beside it, looking out in rings. It is false where there is
// none within springRange.
func (f *farm) springNear(from pathfind.Cell) (spring, bool) {
	return f.springWhere(from, f.isWater)
}

// springOnFarm is the nearest water of the farm's own, in the same row as
// cell from, as springNear finds it: one the demo's bot, which walks
// across and then up or down, can reach and come back from.
func (f *farm) springOnFarm(from pathfind.Cell) (spring, bool) {
	return f.springWhere(from, func(col, row int) bool { return f.onFarm(col, row) && row == from.Row && f.isWater(col, row) })
}

// springWhere is the nearest cell for which water is true, as springNear
// finds it.
func (f *farm) springWhere(from pathfind.Cell, water func(col, row int) bool) (spring, bool) {
	for r := 1; r <= springRange; r++ {
		for dc := -r; dc <= r; dc++ {
			for dr := -r; dr <= r; dr++ {
				if max(abs(dc), abs(dr)) != r || !water(from.Col+dc, from.Row+dr) {
					continue
				}
				if s, ok := f.bankOf(pathfind.Cell{Col: from.Col + dc, Row: from.Row + dr}); ok {
					return s, true
				}
			}
		}
	}
	return spring{}, false
}

// bankOf is a dry cell beside water cell w to stand on, facing it.
func (f *farm) bankOf(w pathfind.Cell) (spring, bool) {
	for _, b := range banks {
		col, row := w.Col-b.dc, w.Row-b.dr
		if f.open(col, row) && !f.isWater(col, row) {
			return spring{stand: pathfind.Cell{Col: col, Row: row}, face: b.face}, true
		}
	}
	return spring{}, false
}
