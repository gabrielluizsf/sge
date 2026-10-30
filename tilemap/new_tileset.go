package tilemap

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
)

// NewTileSet cuts sheet into tiles of width by height pixels, numbered
// from 0 left to right and top to bottom.
func NewTileSet(sheet *canvas.Canvas, width, height int) (*TileSet, error) {
	tiles, err := sprite.Sheet(sheet, width, height)
	if err != nil {
		return nil, err
	}
	return &TileSet{Tiles: tiles, Width: width, Height: height, solid: make([]bool, len(tiles)), opaque: make([]int8, len(tiles))}, nil
}
