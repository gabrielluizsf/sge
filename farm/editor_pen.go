package farm

import "github.com/gabrielluizsf/sge/i18n"

// penReach is how far from the player, in cells, a pen may be moved.
const penReach = 14

// penUnder is the pen of the farm whose cells include col, row, or nil.
func (f *farm) penUnder(col, row int) *pen {
	if f.herd == nil || f.room != nil {
		return nil
	}
	for _, p := range f.herd.pens {
		if p.town == nil && p.r.has(col, row) {
			return p
		}
	}
	return nil
}

// liftPen starts moving the pen under the cursor, if there is one, and
// reports whether it did. Its fence stays up until it is put down.
func (f *farm) liftPen() bool {
	e := f.edit
	p := f.penUnder(e.col, e.row)
	if p == nil {
		return false
	}
	e.pen, e.penCol, e.penRow = p, e.col-p.r.col, e.row-p.r.row
	f.sound("ui_tick", nil)
	return true
}

// penGhost is where the pen being moved would stand, and whether it may.
func (f *farm) penGhost() (rect, bool) {
	e := f.edit
	r := rect{e.col - e.penCol, e.row - e.penRow, penCols, penRows}
	return r, f.penFits(r, e.pen)
}

// putPen stands the pen being moved where the cursor has it, if it fits.
func (f *farm) putPen() {
	r, ok := f.penGhost()
	if !ok {
		f.noFit()
		return
	}
	f.movePen(f.edit.pen, r)
	f.edit.pen = nil
	f.notify(i18n.T("Pen moved"))
	f.sound("ui_open", nil)
}
