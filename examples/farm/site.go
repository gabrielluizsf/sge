package farm

import (
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// treeRoom is how many cells around its buildings the farm keeps clear of
// trees, whose crowns would hide them.
const treeRoom = 3

// site is where the farm's own buildings stand: the stall, where one
// buys and sells; the player's house; its doorstep, which leads inside;
// and where the player starts.
type site struct {
	stall, home rect
	door, spawn town.Cell
	// room is how many cells around its buildings and its path the farm
	// keeps clear of trees.
	room int
}

// siteFor is the farm's site for the footprints a gives its buildings.
// The stall stands on the farm's path as ever, apart from the house. A
// house of two by two cells is north of the path, as the stylized look
// has it, the player starting in front of the stall; a larger house, with
// no room there, stands in the open ground south of the path, its
// doorstep in front, where the player starts.
func siteFor(a Art) site {
	s := site{stall: stallRect(a, farmStall), room: treeRoom}
	if a.Reach() > 0 { // trees grown large would hide the house and the way
		s.room = 10
	}
	cols, rows := a.Footprint("house0")
	s.home = rect{16, 1, cols, rows}
	s.spawn = town.Cell{Col: stallCol, Row: stallRow + 1}
	if cols > 2 || rows > 2 {
		s.home = rect{16, 6, cols, rows}
		s.spawn = town.Cell{}
	}
	s.door = town.Cell{Col: s.home.col + (cols-1)/2, Row: s.home.row + rows}
	if s.spawn == (town.Cell{}) {
		s.spawn = s.door
	}
	return s
}

// apply puts the site's buildings on the things of l: the cells they
// stand on solid, and nothing else in them, nor a tree beside them.
func (s site) apply(l layers) {
	block := tileID["block"]
	for row := range l.things.Rows {
		for col := range l.things.Cols {
			if l.things.At(col, row) == block {
				l.things.Set(col, row, tilemap.Empty)
			}
		}
	}
	for _, r := range []rect{s.stall, s.home} {
		for row := r.row - s.room; row <= r.row+r.rows+s.room; row++ {
			for col := r.col - s.room; col <= r.col+r.cols+s.room; col++ {
				if isTree(l.things.At(col, row)) || r.has(col, row) {
					l.things.Set(col, row, tilemap.Empty)
				}
			}
		}
		for row := r.row; row < r.row+r.rows; row++ {
			for col := r.col; col < r.col+r.cols; col++ {
				l.things.Set(col, row, block)
			}
		}
	}
	if s.room > treeRoom { // only a look of large trees needs the way kept clear
		s.clearAlongPath(l)
	}
}

// clearAlongPath takes the trees within room cells of the path, which
// leads to the towns.
func (s site) clearAlongPath(l layers) {
	path := tileID["path"]
	for row := range l.things.Rows {
		for col := range l.things.Cols {
			if l.ground.At(col, row) != path {
				continue
			}
			for r := row - s.room; r <= row+s.room; r++ {
				for c := col - s.room; c <= col+s.room; c++ {
					if isTree(l.things.At(c, r)) {
						l.things.Set(c, r, tilemap.Empty)
					}
				}
			}
		}
	}
}
