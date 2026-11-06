package farm

import "github.com/gabrielluizsf/sge/vector"

// unstickRange is how far, in cells, a player found standing in a fence
// is moved to the nearest free cell.
const unstickRange = 8

// freed is p, or, where p is in a cell one cannot stand in (a fence that
// was built where a saved game left the player), the middle of the nearest
// cell one can, looking ring by ring.
func (f *farm) freed(p vector.Vector3) vector.Vector3 {
	col, row := f.cellOf(p)
	if f.open(col, row) {
		return p
	}
	for d := 1; d <= unstickRange; d++ {
		for dr := -d; dr <= d; dr++ {
			for dc := -d; dc <= d; dc++ {
				if max(abs(dc), abs(dr)) == d && f.open(col+dc, row+dr) {
					c := f.ground.CellCenter(col+dc, row+dr)
					return vector.Vector3{X: c.X, Y: c.Y, Z: p.Z}
				}
			}
		}
	}
	return p
}
