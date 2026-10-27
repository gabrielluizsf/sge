package tilemap

import "github.com/gabrielluizsf/sge/collider"

// Blocks reports whether c touches a solid tile of the map.
func (m *Map) Blocks(c collider.Collider) bool {
	return (&solidTiles{m}).Overlaps(c)
}

// Stride is half a tile, as the map's scale draws it.
func (m *Map) Stride() float64 {
	cw, ch := m.Tiles.Cell()
	return min(float64(cw)*m.Transform.Scale.X, float64(ch)*m.Transform.Scale.Y) / 2
}
