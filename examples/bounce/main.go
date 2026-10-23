// Command bounce opens a window with a ball bouncing around it. It is the
// same program in Go and, transpiled by gots, on Deno.
//
// BOUNCE_FRAMES, when set, closes the window after that many frames and
// steps time by a fixed 1/60 s, so a run is the same every time;
// BOUNCE_PNG then names a file the last frame is written to.
package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func main() {
	win, err := antui.Open("sge bounce", 480, 320)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	limit, _ := strconv.Atoi(os.Getenv("BOUNCE_FRAMES"))
	x, y, vx, vy := 60.0, 60.0, 180.0, 140.0
	const r = 24
	frames := 0
	for win.Begin() {
		dt := win.Delta()
		if limit > 0 {
			dt = 1.0 / 60
		}
		x, y = x+vx*dt, y+vy*dt
		if x < r || x > float64(win.Width()-r) {
			vx = -vx
		}
		if y < r || y > float64(win.Height()-r) {
			vy = -vy
		}
		win.Clear(canvas.RGB(24, 26, 33))
		win.Canvas().FillCircle(canvas.Circle{X: int(x), Y: int(y), R: r}, canvas.RGB(240, 180, 40))
		win.Canvas().Text(canvas.Point{X: 10, Y: 10}, fmt.Sprintf("frame %d", frames), canvas.RGB(220, 220, 220))
		frames++
		if limit > 0 && frames >= limit {
			if path := os.Getenv("BOUNCE_PNG"); path != "" {
				if err := writePNG(win.Canvas(), path); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
			win.End()
			break
		}
		win.End()
	}
	fmt.Println("frames:", frames)
}

// writePNG writes the canvas's pixels, 0xAARRGGBB, to a PNG file.
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
