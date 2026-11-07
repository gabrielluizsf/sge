package farm

// fit lays the piece ch, cols by rows cells, on the free cells of t over
// place, the first way round it that fits against the walls: drawn from
// its back cell, the one of the least column and row, the others blocks.
// Behind each of its cells is the piece or a wall, so no one stands
// behind any part of it and whoever stands beside it is drawn over it.
// It reports whether it found room.
func fit(t [][]byte, place [2]int, cols, rows int, ch byte) bool {
	cols, rows = max(cols, 1), max(rows, 1)
	for dy := range rows {
		for dx := range cols {
			at := span{c0: place[0] - dx, r0: place[1] - dy, cols: cols, rows: rows}
			if !free(t, at) || !backed(t, at) {
				continue
			}
			for r := at.r0; r < at.r0+rows; r++ {
				for c := at.c0; c < at.c0+cols; c++ {
					t[r][c] = 'x'
				}
			}
			t[at.r0][at.c0] = ch
			return true
		}
	}
	return false
}

// span is the cells a piece covers: cols by rows from c0, r0.
type span struct{ c0, r0, cols, rows int }

// has reports whether cell c, r is one of s's.
func (s span) has(c, r int) bool {
	return c >= s.c0 && c < s.c0+s.cols && r >= s.r0 && r < s.r0+s.rows
}

// free reports whether the cells of s are all in the room and empty.
func free(t [][]byte, s span) bool {
	for r := s.r0; r < s.r0+s.rows; r++ {
		for c := s.c0; c < s.c0+s.cols; c++ {
			if r < 0 || c < 0 || r >= len(t) || c >= len(t[r]) || t[r][c] != ' ' {
				return false
			}
		}
	}
	return true
}

// backed reports whether behind each cell of s, a column or a row back,
// is s itself or a wall; one cell is always backed, being drawn from its
// own cell.
func backed(t [][]byte, s span) bool {
	if s.cols*s.rows == 1 {
		return true
	}
	for r := s.r0; r < s.r0+s.rows; r++ {
		for c := s.c0; c < s.c0+s.cols; c++ {
			for _, b := range [][2]int{{c - 1, r}, {c, r - 1}} {
				if !s.has(b[0], b[1]) && !wallAt(t, b[0], b[1]) {
					return false
				}
			}
		}
	}
	return true
}

// wallAt reports whether cell c, r of t is a back wall, or off the room.
func wallAt(t [][]byte, c, r int) bool {
	if r < 0 || c < 0 || r >= len(t) || c >= len(t[r]) {
		return true
	}
	switch t[r][c] {
	case 'C', 'n', 'W', 'w':
		return true
	}
	return false
}
