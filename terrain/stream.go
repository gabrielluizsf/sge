package terrain

import (
	"math"

	"github.com/gabrielluizsf/sge/game"
)

// stream is the script that keeps a world's chunks made around the view.
type stream struct{ w *World }

// marginPerFrame is how many chunks of the margin are made a frame, so
// walking does not make a frame late making a row of them.
const marginPerFrame = 2

func (s stream) Update(o *game.Object, dt float64) {
	w := s.w
	g := o.Game()
	cw, ch := w.chunkSize()
	// The chunks the view shows, slanted as a projection may see them, and
	// those around them.
	v := g.VisibleWorld()
	x0, y0 := int(math.Floor(v.X/cw)), int(math.Floor(v.Y/ch))
	x1, y1 := int(math.Floor((v.X+v.Width)/cw)), int(math.Floor((v.Y+v.Height)/ch))
	m := w.Margin
	made := 0
	for cy := y0 - m; cy <= y1+m; cy++ {
		for cx := x0 - m; cx <= x1+m; cx++ {
			key, ok := cellKey(cx, cy)
			if !ok {
				continue // the world ends a million chunks out
			}
			if _, ok := w.chunks[key]; ok {
				continue
			}
			seen := cx >= x0 && cx <= x1 && cy >= y0 && cy <= y1
			if seen || made < marginPerFrame {
				w.makeChunk(key, cx, cy)
				made++
			}
		}
	}
	// Chunks a chunk further than the margin are let go: a step back
	// across a chunk's edge does not make it again.
	for k, c := range w.chunks {
		if c.cx < x0-m-1 || c.cx > x1+m+1 || c.cy < y0-m-1 || c.cy > y1+m+1 {
			for _, l := range c.layers {
				l.Destroy()
			}
			for _, o := range c.props {
				o.Destroy()
			}
			delete(w.chunks, k)
		}
	}
}
