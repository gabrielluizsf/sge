package terrain

import (
	"fmt"
	"slices"

	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// makeChunk makes chunk cx, cy, keyed key, a map for each layer filled
// from the generator, and adds them to the world's game.
func (w *World) makeChunk(key, cx, cy int) {
	cw, ch := w.chunkSize()
	layers := make([]*tilemap.Map, w.Gen.Layers)
	for l := range layers {
		m, err := tilemap.New(fmt.Sprint(w.Name, " ", cx, ",", cy, " layer ", l), w.Gen.Tiles, w.Size, w.Size)
		if err != nil {
			panic(err)
		}
		m.Transform.Position = vector.Vector3{X: float64(cx) * cw, Y: float64(cy) * ch}
		if l < len(w.Z) {
			m.Z = w.Z[l]
		}
		m.Sorted = slices.Contains(w.Sorted, l)
		for r := range w.Size {
			for c := range w.Size {
				if id := w.Gen.Tile(l, cx*w.Size+c, cy*w.Size+r); id != tilemap.Empty {
					m.Set(c, r, id)
				}
			}
		}
		layers[l] = m
		w.Game().Add(m.Object)
	}
	c := &chunk{cx: cx, cy: cy, layers: layers}
	if w.Props != nil {
		c.props = w.Props(cx, cy)
		for _, o := range c.props {
			w.Game().Add(o)
		}
	}
	w.chunks[key] = c
}
