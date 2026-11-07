package farm

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/vector"
)

// editMouse is the mouse while arranging: moved, the cursor goes to the
// cell under it, as far as it may; a click does what Enter does, and a
// right click what R does. Until it moves the keys have the cursor.
func (f *farm) editMouse(win *antui.Window) {
	e := f.edit
	x, y := win.MouseX(), win.MouseY()
	if !e.mouseKnown {
		e.mouseX, e.mouseY, e.mouseKnown = x, y, true
	}
	if x != e.mouseX || y != e.mouseY {
		e.mouseX, e.mouseY = x, y
		w := f.player.Game().ToWorld(vector.Vector3{X: float64(x), Y: float64(y)})
		f.editPoint(f.cellOf(w))
	}
	switch {
	case win.MousePressed(antui.MouseLeft):
		f.editAct()
	case win.MousePressed(antui.MouseRight):
		f.editTurn()
	}
}

// editPoint has the cursor go to cell col, row, as near as it may.
func (f *farm) editPoint(col, row int) {
	e := f.edit
	e.col, e.row = col, row
	f.clampCursor()
}
