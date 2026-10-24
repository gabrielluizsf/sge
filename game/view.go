package game

// View is the size of the canvas the world was last drawn on: the
// window's, which a player may resize.
func (g *Game) View() (width, height int) {
	return g.viewW, g.viewH
}
