package tilemap

import (
	"math"

	"github.com/gabrielluizsf/sge/vector"
)

// isoAt is where, in a projected world, the picture of cell col, row's
// tile goes on the screen: its bottom middle on the bottom corner of the
// cell's diamond, so a tall tile stands up from its cell.
func (m *Map) isoAt(col, row int) (x, y int) {
	b := m.cellBox(col, row)
	bottom := m.Game().ToScreen(vector.Vector3{X: b.X + b.Width, Y: b.Y + b.Height})
	return int(math.Floor(bottom.X+0.5)) - m.Tiles.Width/2, int(math.Floor(bottom.Y+0.5)) - m.Tiles.Height
}

// isoDepth is how near the screen cell col, row stands: its middle's
// screen y, as objects' feet are.
func (m *Map) isoDepth(col, row int) float64 {
	return m.Game().ToScreen(m.CellCenter(col, row)).Y
}
