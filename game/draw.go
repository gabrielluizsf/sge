package game

import (
	"cmp"
	"slices"

	"github.com/gabrielluizsf/antui/canvas"
)

// Draw paints the world onto cv as the camera sees it: the background,
// unless the world is Covered, then every visible object by Z, lowest
// first, in the order added among equals; or, when the world is
// projected, nearer over farther among equals. Painters draw on the
// canvas as it is, at points of the screen.
func (g *Game) Draw(cv *canvas.Canvas) {
	g.viewW, g.viewH = cv.Width, cv.Height
	if !g.Covered {
		cv.Clear(g.Background)
	}
	g.items = g.items[:0]
	for _, o := range g.Objects() {
		if !o.Hidden {
			g.collect(o)
		}
	}
	slices.SortStableFunc(g.items, func(a, b item) int {
		if a.z != b.z {
			return a.z - b.z
		}
		return cmp.Compare(a.depth, b.depth)
	})
	for _, it := range g.items {
		it.paint(cv)
	}
	clear(g.items) // let the painters go
}

// collect adds what o is drawn as to the items: its pieces, if a script
// spreads it, or itself.
func (g *Game) collect(o *Object) {
	add := func(depth float64, paint func(*canvas.Canvas)) {
		g.items = append(g.items, item{o.Z, depth, paint})
	}
	for _, s := range o.scripts {
		if sp, ok := s.(Spreader); ok && sp.Spread(o, add) {
			return
		}
	}
	add(g.depth(o), func(cv *canvas.Canvas) { g.drawObject(o, cv) })
}

// depth is where o stands on the screen, nearer the lower, in a projected
// world; in a flat one objects keep the order they were added in.
func (g *Game) depth(o *Object) float64 {
	if g.Projection == nil || o.Screen {
		return 0
	}
	return g.ToScreen(o.Transform.Position.Add(o.Footing)).Y
}

// drawObject draws o's look where the camera sees it, and its painters.
func (g *Game) drawObject(o *Object, cv *canvas.Canvas) {
	t := o.Transform
	if !o.Screen {
		t.Position = g.ToScreen(t.Position)
	}
	o.Look().Draw(cv, t)
	for _, s := range o.scripts {
		if p, ok := s.(Painter); ok {
			p.Paint(o, cv)
		}
	}
}
