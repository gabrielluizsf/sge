package game

import (
	"slices"

	"github.com/gabrielluizsf/antui/canvas"
)

// Draw paints the world onto cv: the background, then every visible
// object by Z, lowest first, in the order added among equals.
func (g *Game) Draw(cv *canvas.Canvas) {
	cv.Clear(g.Background)
	objs := g.Objects()
	slices.SortStableFunc(objs, func(a, b *Object) int { return a.Z - b.Z })
	for _, o := range objs {
		if o.Hidden {
			continue
		}
		o.Look().Draw(cv, o.Transform)
		for _, s := range o.scripts {
			if p, ok := s.(Painter); ok {
				p.Paint(o, cv)
			}
		}
	}
}
