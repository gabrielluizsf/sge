package game

import "github.com/gabrielluizsf/antui"

// firstStep is the most time the first frame of a game shown in a window
// lets pass: the time the game took to be made, before it was shown, is
// not lived by its world.
const firstStep = 1.0 / 30

// Show steps and draws the world in win every frame until win is closed
// or Quit is called, and leaves win open for whatever is shown next, so
// a game's scenes, a loading screen, a menu and the play, follow one
// another in one window. It reports whether win is still open.
func (g *Game) Show(win *antui.Window) bool {
	g.window = win
	defer func() { g.window = nil }()
	if a := g.openSound(); a != nil {
		defer a.Close()
	}
	first := true
	for !g.quit {
		if !win.Begin() {
			return false
		}
		dt := win.Delta()
		if first {
			dt, first = min(dt, firstStep), false
		}
		if g.FixedStep > 0 {
			dt = g.FixedStep
		}
		// A frame the same as the one on the screen is not drawn, and End
		// presents nothing for it.
		g.Frame(win.Canvas(), dt)
		win.End()
	}
	return true
}
