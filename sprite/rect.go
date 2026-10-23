package sprite

import "github.com/gabrielluizsf/antui/canvas"

// Rect is a sprite of a filled rectangle.
func Rect(width, height int, c canvas.Color) (*Sprite, error) {
	img, err := canvas.NewCanvas(width, height)
	if err != nil {
		return nil, err
	}
	img.Clear(c)
	return New(img), nil
}
