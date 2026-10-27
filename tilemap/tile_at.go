package tilemap

// TileAt is the tile under the world point x, y, or Empty off the map.
func (m *Map) TileAt(x, y float64) int {
	col, row, _ := m.CellAt(x, y)
	return m.At(col, row)
}
