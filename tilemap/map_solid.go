package tilemap

import "github.com/gabrielluizsf/sge/collider"

// Blocks reports whether c touches a solid tile of the map.
func (m *Map) Blocks(c collider.Collider) bool {
	return (&solidTiles{m}).Overlaps(c)
}

// Stride is half a tile, as the map's scale draws it.
func (m *Map) Stride() float64 {
	return min(float64(m.Tiles.Width)*m.Transform.Scale.X, float64(m.Tiles.Height)*m.Transform.Scale.Y) / 2
}
