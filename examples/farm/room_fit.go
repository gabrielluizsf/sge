package farm

// fit lays the piece ch, cols by rows cells, on the free cells of t over
// place, the first way round it that fits: drawn from its front cell,
// the one of the most column and row, the others blocks. It reports
// whether it found room.
func fit(t [][]byte, place [2]int, cols, rows int, ch byte) bool {
	cols, rows = max(cols, 1), max(rows, 1)
	for dy := range rows {
		for dx := range cols {
			c0, r0 := place[0]-dx, place[1]-dy
			if !free(t, c0, r0, cols, rows) {
				continue
			}
			for r := r0; r < r0+rows; r++ {
				for c := c0; c < c0+cols; c++ {
					t[r][c] = 'x'
				}
			}
			t[r0+rows-1][c0+cols-1] = ch
			return true
		}
	}
	return false
}

// free reports whether the cols by rows cells of t from c0, r0 are all
// in the room and empty.
func free(t [][]byte, c0, r0, cols, rows int) bool {
	for r := r0; r < r0+rows; r++ {
		for c := c0; c < c0+cols; c++ {
			if r < 0 || c < 0 || r >= len(t) || c >= len(t[r]) || t[r][c] != ' ' {
				return false
			}
		}
	}
	return true
}
