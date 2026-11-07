package farm

import "github.com/gabrielluizsf/sge/i18n"

// editor is the arranging of furniture: a cursor on a cell, the piece
// being put down on it, if any, and the item of the bag to put down next.
// In the player's house any piece may be placed, moved, turned and taken
// back into the bag; out of doors only chests may.
type editor struct {
	col, row int
	// hold is the piece being put down, turned as it will stand; lifted
	// says it came off the floor, from origin, not out of the bag.
	hold   *piece
	lifted bool
	origin piece
	// item is the piece of the bag put down next.
	item string
	// pen is the pen of animals being moved, and penCol and penRow where
	// on it the cursor took it.
	pen            *pen
	penCol, penRow int
	// mouseX and mouseY are where the pointer last was, once mouseKnown.
	mouseX, mouseY int
	mouseKnown     bool
}

// editing reports whether the player is arranging furniture.
func (f *farm) editing() bool { return f.edit != nil }

// pieceKinds are the pieces that may be placed where the player is: any in
// their own house, only chests out of doors, none in others' houses.
func (f *farm) pieceKinds() []string {
	switch {
	case f.room == nil:
		return []string{"chest"}
	case f.room.own:
		return furnishings
	}
	return nil
}

// startEdit begins arranging, with item chosen if it is a piece of the
// bag that may be placed here, and reports whether it did: not in a
// house that is not the player's, nor while seated or in the water.
func (f *farm) startEdit(item string) bool {
	if len(f.pieceKinds()) == 0 || f.seat != nil || f.hands.swims || f.edit != nil {
		return false
	}
	col, row := f.cellOf(feet(f.player))
	d := f.hands.facing
	e := &editor{col: col + int(d.X), row: row + int(d.Y)}
	f.edit = e
	f.clampCursor()
	f.chooseItem(item)
	f.sound("ui_open", nil)
	return true
}

// chooseItem picks the item to put down: item if it can be, else the first
// piece of the bag that can.
func (f *farm) chooseItem(item string) {
	if f.canChoose(item) {
		f.edit.item = item
		return
	}
	f.edit.item = ""
	for _, n := range f.pieceKinds() {
		if f.bag.Count(n) > 0 {
			f.edit.item = n
			return
		}
	}
}

// canChoose reports whether item is a piece of the bag that may be put
// down here.
func (f *farm) canChoose(item string) bool {
	for _, n := range f.pieceKinds() {
		if n == item {
			return f.bag.Count(n) > 0
		}
	}
	return false
}

// stopEdit ends the arranging, putting back what was lifted.
func (f *farm) stopEdit() {
	if f.edit == nil {
		return
	}
	f.dropHold()
	f.edit = nil
	f.sound("ui_close", nil)
}

// dropHold lets go of the piece held: back where it was lifted from, or
// back into the bag, where it never left.
func (f *farm) dropHold() {
	e := f.edit
	e.pen = nil
	if e.hold == nil {
		return
	}
	if e.lifted {
		f.setPiece(e.origin)
	}
	e.hold, e.lifted = nil, false
}

// ghost is the piece as it would stand on the cursor, if one is held or
// one of the bag is chosen, and whether it may stand there.
func (f *farm) ghost() (p piece, ok, shown bool) {
	e := f.edit
	if e == nil {
		return piece{}, false, false
	}
	switch {
	case e.hold != nil:
		p = *e.hold
	case e.item != "":
		p = newPiece(e.item, 0, 0, 0)
	default:
		return piece{}, false, false
	}
	p.col, p.row = e.col, e.row
	return p, f.canStand(p), true
}

// canStand reports whether p may stand where it is.
func (f *farm) canStand(p piece) bool {
	if f.room != nil {
		return f.canPlaceHome(p, f.home)
	}
	return p.name == "chest" && f.canPlaceYard(p)
}

// noFit tells the player the piece does not go there.
func (f *farm) noFit() {
	f.notify(i18n.T("It does not fit there"))
	f.sound("deny", nil)
}
