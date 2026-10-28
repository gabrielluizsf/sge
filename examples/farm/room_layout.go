package main

import "github.com/gabrielluizsf/sge/random"

// roomLayout is a room's floor and things as text, a character a cell:
// the walls, a window in the back one, invisible walls round the open
// sides, the doormat, a rug, and the furniture: a bed, a shelf, a stove,
// a table and a chair, a plant and a chest, each somewhere of its own as
// seed falls.
func roomLayout(seed uint32) (ground, things string) {
	src := random.New(seed)
	g := make([][]byte, roomRows)
	t := make([][]byte, roomRows)
	for r := range roomRows {
		g[r], t[r] = make([]byte, roomCols), make([]byte, roomCols)
		for c := range roomCols {
			g[r][c], t[r][c] = '.', ' '
			switch {
			case r == roomRows-1 || c == roomCols-1:
				g[r][c], t[r][c] = ' ', 'x'
			case r == 0 && c == 0:
				t[r][c] = 'C'
			case r == 0:
				t[r][c] = 'n'
			case c == 0:
				t[r][c] = 'w'
			}
		}
	}
	g[roomRows-1][matCol], t[roomRows-1][matCol] = 'm', ' '
	t[0][2+src.Intn(4)] = 'W'
	// The furniture, each on a free cell of the places it may stand.
	put := func(ch byte, places [][2]int) {
		for range 8 {
			p := places[src.Intn(len(places))]
			if t[p[1]][p[0]] == ' ' {
				t[p[1]][p[0]] = ch
				return
			}
		}
	}
	put("bB"[src.Intn(2)], [][2]int{{1, 1}, {6, 1}})
	put('s', [][2]int{{2, 1}, {4, 1}, {5, 1}})
	put('o', [][2]int{{1, 3}, {6, 3}, {1, 4}})
	put('t', [][2]int{{4, 3}, {2, 3}})
	put('c', [][2]int{{4, 4}, {2, 4}, {5, 3}})
	put('p', [][2]int{{6, 5}, {1, 5}, {6, 4}})
	put('h', [][2]int{{1, 2}, {6, 2}, {5, 5}})
	g[2+src.Intn(2)][2+src.Intn(3)] = 'r'
	join := func(rows [][]byte) string {
		s := ""
		for _, r := range rows {
			s += string(r) + "\n"
		}
		return s
	}
	return join(g), join(t)
}
