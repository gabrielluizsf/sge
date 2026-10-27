package main

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// keyboard is a player at the keys and the mouse: arrows or WASD walk the
// way they point on the screen, 1 to 3 pick what to hold, space or a
// click uses it, on the cell under the mouse when it is in reach and on
// the one in front otherwise, and E opens the shop at the stall.
type keyboard struct{}

func (keyboard) control(f *farm, o *game.Object) intent {
	win := o.Game().Window()
	if win == nil {
		return intent{}
	}
	if f.shop.open {
		f.shopKeys(win)
		return intent{}
	}
	if win.KeyPressed(antui.KeyE) && f.atStall(o) {
		f.shop.open, f.shop.choice = true, 0
		f.sound("ui_open", nil)
		return intent{}
	}
	for i, k := range []antui.Key{antui.Key1, antui.Key2, antui.Key3} {
		if win.KeyPressed(k) && f.slot != i {
			f.slot = i
			f.sound("ui_tick", nil)
		}
	}
	dir := arrows(win)
	if pr := o.Game().Projection; pr != nil {
		dir = pr.Unproject(dir) // the screen's way, on the ground
	}
	return intent{dir: dir, aim: mouseCell(f, o), use: win.KeyPressed(antui.KeySpace) || win.MousePressed(antui.MouseLeft)}
}

// mouseCell is the cell under the mouse, if its middle is within the
// farmer's reach, or nil.
func mouseCell(f *farm, o *game.Object) *[2]int {
	g := o.Game()
	win := g.Window()
	w := g.ToWorld(vector.Vector3{X: float64(win.MouseX()), Y: float64(win.MouseY())})
	col, row, ok := f.ground.CellAt(w.X, w.Y)
	if !ok || f.ground.CellCenter(col, row).Sub(feet(o)).Length() > reach {
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
