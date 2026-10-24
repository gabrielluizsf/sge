package tilemap

// Fill puts tile id in every cell of the w by h cells from col, row.
func (m *Map) Fill(col, row, w, h, id int) {
	for r := row; r < row+h; r++ {
		for c := col; c < col+w; c++ {
			m.Set(c, r, id)
		}
	}
}
