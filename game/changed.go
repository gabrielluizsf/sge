package game

import "github.com/gabrielluizsf/antui/canvas"

// Changed reports whether drawing the world onto cv now would put other
// pixels there than the last frame Changed answered true for did, as cv
// still holds them: the camera, the view, any object's place, look or
// visibility, the background, and every painter of the world are the same.
// When any is not, or cannot be told, as a painter that is not a Versioned
// cannot, it is true. The first frame is always changed.
//
// A frame that is not changed is left undrawn and unpresented, and costs
// nothing. Whoever asks must draw, then, when it is true, onto the same cv
// that no one else has drawn on; Frame does.
func (g *Game) Changed(cv *canvas.Canvas) bool {
	steady := g.next.fill(g, cv)
	if steady && g.next.same(&g.last) {
		return false
	}
	g.last, g.next = g.next, g.last
	return true
}

// Invalidate makes the next Changed true: for a game that has changed
// what it draws in a way it cannot tell, a sprite's image painted over in
// place, say.
func (g *Game) Invalidate() { g.last.valid = false }
