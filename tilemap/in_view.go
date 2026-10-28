package tilemap

import "github.com/gabrielluizsf/sge/vector"

// inView reports whether any of the map can be seen: it is drawn in the
// game's view, or the view is not known yet.
func (m *Map) inView() bool {
	g := m.Game()
	if g == nil {
		return true
	}
	w, h := g.View()
	if w == 0 || h == 0 {
		return true
	}
	if g.Projection != nil && !m.Screen {
		c0, r0, c1, r1 := m.isoCells()
		return c0 <= c1 && r0 <= r1
	}
	b := m.cellBox(0, 0)
	b.Width *= float64(m.Cols)
	b.Height *= float64(m.Rows)
	if !m.Screen {
		p := g.ToScreen(vector.Vector3{X: b.X, Y: b.Y})
		b.X, b.Y = p.X, p.Y
	}
	return b.X < float64(w) && b.Y < float64(h) && b.X+b.Width > 0 && b.Y+b.Height > 0
}
