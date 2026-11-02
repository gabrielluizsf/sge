package tilemap

import (
	"math"

	"github.com/gabrielluizsf/sge/vector"
)

// isoAt is where, in a projected world, the picture of cell col, row's
// tile goes on the screen: its bottom middle on the bottom corner of the
// cell's diamond, so a tall tile stands up from its cell. A picture of
// any size, the tile set's or its own: a tree wider and
// taller than a cell is placed by its own bottom middle.
func (m *Map) isoAt(col, row int, t trimmed) (x, y int) {
	b := m.cellBox(col, row)
	bottom := m.Game().ToScreen(vector.Vector3{X: b.X + b.Width, Y: b.Y + b.Height})
	return int(math.Floor(bottom.X+0.5)) - t.w/2, int(math.Floor(bottom.Y+0.5)) - t.h
}

// isoDepth is how near the screen cell col, row stands: its middle's
// screen y, as objects' feet are.
func (m *Map) isoDepth(col, row int) float64 {
	return m.Game().ToScreen(m.CellCenter(col, row)).Y
}
