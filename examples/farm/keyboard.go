package farm

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// keyboard is a player at the keys and the mouse: arrows or WASD walk the
// way they point on the screen, faster with Shift held, 1 to 6 take in hand what the hotbar holds (what is on it is chosen in the inventory, with I or Tab), space or a
// click uses it, on the cell under the mouse when it is in reach and on
// the one in front otherwise, and E talks, goes in, sleeps or shops, as
// where they stand allows. In a room, the doormat is the way out. A
// question asked waits for its answer first; F9 passes a year at once.
type keyboard struct{}

func (keyboard) control(f *farm, o *game.Object, _ float64) intent {
	win := o.Game().Window()
	if win == nil {
		return intent{}
	}
	switch {
	case f.question != nil:
		f.questionKeys(win)
		return intent{}
	case win.KeyPressed(antui.KeyF9):
		f.passDays(daysPerYear)
		return intent{}
	case f.talk != nil:
		return intent{}
	case f.shop.open:
		f.shopKeys(win)
		return intent{}
	case f.edit != nil:
		f.editKeys(win)
		return intent{}
	case f.seat != nil:
		f.seatKeys(win)
		return intent{}
	}
	if f.riding(o) && (win.KeyPressed(antui.KeyE) || win.KeyPressed(antui.KeySpace)) {
		f.getDown()
		return intent{}
	}
	if _, row := f.cellOf(feet(o)); f.room != nil && row >= roomRows-1 {
		f.leave()
		return intent{}
	}
	if win.KeyPressed(antui.KeyB) && f.startEdit("") {
		return intent{}
	}
	if win.KeyPressed(antui.KeyE) {
		if _, do, ok := f.interaction(o); ok {
			do()
			if f.hands.useNow {
				f.hands.useNow = false
				return intent{use: true, aim: mouseCell(f, o)}
			}
			return intent{}
		}
	}
	for i, k := range numberKeys {
		if win.KeyPressed(k) && isFurnishing(f.pack.Hotbar[i]) && f.startEdit(f.pack.Hotbar[i]) {
			return intent{}
		}
		if win.KeyPressed(k) && f.hold(i) {
			f.sound("ui_tick", nil)
		}
	}
	dir := arrows(win)
	if pr := o.Game().Projection; pr != nil {
		dir = pr.Unproject(dir) // the screen's way, on the ground
	}
	run := win.KeyDown(antui.KeyLeftShift) || win.KeyDown(antui.KeyRightShift)
	return intent{dir: dir, aim: mouseCell(f, o), use: win.KeyPressed(antui.KeySpace) || win.MousePressed(antui.MouseLeft), run: run,
		strike: win.KeyPressed(antui.KeyF)}
}

// mouseCell is the cell under the mouse, if its middle is within the
// farmer's reach, or nil.
func mouseCell(f *farm, o *game.Object) *[2]int {
	g := o.Game()
	win := g.Window()
	w := g.ToWorld(vector.Vector3{X: float64(win.MouseX()), Y: float64(win.MouseY())})
	col, row := f.cellOf(w)
	if f.ground.CellCenter(col, row).Sub(feet(o)).Length() > reach {
		return nil
	}
	return &[2]int{col, row}
}

// arrows is the way the arrow or WASD keys held down point on the screen.
func arrows(win *antui.Window) vector.Vector3 {
	var d vector.Vector3
	if win.KeyDown(antui.KeyLeft) || win.KeyDown(antui.KeyA) {
		d.X--
	}
	if win.KeyDown(antui.KeyRight) || win.KeyDown(antui.KeyD) {
		d.X++
	}
	if win.KeyDown(antui.KeyUp) || win.KeyDown(antui.KeyW) {
		d.Y--
	}
	if win.KeyDown(antui.KeyDown) || win.KeyDown(antui.KeyS) {
		d.Y++
	}
	return d
}
