package building

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/random"
)

// stripes paints an awning's canvas: broad stripes, red and cream,
// running down it, the weave's threads showing and the cloth a little
// faded in places; it tiles, eight stripes to a repeat.
func stripes() *image.NRGBA {
	const size, bands = 256, 8
	red, cream := [3]float64{168, 44, 38}, [3]float64{226, 214, 186}
	src := random.New(7)
	fade := [16]float64{}
	for i := range fade {
		fade[i] = 0.92 + 0.08*src.Float()
	}
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			c := cream
			if (x*bands/size)%2 == 0 {
				c = red
			}
			// The weave: threads across and down, alternately raised.
			weave := 0.94 + 0.06*math.Sin(float64(x)*math.Pi/1.5)*math.Sin(float64(y)*math.Pi/1.5)
			k := weave * fade[(y*16/size+x*3/size)%16] * (0.97 + 0.03*src.Float())
			img.SetNRGBA(x, y, color.NRGBA{uint8(c[0] * k), uint8(c[1] * k), uint8(c[2] * k), 255})
		}
	}
	return img
}
