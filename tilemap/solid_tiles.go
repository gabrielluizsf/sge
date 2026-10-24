package tilemap

import (
	"math"

	"github.com/gabrielluizsf/sge/collider"
)

// solidTiles is a map's collider: its solid tiles, wherever the map is.
// It works out where from the map itself, so moving it changes nothing.
type solidTiles struct{ m *Map }

// Bounds is the whole map.
func (s *solidTiles) Bounds() collider.Box {
	b := s.m.cellBox(0, 0)
	return collider.Box{X: b.X, Y: b.Y, Width: b.Width * float64(s.m.Cols), Height: b.Height * float64(s.m.Rows)}
}

// Moved is the same tiles: they are where the map is.
func (s *solidTiles) Moved(dx, dy float64) collider.Collider { return s }

// Overlaps reports whether other touches a solid tile, testing only the
// cells under other's bounds.
func (s *solidTiles) Overlaps(other collider.Collider) bool {
	m := s.m
	b := other.Bounds()
	c0, r0, _ := m.CellAt(b.X, b.Y)
	c1, r1, _ := m.CellAt(math.Nextafter(b.X+b.Width, b.X), math.Nextafter(b.Y+b.Height, b.Y))
	for r := max(r0, 0); r <= min(r1, m.Rows-1); r++ {
		for c := max(c0, 0); c <= min(c1, m.Cols-1); c++ {
			if m.Tiles.IsSolid(m.At(c, r)) && collider.Collides(m.cellBox(c, r), other) {
				return true
			}
		}
	}
	return false
}
