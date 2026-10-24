package tilemap

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
)

// Empty is a cell with no tile.
const Empty = -1

// Map is a grid of Cols by Rows cells, each a tile of Tiles or Empty. It
// is a game.Object: its position is the map's top-left corner, its Z and
// Hidden work as for any object, and its scale stretches it. Rotating a
// map turns its picture, not its collisions.
//
// A map at least half full of tiles is drawn as one picture of them all,
// its Sprite; a sparser one, such as a layer of things over the ground,
// has no Sprite and draws its tiles one by one, sparing the blending of
// all the empty cells.
type Map struct {
	*game.Object
	Tiles      *TileSet
	Cols, Rows int
	cells      []int
	painted    []int   // the tile each cell shows now
	time       float64 // seconds the map has been in a world, for animations
	picture    *sprite.Sprite
	filled     int // cells holding a tile
}
