package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
)

// townH is how tall the preview's town is: a row of houses, front and
// back, and a row of villagers under it.
const townH = 2*art.HouseH*2 + art.FarmerH*3 + 16

// drawTown draws every kind of house, its door to the viewer and away,
// at twice the size, and villagers of every skin tone, women and men, at
// three times, from x, y.
func drawTown(cv *canvas.Canvas, x, y int) error {
	for k := range art.Houses {
		for i, front := range []bool{true, false} {
			img, err := render(art.House(k, front), 2)
			if err != nil {
				return err
			}
			cv.Blit(x+(2*k+i)*(art.HouseW*2+4), y, img)
		}
	}
	y += art.HouseH*2 + 8
	for i := range 2 * len(art.Skins) {
		l := art.Villager(i%2 == 0, [5]int{i / 2, i, i / 2, i, i / 3})
		frames := art.PersonFrames(l)
		for j, name := range []string{"down", "right"} {
			img, err := render(frames[name], 3)
			if err != nil {
				return err
			}
			cv.Blit(x+(2*i+j)*(art.FarmerW*3+2), y, img)
		}
	}
	return nil
}
