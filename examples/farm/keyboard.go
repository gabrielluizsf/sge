package main

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// keyboard is a player at the keys: arrows or WASD walk, 1 to 3 pick what
// to hold, space uses it, and E opens the shop at the stall.
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
		return intent{}
	}
	for i, k := range []antui.Key{antui.Key1, antui.Key2, antui.Key3} {
		if win.KeyPressed(k) {
			f.slot = i
		}
	}
	return intent{dir: arrows(win), use: win.KeyPressed(antui.KeySpace)}
}

// arrows is the way the arrow or WASD keys held down point.
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
