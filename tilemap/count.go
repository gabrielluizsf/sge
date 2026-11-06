package tilemap

// count adds n to the count of cells showing tiles like id: opaque, or
// with transparent pixels. Nothing is counted until the counts are first
// wanted, as telling a tile's kind means drawing it, which a map only
// ever drawn from its cells, as an isometric one is, need not.
func (m *Map) count(id, n int) {
	if !m.counting {
		return
	}
	switch {
	case id == Empty:
	case m.Tiles.isOpaque(id):
		m.opaques += n
	default:
		m.glassy += n
	}
}

// tally counts the cells' tiles the first time the counts are wanted; they
// are kept up to date after.
func (m *Map) tally() {
	if m.counting {
		return
	}
	m.counting = true
	for _, id := range m.painted {
		m.count(id, 1)
	}
}
