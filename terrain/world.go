package terrain

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/warm"
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
	// Props, when set, are the objects that stand on chunk cx, cy, houses
	// say: added with its maps and destroyed with them, so they are only
	// in the game while near the view.
	Props func(cx, cy int) []*game.Object
	// Warm, when set, is given the pictures of the tiles of each chunk
	// made, to draw a few a frame: a tile of a kind not seen yet is drawn
	// from SVG, say, and would hold a frame up were it first drawn when
	// it comes into view; the margin of chunks gives it the time.
	Warm   *warm.Queue
	chunks map[int]*chunk // by cellKey of the chunk's place
}

// chunk is the maps of a chunk's layers, and its place among chunks.
type chunk struct {
	cx, cy int
	layers []*tilemap.Map
	props  []*game.Object
}
