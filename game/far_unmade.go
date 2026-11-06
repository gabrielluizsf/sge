package game

import "github.com/gabrielluizsf/sge/transform"

// farUnmade reports whether o's sprite is not made yet and, drawn as t
// places it, is further than LazyReach outside the canvas: making it
// would cost a picture nobody can see.
func (g *Game) farUnmade(o *Object, t transform.Transform) bool {
	r := g.LazyReach
	if r <= 0 || o.Screen {
		return false
	}
	s := o.Look()
	if s == nil || s.IsReady() {
		return false
	}
	p := t.Position
	return p.X < -r || p.Y < -r || p.X > float64(g.viewW)+r || p.Y > float64(g.viewH)+r
}
