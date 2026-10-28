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
	m.cells[row*m.Cols+col] = id
	m.paint(col, row)
}
