package terrain

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
)

// World is an endless map made by a generator: chunks of Size by Size
// cells, a tilemap.Map for each layer, made as the camera comes near them
// and dropped when it has gone, so a window of any size, resized or not,
// is always filled. Its layers are drawn at the Zs given; its solid tiles
// are walls for tilemap.Move.
type World struct {
	*game.Object
	Gen  *Generator
	Size int
	Z    []int
	// Margin is how many chunks beyond the view are kept made, so walking
	// meets chunks made already. Those within the view are made at once,
	// those in the margin a few a frame.
	Margin int
	// Sorted are the layers whose maps are tilemap.Map's Sorted: in a
	// projected world their tiles stand among the objects of their Z.
	Sorted []int
	chunks map[int]*chunk // by cellKey of the chunk's place
}

// chunk is the maps of a chunk's layers, and its place among chunks.
type chunk struct {
	cx, cy int
	layers []*tilemap.Map
}
