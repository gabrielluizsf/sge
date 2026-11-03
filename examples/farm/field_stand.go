package farm

import "github.com/gabrielluizsf/sge/pathfind"

// standBy is a cell to stand on beside cell col, row, and the way to face
// it from there: an open one not on soil if there is one, never the one
// the player is on.
func (f *farm) standBy(col, row int) (pathfind.Cell, string, bool) {
	mine := pathfind.Cell{}
	mine.Col, mine.Row = f.cellOf(feet(f.player))
	for _, bare := range []bool{true, false} {
		for _, b := range banks {
			c := pathfind.Cell{Col: col - b.dc, Row: row - b.dr}
			if f.open(c.Col, c.Row) && c != mine && f.onFarm(c.Col, c.Row) && !f.isWater(c.Col, c.Row) && (!bare || f.plotAt(c.Col, c.Row) == nil) {
				return c, b.face, true
			}
		}
	}
	return pathfind.Cell{}, "", false
}
