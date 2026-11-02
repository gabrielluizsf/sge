package farm

import (
	"github.com/gabrielluizsf/sge/town"
)

// townRegion is the side of the regions of the world, each with a town or
// none.
const townRegion = 40

// newTowns is where the towns are: one a short walk east of the farm
// along its path, one west, and others about the world, none on the farm.
// The path's row, 4, runs on as each one's main street. Their houses
// stand on the footprints a gives them, with room between for the people
// walking there to be seen.
func newTowns(a Art) *town.Layout {
	l := &town.Layout{Seed: seed + 7, Region: townRegion, Chance: 0.45, Kinds: a.Houses(), Gap: 1,
		Footprint: func(kind int) town.Cell {
			cols, rows := a.Footprint(spriteName("house", kind, true))
			return town.Cell{Col: cols, Row: rows}
		},
		Fixed: map[town.Cell]town.Cell{{Col: 1, Row: 0}: {Col: 58, Row: pathRow}, {Col: -1, Row: 0}: {Col: -22, Row: pathRow}},
		Skip:  func(rx, ry int) bool { return rx == 0 && ry == 0 }}
	if cols, _ := a.Footprint("house0"); cols > 2 { // houses with room between and around the plaza
		l.Gap, l.Lead, l.Margin = 3, 7, 3
	}
	return l
}

// pathRow is the row of the farm's path, east to west.
const pathRow = 4

// highway is the road joining the farm's path to the two towns beside it,
// whose streets run reach cells each way: it reports whether cell col,
// row is on it, within margin rows of it.
func highway(reach int) func(col, row, margin int) bool {
	east, west := 58-reach, -22+reach
	return func(col, row, margin int) bool {
		return abs(row-pathRow) <= margin && (col >= 30 && col < east || col > west && col < 0)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
