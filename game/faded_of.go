package game

import "github.com/gabrielluizsf/antui/canvas"

// fadedKeep is how many faded pictures each of the two generations holds.
const fadedKeep = 64

// faded holds the faded copies of the pictures drawn faded lately, in two
// generations as a cache that lets the older go when the newer is full.
var faded struct {
	newer, older map[*canvas.Canvas]*canvas.Canvas
}

// fadedOf is img with each pixel's alpha at RevealOpacity of what it was.
func fadedOf(img *canvas.Canvas) *canvas.Canvas {
	if f, ok := faded.newer[img]; ok {
		return f
	}
	f, ok := faded.older[img]
	if !ok {
		f = fade(img)
	}
	if faded.newer == nil || len(faded.newer) >= fadedKeep {
		faded.older, faded.newer = faded.newer, map[*canvas.Canvas]*canvas.Canvas{}
	}
	faded.newer[img] = f
	return f
}

// fade is a copy of img at RevealOpacity.
func fade(img *canvas.Canvas) *canvas.Canvas {
	out, _ := canvas.NewCanvas(img.Width, img.Height)
	for y := range img.Height {
		for x := range img.Width {
			if c := img.At(x, y); c.A() > 0 {
				out.Pixels[y*out.Stride+x] = canvas.Fade(c, int(float64(c.A())*RevealOpacity+0.5))
			}
		}
	}
	return out
}
