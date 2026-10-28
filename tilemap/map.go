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
// A map draws itself: one wholly of opaque tiles, ground, is copied in
// one go from a picture of them all, and one mostly of tiles with
// transparent pixels is blended in one go; any other, such as a layer of
// things over the ground, draws its tiles one by one, copying the opaque
// ones and blending only those that let the layers under them show.
type Map struct {
	*game.Object
	Tiles      *TileSet
	Cols, Rows int
	// Sorted has the map's tiles, in a projected world, each drawn among
	// the objects of its Z by where it stands, as the trees of a wood a
	// farmer walks between are; otherwise the map is drawn whole.
	Sorted  bool
	cells   []int
	painted []int   // the tile each cell shows now
	time    float64 // seconds the map has been in a world, for animations
	picture *sprite.Sprite
	opaques int // cells showing a tile with no transparent pixel
	glassy  int // cells showing a tile with some
}
