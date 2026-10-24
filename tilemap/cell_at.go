package tilemap

import "math"

// CellAt is the cell the world point x, y falls in, and whether it is on
// the map.
func (m *Map) CellAt(x, y float64) (col, row int, ok bool) {
	p, s := m.Transform.Position, m.Transform.Scale
	col = int(math.Floor((x - p.X) / (float64(m.Tiles.Width) * s.X)))
	row = int(math.Floor((y - p.Y) / (float64(m.Tiles.Height) * s.Y)))
	return col, row, col >= 0 && row >= 0 && col < m.Cols && row < m.Rows
}
