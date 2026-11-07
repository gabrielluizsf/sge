package main

import "github.com/gabrielluizsf/antui/canvas"

// The pale warm edge round an icon, so that a dark thing, an iron tool, is
// seen on the dark panel: its colour, how strong at its strongest, and how
// many pixels it reaches.
const (
	haloR, haloG, haloB = 240, 222, 176
	haloStrength        = 0.6
	haloReach           = 2
)

// halo is the icon cv with a pale edge laid behind it where it is clear
// beside what it shows. What it shows is left as it was, pixel for pixel.
func halo(cv *canvas.Canvas) *canvas.Canvas {
	out, err := canvas.NewCanvas(cv.Width, cv.Height)
	if err != nil {
		return cv
	}
	for y := range cv.Height {
		for x := range cv.Width {
			p := cv.Pixels[y*cv.Stride+x]
			out.Pixels[y*out.Stride+x] = p
			if p.A() == 255 {
				continue
			}
			if a := rim(cv, x, y); a > 0 {
				out.Pixels[y*out.Stride+x] = over(p, a)
			}
		}
	}
	return out
}

// rim is how strong the edge is at x, y: the most of what is shown near it,
// fainter the further off.
func rim(cv *canvas.Canvas, x, y int) float64 {
	best := 0.0
	for dy := -haloReach; dy <= haloReach; dy++ {
		for dx := -haloReach; dx <= haloReach; dx++ {
			nx, ny := x+dx, y+dy
			if nx < 0 || ny < 0 || nx >= cv.Width || ny >= cv.Height {
				continue
			}
			far := max(abs(dx), abs(dy))
			if far == 0 {
				continue
			}
			a := float64(cv.Pixels[ny*cv.Stride+nx].A()) / 255 * haloStrength / float64(far)
			best = max(best, a)
		}
	}
	return best
}

// over is the pixel p over the edge's colour at strength a.
func over(p canvas.Color, a float64) canvas.Color {
	pa := float64(p.A()) / 255
	total := pa + a*(1-pa)
	if total <= 0 {
		return p
	}
	mix := func(c uint8, h float64) uint8 {
		return uint8((float64(c)*pa + h*a*(1-pa)) / total)
	}
	return canvas.RGBA(canvas.RGBAParams{R: mix(p.R(), haloR), G: mix(p.G(), haloG), B: mix(p.B(), haloB), A: uint8(total * 255)})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
