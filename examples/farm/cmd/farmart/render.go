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

// render paints an SVG drawing scale times its own size, its shapes
// smooth at that size, as the game paints them at its own.
func render(src string, scale int) (*canvas.Canvas, error) {
	img, err := svg.Parse(src)
	if err != nil {
		return nil, err
	}
	w, h := svgSize(src)
	return img.Render(w*scale, h*scale), nil
}
