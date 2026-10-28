package tilemap

import "github.com/gabrielluizsf/sge/collider"

// cellBox is the world box cell col, row covers.
func (m *Map) cellBox(col, row int) collider.Box {
	p, s := m.Transform.Position, m.Transform.Scale
	cw, ch := m.Tiles.Cell()
	w, h := float64(cw)*s.X, float64(ch)*s.Y
	return collider.Box{X: p.X + float64(col)*w, Y: p.Y + float64(row)*h, Width: w, Height: h}
}
