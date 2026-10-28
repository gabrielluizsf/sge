package tilemap

import "github.com/gabrielluizsf/sge/vector"

// CellCenter is the world point at the middle of cell col, row.
func (m *Map) CellCenter(col, row int) vector.Vector3 {
	b := m.cellBox(col, row)
	return vector.Vector3{X: b.X + b.Width/2, Y: b.Y + b.Height/2}
}
