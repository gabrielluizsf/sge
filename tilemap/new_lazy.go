package tilemap

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
)

// Lazy is a tile set whose tiles are drawn as they are first wanted:
// Count tiles of Width by Height pixels, Make drawing tile id. A game of
// many tiles, drawn from SVG say, starts without drawing those it does
// not show yet.
type Lazy struct {
	Count, Width, Height int
	Make                 func(id int) *canvas.Canvas
}

// NewLazyTileSet is the tile set l describes, none of it drawn yet.
func NewLazyTileSet(l Lazy) *TileSet {
	tiles := make([]*sprite.Sprite, l.Count)
	for i := range tiles {
		tiles[i] = &sprite.Sprite{Make: func() *canvas.Canvas { return l.Make(i) }}
	}
	return &TileSet{Tiles: tiles, Width: l.Width, Height: l.Height, solid: make([]bool, l.Count), opaque: make([]int8, l.Count)}
}
