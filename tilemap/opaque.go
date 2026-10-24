package tilemap

// opaque reports whether every cell shows an opaque tile, so the map
// hides whatever is under it.
func (m *Map) opaque() bool {
	return m.opaques == len(m.cells)
}
