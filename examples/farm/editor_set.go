package farm

// setPiece puts p on the floor of the house, or out of doors, and shows it.
func (f *farm) setPiece(p piece) {
	if f.room != nil {
		f.home = append(f.home, p)
		f.relay()
		return
	}
	f.putChest(p)
}

// pieceUnder is the piece on cell col, row of the house, or the chest out
// of doors, and whether there is one.
func (f *farm) pieceUnder(col, row int) (piece, bool) {
	if f.room != nil {
		if i := f.pieceAt(col, row); i >= 0 {
			return f.home[i], true
		}
		return piece{}, false
	}
	if y := f.chestAt(col, row); y != nil {
		return y.p, true
	}
	return piece{}, false
}

// removePiece takes the piece on cell col, row off the floor.
func (f *farm) removePiece(col, row int) {
	if f.room != nil {
		if i := f.pieceAt(col, row); i >= 0 {
			f.home = append(f.home[:i], f.home[i+1:]...)
			f.relay()
		}
		return
	}
	if y := f.chestAt(col, row); y != nil {
		f.takeChest(y)
	}
}

// clampCursor keeps the cursor on the floor of the house, or within the
// player's reach out of doors.
func (f *farm) clampCursor() {
	e := f.edit
	if f.room != nil {
		e.col, e.row = max(1, min(roomCols-2, e.col)), max(1, min(roomRows-2, e.row))
		return
	}
	pc, pr := f.cellOf(feet(f.player))
	e.col, e.row = max(pc-editReach, min(pc+editReach, e.col)), max(pr-editReach, min(pr+editReach, e.row))
}

// lifting is the piece being moved, as it stood before it was lifted, if
// one is, so a save made while moving it keeps it where it was.
func (f *farm) lifting() (piece, bool) {
	if e := f.edit; e != nil && e.hold != nil && e.lifted {
		return e.origin, true
	}
	return piece{}, false
}
