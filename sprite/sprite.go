// Package sprite is a picture a game object is drawn with: an antui
// canvas and the point of it that sits on the object's position.
package sprite

import "github.com/gabrielluizsf/antui/canvas"

// Sprite is an image with a pivot: the point, as a fraction of the image's
// width and height, that is placed on the position it is drawn at and
// that it rotates and scales around.
type Sprite struct {
	Image          *canvas.Canvas
	PivotX, PivotY float64
}
