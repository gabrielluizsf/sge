package town

import "github.com/gabrielluizsf/sge/game"

// reach is how many cells from the camera a town wakes: Range, or, with
// Ahead, as far as the view shows and Ahead more.
func (k *Keeper) reach(g *game.Game) int {
	if k.Ahead <= 0 {
		return k.Range
	}
	v := g.VisibleWorld()
	return max(k.Range, int(max(v.Width, v.Height)/2/k.Cell)+k.Ahead)
}
