package farm

import "github.com/gabrielluizsf/sge/examples/i18n"

// editMove moves the cursor dx cells along the columns and dy along the
// rows: seen from the corner, right goes down and right, down goes down
// and left, so every cell can be reached by a key of its own.
func (f *farm) editMove(dx, dy int) {
	e := f.edit
	e.col += dx
	e.row += dy
	f.clampCursor()
}

// editAct is Enter: put down the piece held on the cursor, or lift the
// piece under it, or take one of the bag in hand.
func (f *farm) editAct() {
	e := f.edit
	if e.hold != nil {
		f.putDown()
		return
	}
	if p, ok := f.pieceUnder(e.col, e.row); ok {
		f.removePiece(p.col, p.row)
		e.hold, e.lifted, e.origin = &p, true, p
		e.col, e.row = p.col, p.row
		f.sound("ui_tick", nil)
		return
	}
	if e.item == "" || f.bag.Count(e.item) == 0 {
		f.notify(i18n.T("Nothing to place: craft furniture in the backpack"))
		f.sound("deny", nil)
		return
	}
	p := newPiece(e.item, e.col, e.row, 0)
	e.hold = &p
	f.sound("ui_tick", nil)
}

// putDown stands the piece held where the cursor is, if it fits.
func (f *farm) putDown() {
	e := f.edit
	p, ok, _ := f.ghost()
	if !ok {
		f.noFit()
		return
	}
	if !e.lifted {
		f.bag.Take(p.name, 1)
	}
	f.setPiece(p)
	e.hold, e.lifted = nil, false
	f.chooseItem(e.item)
	f.sound("ui_open", nil)
}

// editTurn is R: turn the piece held a quarter, or the piece under the
// cursor where it stands, if it still fits.
func (f *farm) editTurn() {
	e := f.edit
	if e.hold != nil {
		t := e.hold.turned()
		e.hold = &t
		f.sound("ui_tick", nil)
		return
	}
	p, ok := f.pieceUnder(e.col, e.row)
	if !ok {
		return
	}
	f.removePiece(p.col, p.row)
	t := p.turned()
	if _, fits, _ := f.ghostOf(t); fits {
		f.setPiece(t)
		f.sound("ui_tick", nil)
		return
	}
	f.setPiece(p)
	f.noFit()
}

// ghostOf is whether piece p may stand where it is.
func (f *farm) ghostOf(p piece) (piece, bool, bool) { return p, f.canStand(p), true }

// editStore is X: put the piece held, or under the cursor, back in the
// bag. A chest must be empty, and the house's last bed stays: the player
// must be able to sleep.
func (f *farm) editStore() {
	e := f.edit
	p, held := piece{}, e.hold != nil
	if held {
		p = *e.hold
	} else if under, ok := f.pieceUnder(e.col, e.row); ok {
		p = under
	} else {
		return
	}
	switch {
	case held && !e.lifted:
		e.hold = nil // it never left the bag
		return
	case !p.isEmpty():
		f.notify(i18n.T("Empty the chest first"))
		f.sound("deny", nil)
		return
	case p.name == "bed" && f.room != nil && f.beds() <= 1:
		f.notify(i18n.T("You need a bed to sleep in"))
		f.sound("deny", nil)
		return
	}
	if held {
		e.hold, e.lifted = nil, false
	} else {
		f.removePiece(p.col, p.row)
	}
	f.bag.Add(p.name, 1)
	f.sound("ui_close", nil)
}

// beds is how many beds the house has, one lifted counted.
func (f *farm) beds() int {
	n := 0
	for _, p := range f.home {
		if p.name == "bed" {
			n++
		}
	}
	if f.edit != nil && f.edit.hold != nil && f.edit.lifted && f.edit.hold.name == "bed" {
		n++
	}
	return n
}

// editCycle is Q and E: the next piece of the bag to put down, d on.
func (f *farm) editCycle(d int) {
	e := f.edit
	kinds := f.pieceKinds()
	at := 0
	for i, n := range kinds {
		if n == e.item {
			at = i
		}
	}
	for k := 1; k <= len(kinds); k++ {
		n := kinds[((at+d*k)%len(kinds)+len(kinds))%len(kinds)]
		if f.bag.Count(n) > 0 {
			e.item = n
			if e.hold != nil && !e.lifted {
				h := newPiece(n, e.col, e.row, e.hold.heading)
				e.hold = &h
			}
			f.sound("ui_tick", nil)
			return
		}
	}
}
