// Package sprite is a picture a game object is drawn with: an antui
// canvas and the point of it that sits on the object's position.
package sprite

import "github.com/gabrielluizsf/antui/canvas"

// Sprite is an image with a pivot: the point, as a fraction of the image's
// width and height, that is placed on the position it is drawn at and
// that it rotates and scales around.
//
// Make, when set, makes the image the first time the sprite is drawn,
// while Image is nil: a character's many frames rendered from SVG as each
// is first shown, rather than all at once when the character appears.
type Sprite struct {
	Image          *canvas.Canvas
	PivotX, PivotY float64
	Make           func() *canvas.Canvas
}
