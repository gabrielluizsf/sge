package game

import (
	"math"

	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/vector"
)

// VisibleWorld is the box of the world the view shows, all of it: under
// a projection the view's corners fall on a slanted patch of ground, and
// the box holds it. Before the view is known it is empty.
func (g *Game) VisibleWorld() collider.Box {
	w, h := float64(g.viewW), float64(g.viewH)
	x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, s := range []vector.Vector3{{}, {X: w}, {Y: h}, {X: w, Y: h}} {
		p := g.ToWorld(s)
		x0, y0, x1, y1 = min(x0, p.X), min(y0, p.Y), max(x1, p.X), max(y1, p.Y)
	}
	return collider.Box{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0}
}
