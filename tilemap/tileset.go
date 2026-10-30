// Package tilemap is a grid of tiles drawn from a tile set: the ground,
// walls and floors of a level. A Map is a game.Object, added to a game as
// any other: it draws its tiles itself, and its solid tiles are its
// collider.
package tilemap

import "github.com/gabrielluizsf/sge/sprite"

// TileSet is the tiles a map is made of, all of one size, and which of
// them are solid: walls the player cannot walk through.
type TileSet struct {
	Tiles         []*sprite.Sprite
	Width, Height int
	// CellWidth and CellHeight, when set, are a tile's cell in the world,
	// when not its picture's size: an isometric tile's picture is wider
	// than the square it stands on, and taller where something stands.
	CellWidth, CellHeight int
	solid                 []bool
	opaque                []int8       // by tile id: whether no pixel of it lets through, once known
	anims                 []*Animation // by tile id; nil for a still tile
	animated              int          // how many tiles are animated
	trims                 []trimmed    // by tile id: its picture cut to what shows
}
