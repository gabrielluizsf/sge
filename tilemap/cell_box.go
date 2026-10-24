package tilemap

import "github.com/gabrielluizsf/sge/collider"

// cellBox is the world box cell col, row covers.
func (m *Map) cellBox(col, row int) collider.Box {
	p, s := m.Transform.Position, m.Transform.Scale
	w, h := float64(m.Tiles.Width)*s.X, float64(m.Tiles.Height)*s.Y
	return collider.Box{X: p.X + float64(col)*w, Y: p.Y + float64(row)*h, Width: w, Height: h}
}
