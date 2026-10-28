package tilemap

// SolidAt reports whether the world point x, y is on a solid tile.
func (m *Map) SolidAt(x, y float64) bool {
	col, row, ok := m.CellAt(x, y)
	return ok && m.Tiles.IsSolid(m.At(col, row))
}
