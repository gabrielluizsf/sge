package sprite

import (
	"image"
	"image/color"

	"github.com/gabrielluizsf/antui/canvas"
)

// FromImage is a sprite of a decoded image.
func FromImage(img image.Image) (*Sprite, error) {
	b := img.Bounds()
	cv, err := canvas.NewCanvas(b.Dx(), b.Dy())
	if err != nil {
		return nil, err
	}
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			c := color.NRGBAModel.Convert(img.At(b.Min.X+x, b.Min.Y+y)).(color.NRGBA)
			cv.Pixels[y*cv.Stride+x] = canvas.RGBA(canvas.RGBAParams{R: c.R, G: c.G, B: c.B, A: c.A})
		}
	}
	return New(cv), nil
}
