package game

import (
	"slices"

	"github.com/gabrielluizsf/antui/canvas"
)

// Draw paints the world onto cv as the camera sees it: the background,
// unless the world is Covered, then every visible object by Z, lowest
// first, in the order added among equals. Painters draw on the canvas as it is, at points of the screen.
func (g *Game) Draw(cv *canvas.Canvas) {
	g.viewW, g.viewH = cv.Width, cv.Height
	cam := g.camera()
	if !g.Covered {
		cv.Clear(g.Background)
	}
	objs := g.Objects()
	slices.SortStableFunc(objs, func(a, b *Object) int { return a.Z - b.Z })
	for _, o := range objs {
		if o.Hidden {
			continue
		}
		t := o.Transform
		if !o.Screen {
			t.Position = t.Position.Sub(cam)
		}
		o.Look().Draw(cv, t)
		for _, s := range o.scripts {
			if p, ok := s.(Painter); ok {
				p.Paint(o, cv)
			}
		}
	}
}
