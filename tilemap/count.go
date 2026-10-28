package tilemap

// count adds n to the count of cells showing tiles like id: opaque, or
// with transparent pixels.
func (m *Map) count(id, n int) {
	switch {
	case id == Empty:
	case m.Tiles.isOpaque(id):
		m.opaques += n
	default:
		m.glassy += n
	}
}
