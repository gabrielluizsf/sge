package main

import (
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/svg"
)

// load renders the sprite name from its SVG file in dir, as the game
// does, scaled up by scale with every pixel kept square.
func load(dir, name string, scale int) (*canvas.Canvas, error) {
	src, err := os.ReadFile(filepath.Join(dir, name+".svg"))
	if err != nil {
		return nil, err
	}
	return render(string(src), scale)
}

// render paints an SVG drawing at its own size, then enlarges it.
func render(src string, scale int) (*canvas.Canvas, error) {
	img, err := svg.Parse(src)
	if err != nil {
		return nil, err
	}
	w, h := svgSize(src)
	small := img.Render(w, h)
	big, err := canvas.NewCanvas(w*scale, h*scale)
	if err != nil {
		return nil, err
	}
	for y := 0; y < h*scale; y++ {
		for x := 0; x < w*scale; x++ {
			big.Pixels[y*big.Stride+x] = small.Pixels[(y/scale)*small.Stride+x/scale]
		}
	}
	return big, nil
}
