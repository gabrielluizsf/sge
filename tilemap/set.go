package tilemap

// Set puts tile id in cell col, row, Empty clearing it, and paints the
// cell into the map's picture. A cell off the map is left alone.
func (m *Map) Set(col, row, id int) {
	if col < 0 || row < 0 || col >= m.Cols || row >= m.Rows {
		return
	}
	if id < Empty || id >= len(m.Tiles.Tiles) {
		id = Empty
	}
	i := row*m.Cols + col
	switch {
	case m.cells[i] == Empty && id != Empty:
		m.filled++
	case m.cells[i] != Empty && id == Empty:
		m.filled--
	}
	m.cells[i] = id
	m.paint(col, row)
	m.Sprite = nil
	if 2*m.filled >= len(m.cells) {
		m.Picture()
		m.Sprite = m.picture
	}
}
