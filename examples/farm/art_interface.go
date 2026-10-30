package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
)

// Art draws every picture the farm shows, so that the one game can be
// played in more than one look. The farm asks for each picture by name,
// at the screen's size, Scale pixels to a pixel of art, and keeps what it
// is given; SVGArt is the stylized look the farm command plays in.
type Art interface {
	// Tile is the picture of the tile name, one of TileNames: a cell's
	// diamond, TileW wide and TileW/2 tall, at the bottom of a picture
	// TileH tall, with what stands on the cell reaching up out of it.
	Tile(name string) *canvas.Canvas
	// Icon is the HUD's icon name, 16 by 16 pixels of art: "icon_coin",
	// "icon_bubble", "icon_can", "icon_axe_" and an axe's name, or
	// "icon_" and a seed's, a crop's or "wood".
	Icon(name string) *canvas.Canvas
	// Cursor is the outline of the cell being worked: its diamond, TileW
	// by Cell pixels.
	Cursor() *canvas.Canvas
	// Building is a sprite of the building name, a new one each time
	// though its picture may be shared, its pivot where it stands:
	// "stall", on its one cell, or a house, "house" and its kind, from 0
	// to Houses, with "_back" after it when its door faces away, on the
	// corner its two by two cells share.
	Building(name string) *sprite.Sprite
	// Houses is how many kinds of house there are; the towns are laid out
	// with as many.
	Houses() int
	// Person is how who looks and moves.
	Person(who Who) Body
	// Text is s in the HUD's letters, in the palette's ink: 'k' for text
	// on a panel, and over the world 'w' the lightest, 'y' gold, 'r'
	// coral or 'g' leaf.
	Text(s string, ink rune) *canvas.Canvas
	// Panel is a HUD panel of w by h pixels of art, highlighted when
	// selected.
	Panel(w, h int, selected bool) *canvas.Canvas
}
