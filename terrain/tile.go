package terrain

import "github.com/gabrielluizsf/sge/tilemap"

// Tile is the tile cell col, row shows on layer: the one Edit set there,
// or its areas'.
func (g *Generator) Tile(layer, col, row int) int {
	if len(g.edits) > 0 {
		if id, ok := g.edited(layer, col, row); ok {
			return id
		}
	}
	for i := len(g.Areas) - 1; i >= 0; i-- {
		if a := g.Areas[i]; a.Layer == layer && a.In(col, row) {
			return a.Tile(col, row)
		}
	}
	return tilemap.Empty
}
