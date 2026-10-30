package farm

import (
	"github.com/gabrielluizsf/sge/town"
)

// townRegion is the side of the regions of the world, each with a town or
// none.
const townRegion = 40

// newTowns is where the towns are: one a short walk east of the farm
// along its path, one west, and others about the world, none on the farm.
// The path's row, 4, runs on as each one's main street.
func newTowns(kinds int) *town.Layout {
	return &town.Layout{Seed: seed + 7, Region: townRegion, Chance: 0.45, Kinds: kinds,
		Fixed: map[town.Cell]town.Cell{{Col: 1, Row: 0}: {Col: 58, Row: pathRow}, {Col: -1, Row: 0}: {Col: -22, Row: pathRow}},
		Skip:  func(rx, ry int) bool { return rx == 0 && ry == 0 }}
}

// pathRow is the row of the farm's path, east to west.
const pathRow = 4

// highway reports whether cell col, row is on the road joining the farm's
// path to the two towns beside it, within margin rows of it.
func highway(col, row, margin int) bool {
	east, west := 58-town.Reach, -22+town.Reach
	return abs(row-pathRow) <= margin && (col >= 30 && col < east || col > west && col < 0)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
