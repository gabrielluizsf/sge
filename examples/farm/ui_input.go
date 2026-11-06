package farm

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/dialogue"
	"github.com/gabrielluizsf/sge/menu"
)

// step is -1, 0 or 1: the way a pair of keys held opposite ways point.
func step(win *antui.Window, less, more antui.Key) int {
	d := 0
	if win.KeyPressed(less) {
		d--
	}
	if win.KeyPressed(more) {
		d++
	}
	return d
}

// ups and downs are the keys that move a choice up and down, and lefts
// and rights across.
func vertical(win *antui.Window) int {
	return step(win, antui.KeyUp, antui.KeyDown) + step(win, antui.KeyW, antui.KeyS)
}

func horizontal(win *antui.Window) int {
	return step(win, antui.KeyLeft, antui.KeyRight) + step(win, antui.KeyA, antui.KeyD)
}

// numberKeys are the keys of the hotbar's slots, 1 to 6.
var numberKeys = []antui.Key{antui.Key1, antui.Key2, antui.Key3, antui.Key4, antui.Key5, antui.Key6}

// packInput is what the player did to the inventory this frame.
func packInput(win *antui.Window) backpack.Input {
	in := backpack.Input{
		Dx: horizontal(win), Dy: vertical(win),
		Select: win.KeyPressed(antui.KeyEnter) || win.KeyPressed(antui.KeySpace),
		Cancel: win.KeyPressed(antui.KeyEscape),
		Quick:  win.KeyDown(antui.KeyLeftShift) || win.KeyDown(antui.KeyRightShift),
		Tab:    step(win, antui.KeyQ, antui.KeyE),
		Mouse:  backpack.Point{X: win.MouseX(), Y: win.MouseY()},
		Press:  win.MousePressed(antui.MouseLeft), Release: win.MouseReleased(antui.MouseLeft),
	}
	for i, k := range numberKeys {
		if win.KeyPressed(k) {
			in.Number = i + 1
		}
	}
	return in
}

// menuInput is what the player did to a menu this frame.
func menuInput(win *antui.Window) menu.Input {
	return menu.Input{
		Dy:     vertical(win),
		Page:   step(win, antui.KeyPageUp, antui.KeyPageDown),
		Wheel:  -win.Wheel(),
		Select: win.KeyPressed(antui.KeyEnter) || win.KeyPressed(antui.KeySpace),
		Cancel: win.KeyPressed(antui.KeyEscape),
		Mouse:  menu.Point{X: win.MouseX(), Y: win.MouseY()},
		Press:  win.MousePressed(antui.MouseLeft),
	}
}

// talkInput is what the player did to a conversation this frame; E and
// Escape leave it.
func talkInput(win *antui.Window) dialogue.Input {
	return dialogue.Input{
		Dy:     vertical(win),
		Select: win.KeyPressed(antui.KeyEnter) || win.KeyPressed(antui.KeySpace),
		Cancel: win.KeyPressed(antui.KeyEscape) || win.KeyPressed(antui.KeyE),
		Mouse:  dialogue.Point{X: win.MouseX(), Y: win.MouseY()},
		Press:  win.MousePressed(antui.MouseLeft),
	}
}
