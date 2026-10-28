package terrain

import "github.com/gabrielluizsf/sge/tilemap"

// Generator is a world's areas, over its tiles, in Layers layers. On each
// layer the last area a cell is in decides its tile; a cell in none of a
// layer's areas is Empty there.
type Generator struct {
	Tiles  *tilemap.TileSet
	Layers int
	Areas  []Area
	// edits are tiles set by Edit over what the areas give, by layer and
	// cell: a tree cut down stays cut.
	edits map[int]int
}

// Add puts a over the areas there are, and returns it.
func (g *Generator) Add(a Area) Area {
	g.Areas = append(g.Areas, a)
	return a
}
