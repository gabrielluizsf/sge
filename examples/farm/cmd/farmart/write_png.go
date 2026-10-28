package main

import (
	"image"
	"image/png"
	"os"

	"github.com/gabrielluizsf/antui/canvas"
)

// writePNG writes the canvas's pixels to a PNG file.
func writePNG(cv *canvas.Canvas, path string) error {
	img := image.NewNRGBA(image.Rect(0, 0, cv.Width, cv.Height))
	for y := 0; y < cv.Height; y++ {
		for x := 0; x < cv.Width; x++ {
			c := cv.Pixels[y*cv.Stride+x]
			i := img.PixOffset(x, y)
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = c.R(), c.G(), c.B(), c.A()
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
