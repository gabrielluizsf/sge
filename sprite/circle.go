package sprite

import "github.com/gabrielluizsf/antui/canvas"

// Circle is a sprite of a filled circle of radius r on transparency.
func Circle(r int, c canvas.Color) (*Sprite, error) {
	img, err := canvas.NewCanvas(2*r, 2*r)
	if err != nil {
		return nil, err
	}
	img.Clear(canvas.RGBA(canvas.RGBAParams{}))
	img.FillCircle(canvas.Circle{X: r, Y: r, R: r}, c)
	return New(img), nil
}
