package tilemap

import (
	"math"

	"github.com/gabrielluizsf/sge/vector"
)

// isoCells are the columns and rows, from c0, r0 to c1, r1, of the cells
// of a projected map that may show in the view: those whose ground it
// sees, and those nearer whose tall tiles reach up into it.
func (m *Map) isoCells() (c0, r0, c1, r1 int) {
	g := m.Game()
	v := g.VisibleWorld()
	if v.Width <= 0 || v.Height <= 0 {
		return 0, 0, m.Cols - 1, m.Rows - 1 // the view is not known yet
	}
	// A cell further in, by a column or a row, stands half a diamond
	// lower: a tile reaching up tall cells high shows from that far in.
	o := m.cellBox(0, 0)
	diamond := g.ToScreen(vector.Vector3{X: o.X + o.Width, Y: o.Y + o.Height}).Y - g.ToScreen(vector.Vector3{X: o.X, Y: o.Y}).Y
	tall := 0
	if diamond > 0 {
		tall = int(math.Ceil(float64(m.Tiles.Height) / (diamond / 2)))
	}
	c0, r0, _ = m.CellAt(v.X, v.Y)
	c1, r1, _ = m.CellAt(v.X+v.Width, v.Y+v.Height)
	return max(c0, 0), max(r0, 0), min(c1+tall, m.Cols-1), min(r1+tall, m.Rows-1)
}
