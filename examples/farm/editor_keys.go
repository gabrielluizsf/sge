package farm

import "github.com/gabrielluizsf/antui"

// editKeys is what the player pressed while arranging: arrows or WASD move
// the cursor, Enter or Space lifts the piece under it, puts down the one
// held, or takes one from the bag; R turns it, X puts it in the bag, Q and
// E choose the piece of the bag, 1-6 one on the hotbar, Esc lets go of the
// piece held and then leaves, and so does B.
func (f *farm) editKeys(win *antui.Window) {
	e := f.edit
	f.editMouse(win)
	if f.edit == nil {
		return
	}
	if dx, dy := step(win, antui.KeyLeft, antui.KeyRight)+step(win, antui.KeyA, antui.KeyD), step(win, antui.KeyUp, antui.KeyDown)+step(win, antui.KeyW, antui.KeyS); dx != 0 || dy != 0 {
		f.editMove(dx, dy)
	}
	switch {
	case win.KeyPressed(antui.KeyEscape) && (e.hold != nil || e.pen != nil):
		f.dropHold()
		f.sound("ui_close", nil)
	case win.KeyPressed(antui.KeyEscape) || win.KeyPressed(antui.KeyB):
		f.stopEdit()
	case win.KeyPressed(antui.KeyEnter) || win.KeyPressed(antui.KeySpace):
		f.editAct()
	case win.KeyPressed(antui.KeyR):
		f.editTurn()
	case win.KeyPressed(antui.KeyX) || win.KeyPressed(antui.KeyBackspace):
		f.editStore()
	case win.KeyPressed(antui.KeyQ):
		f.editCycle(-1)
	case win.KeyPressed(antui.KeyE):
		f.editCycle(1)
	}
	for i, k := range numberKeys {
		if win.KeyPressed(k) && i < len(f.pack.Hotbar) && f.canChoose(f.pack.Hotbar[i]) {
			e.item = f.pack.Hotbar[i]
		}
	}
}
