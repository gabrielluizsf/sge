package farm

import (
	"strings"

	"github.com/gabrielluizsf/sge/tilemap"
)

// bareRoom is a room's things as text with no furniture: the walls, the
// window and the open sides only.
func bareRoom(things string) string {
	rows := strings.Split(strings.TrimRight(things, "\n"), "\n")
	for r, row := range rows {
		b := []byte(row)
		for c := range b {
			if _, furniture := pieceNames[b[c]]; furniture || b[c] == 'x' && r < roomRows-1 && c < roomCols-1 {
				b[c] = ' '
			}
		}
		rows[r] = string(b)
	}
	return strings.Join(rows, "\n") + "\n"
}

// layPieces sets the tiles of the room's pieces on its things, a tile for
// each cell they cover.
func (r *room) layPieces(a Art, pieces []piece) {
	for _, p := range pieces {
		s := p.cells(a)
		for j := range s.rows {
			for i := range s.cols {
				r.things.Set(s.c0+i, s.r0+j, pieceTileID(p, i, j))
			}
		}
	}
}

// clearPieces takes every piece off the room's floor.
func (r *room) clearPieces() {
	for row := 1; row < roomRows-1; row++ {
		for col := 1; col < roomCols-1; col++ {
			r.things.Set(col, row, tilemap.Empty)
		}
	}
}

// relay lays the player's house again as its pieces stand.
func (f *farm) relay() {
	f.room.clearPieces()
	f.room.layPieces(f.art, f.home)
}

// pieceAt is the index in the house's pieces of the piece on cell col,
// row, or -1.
func (f *farm) pieceAt(col, row int) int {
	for i, p := range f.home {
		if p.cells(f.art).has(col, row) {
			return i
		}
	}
	return -1
}
