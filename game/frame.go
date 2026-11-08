package game

import "github.com/gabrielluizsf/antui/canvas"

// Frame is a frame of the game on cv: the world steps dt seconds, always,
// so scripts, animations, timers and sounds go on whatever is seen, then
// the world is drawn, unless it would come out the very pixels cv holds
// already (see Changed). It reports whether it drew.
func (g *Game) Frame(cv *canvas.Canvas, dt float64) bool {
	// The view is the canvas's before the step, so scripts that fill it,
	// as a world made up around the camera does, see its size from the
	// first frame and after every resize.
	g.viewW, g.viewH = cv.Width, cv.Height
	g.Step(dt)
	if !g.Changed(cv) {
		return false
	}
	g.Draw(cv)
	return true
}
