package farm

// The cells of the house a piece may stand on: inside the walls.
func insideRoom(col, row int) bool {
	return col >= 1 && col <= roomCols-2 && row >= 1 && row <= roomRows-2
}

// entry is the cell inside the door mat, where the player comes in: it is
// kept free so the way out is never blocked.
var entry = [2]int{matCol, roomRows - 2}

// canPlaceHome reports whether piece p may stand in the house with the
// pieces others: the house still fits, and p is not on the player or on
// anyone else who is there.
func (f *farm) canPlaceHome(p piece, others []piece) bool {
	s := p.cells(f.art)
	pc, pr := f.cellOf(feet(f.player))
	for r := s.r0; r < s.r0+s.rows; r++ {
		for c := s.c0; c < s.c0+s.cols; c++ {
			if c == pc && r == pr || f.residentAt(c, r) {
				return false
			}
		}
	}
	return houseFits(f.art, append(append([]piece(nil), others...), p))
}

// houseFits reports whether pieces make a house: every piece on the floor
// inside the walls, none on another, the door's cell free, and every free
// cell of the floor still reachable from it, so nothing is walled in and
// nobody is trapped.
func houseFits(a Art, pieces []piece) bool {
	taken := map[[2]int]bool{}
	for _, p := range pieces {
		s := p.cells(a)
		for r := s.r0; r < s.r0+s.rows; r++ {
			for c := s.c0; c < s.c0+s.cols; c++ {
				at := [2]int{c, r}
				if !insideRoom(c, r) || taken[at] || at == entry {
					return false
				}
				taken[at] = true
			}
		}
	}
	return connected(taken)
}

// connected reports whether the free cells of the floor, those inside the
// room not in taken, are all reachable from the entry on foot.
func connected(taken map[[2]int]bool) bool {
	seen := map[[2]int]bool{entry: true}
	queue := [][2]int{entry}
	for len(queue) > 0 {
		at := queue[0]
		queue = queue[1:]
		for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			n := [2]int{at[0] + d[0], at[1] + d[1]}
			if insideRoom(n[0], n[1]) && !taken[n] && !seen[n] {
				seen[n] = true
				queue = append(queue, n)
			}
		}
	}
	for r := 1; r <= roomRows-2; r++ {
		for c := 1; c <= roomCols-2; c++ {
			if at := [2]int{c, r}; !taken[at] && !seen[at] {
				return false
			}
		}
	}
	return true
}

// residentAt reports whether someone who lives there stands on cell col, row.
func (f *farm) residentAt(col, row int) bool {
	for _, r := range f.residents {
		if c, rw := f.cellOf(r.obj.Transform.Position); c == col && rw == row {
			return true
		}
	}
	return false
}

// editReach is how many cells from the player, across or down, a chest
// may be put out of doors.
const editReach = 4

// canPlaceYard reports whether chest p may stand out of doors: on dry,
// free ground near the player, not on a plot, a door or a stall's counter, a
// tree, another chest, or on the player.
func (f *farm) canPlaceYard(p piece) bool {
	pc, pr := f.cellOf(feet(f.player))
	switch {
	case abs(p.col-pc) > editReach || abs(p.row-pr) > editReach:
		return false
	case p.col == pc && p.row == pr:
		return false
	case !f.open(p.col, p.row) || f.plotAt(p.col, p.row) != nil:
		return false
	}
	if _, door := f.doorAt(p.col, p.row); door {
		return false
	}
	_, stall := f.stallAt(p.col, p.row)
	return !stall
}
