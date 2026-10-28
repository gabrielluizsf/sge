package tilemap

// At is the tile in cell col, row, or Empty for a cell off the map.
func (m *Map) At(col, row int) int {
	if col < 0 || row < 0 || col >= m.Cols || row >= m.Rows {
		return Empty
	}
	return m.cells[row*m.Cols+col]
}
