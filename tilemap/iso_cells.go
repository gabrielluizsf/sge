package tilemap

import (
	"math"

	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/vector"
)

// isoCells are the cells of a projected map that may show in the view:
// those whose ground it sees, and those nearer whose tall tiles reach up
// into it.
func (m *Map) isoCells() cellRange {
	v := m.Game().VisibleWorld()
	if v.Width <= 0 || v.Height <= 0 {
		return m.whole() // the view is not known yet
	}
	return m.isoBox(v)
}

// isoBox are the cells whose ground is in the world box v, and those
// nearer whose tall tiles reach into it.
func (m *Map) isoBox(v collider.Box) cellRange {
	g := m.Game()
	// A cell further in, by a column or a row, stands half a diamond
	// lower: a tile reaching up tall cells high shows from that far in.
	o := m.cellBox(0, 0)
	diamond := g.ToScreen(vector.Vector3{X: o.X + o.Width, Y: o.Y + o.Height}).Y - g.ToScreen(vector.Vector3{X: o.X, Y: o.Y}).Y
	tall := 0
	if diamond > 0 {
		tall = int(math.Ceil(float64(m.Tiles.Height) / (diamond / 2)))
	}
	c0, r0, _ := m.CellAt(v.X, v.Y)
	c1, r1, _ := m.CellAt(v.X+v.Width, v.Y+v.Height)
	return cellRange{c0: max(c0, 0), r0: max(r0, 0), c1: min(c1+tall, m.Cols-1), r1: min(r1+tall, m.Rows-1)}.checked()
}
