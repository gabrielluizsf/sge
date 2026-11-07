package farm

// bedsFor lays a house's beds in the room t: the couple's double bed in the
// corner of the back walls, and the child's single bed, of the letter
// single, beside it or below it, where each has a wall or a piece behind
// every cell.
func bedsFor(t [][]byte, size func(piece byte) (cols, rows int), single byte) {
	dc, dr := size('D')
	dc, dr = max(dc, 1), max(dr, 1)
	if !fitBed(t, span{c0: 1, r0: 1, cols: dc, rows: dr}, 'D') {
		return
	}
	cols, rows := size(single)
	cols, rows = max(cols, 1), max(rows, 1)
	for _, at := range [][2]int{{1 + dc, 1}, {1, 1 + dr}, {roomCols - 1 - cols, 1}, {1, roomRows - 1 - rows}} {
		if fitBed(t, span{c0: at[0], r0: at[1], cols: cols, rows: rows}, single) {
			return
		}
	}
}

// fitBed puts a bed of the letter ch on the cells of s, if they are free
// and each has a wall or a piece behind it, as its back cell, the others
// blocks.
func fitBed(t [][]byte, s span, ch byte) bool {
	if !free(t, s) {
		return false
	}
	for r := s.r0; r < s.r0+s.rows && s.cols*s.rows > 1; r++ { // one cell is drawn from its own
		for c := s.c0; c < s.c0+s.cols; c++ {
			for _, b := range [][2]int{{c - 1, r}, {c, r - 1}} {
				if !s.has(b[0], b[1]) && (b[0] < 0 || b[1] < 0 || t[b[1]][b[0]] == ' ') {
					return false
				}
			}
		}
	}
	for r := s.r0; r < s.r0+s.rows; r++ {
		for c := s.c0; c < s.c0+s.cols; c++ {
			t[r][c] = 'x'
		}
	}
	t[s.r0][s.c0] = ch
	return true
}
