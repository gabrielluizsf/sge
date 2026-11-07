package farm

// bedsIn are the beds laid in a room's things, in reading order: the
// double bed 'D' and the single beds 'b' and 'B', each by the cells it
// covers.
func bedsIn(things string, size func(piece byte) (cols, rows int)) []bed {
	var out []bed
	row, col := 0, 0
	for i := 0; i < len(things); i++ {
		switch ch := things[i]; ch {
		case '\n':
			row, col = row+1, 0
			continue
		case 'b', 'B', 'D':
			cols, rows := size(ch)
			out = append(out, bed{span: span{c0: col, r0: row, cols: max(cols, 1), rows: max(rows, 1)}, double: ch == 'D'})
		}
		col++
	}
	return out
}
