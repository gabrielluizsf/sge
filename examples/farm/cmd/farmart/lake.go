package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/terrain"
)

// lake is a pond of odd shape, to see the shores meet: each cell takes
// the shore whose mask is its watery neighbours.
var lake = []string{
	"..........",
	"..www.....",
	".wwwww....",
	".wwwwwww..",
	"..wwwwww..",
	"....ww....",
	"..........",
}

// drawLake lays the lake out over grass at three times the size, seen
// from the corner.
func drawLake(cv *canvas.Canvas, dir string, x0, y0 int) error {
	const s = 3
	wet := func(c, r int) bool {
		return r >= 0 && r < len(lake) && c >= 0 && c < len(lake[r]) && lake[r][c] == 'w'
	}
	for r, line := range lake {
		for c := range line {
			names := []string{"grass0"}
			if wet(c, r) {
				names = append(names, fmt.Sprint("shore", terrain.Mask(wet, c, r), "_0"))
			}
			for _, name := range names {
				img, err := load(dir, name, s)
				if err != nil {
					return err
				}
				isoBlit(cv, img, [2]int{c, r}, [2]int{x0, y0 + art.IsoH*s})
			}
		}
	}
	return nil
}
