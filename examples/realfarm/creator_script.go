package main

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
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

// inputOf is what the player did to the creation screen this frame.
func inputOf(win *antui.Window, typing bool) cin {
	if typing { // letters are typed, not keys that move or turn
		return cin{
			Dy:     step(win, antui.KeyUp, antui.KeyDown),
			Select: win.KeyPressed(antui.KeyEnter),
			Typed:  win.TextInput(), Back: win.KeyPressed(antui.KeyBackspace),
			X: win.MouseX(), Y: win.MouseY(),
			Press: win.MousePressed(antui.MouseLeft), Held: win.MouseDown(antui.MouseLeft),
		}
	}
	return cin{
		Dy:     step(win, antui.KeyUp, antui.KeyDown) + step(win, antui.KeyW, antui.KeyS),
		Dx:     step(win, antui.KeyLeft, antui.KeyRight) + step(win, antui.KeyA, antui.KeyD),
		Turn:   step(win, antui.KeyQ, antui.KeyE),
		Select: win.KeyPressed(antui.KeyEnter) || win.KeyPressed(antui.KeySpace),
		X:      win.MouseX(), Y: win.MouseY(),
		Press: win.MousePressed(antui.MouseLeft), Held: win.MouseDown(antui.MouseLeft),
	}
}

// Update runs the screen on what the player does.
func (c *creator) Update(o *game.Object, dt float64) {
	c.obj = o
	c.dt = dt
	if win := o.Game().Window(); win != nil {
		c.handle(inputOf(win, c.typing))
	}
	c.advance(dt)
}

// Paint draws the screen.
func (c *creator) Paint(_ *game.Object, cv *canvas.Canvas) { c.paint(cv) }
