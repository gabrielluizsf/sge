package terrain

import (
	"math"

	"github.com/gabrielluizsf/sge/collider"
)

// Blocks reports whether c touches a solid tile of a chunk made: walls
// are where the player is, and the player is in view.
func (w *World) Blocks(c collider.Collider) bool {
	b := c.Bounds()
	cw, ch := w.chunkSize()
	for cy := int(math.Floor(b.Y / ch)); cy <= int(math.Floor((b.Y+b.Height)/ch)); cy++ {
		for cx := int(math.Floor(b.X / cw)); cx <= int(math.Floor((b.X+b.Width)/cw)); cx++ {
			key, ok := cellKey(cx, cy)
			ch := w.chunks[key]
			if !ok || ch == nil {
				continue
			}
			for _, m := range ch.layers {
				if m.Blocks(c) {
					return true
				}
			}
		}
	}
	return false
}

// Stride is half a tile.
func (w *World) Stride() float64 {
	cw, ch := w.Gen.Tiles.Cell()
	return float64(min(cw, ch)) / 2
}
