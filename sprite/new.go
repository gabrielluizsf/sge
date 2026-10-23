package sprite

import "github.com/gabrielluizsf/antui/canvas"

// New is a sprite of img, pivoting on its centre.
func New(img *canvas.Canvas) *Sprite {
	return &Sprite{Image: img, PivotX: 0.5, PivotY: 0.5}
}
