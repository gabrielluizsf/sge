package game

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// blendFadedSlow is BlendFaded as it was, a pixel at a time.
func blendFadedSlow(cv *canvas.Canvas, p Placed) {
	img, x, y := p.Img, p.X, p.Y
	for py := range img.Height {
		for px := range img.Width {
			if c := img.At(px, py); c.A() > 0 {
				cv.Pixel(x+px, y+py, canvas.Fade(c, int(float64(c.A())*RevealOpacity+0.5)))
			}
		}
	}
}

// Faded pictures come out exactly as the pixel-at-a-time blend made them,
// over any ground, half off the canvas and all.
func TestFadedPicturesAreBlendedAsTheyWere(t *testing.T) {
	rnd := rand.New(rand.NewSource(3))
	img, _ := canvas.NewCanvas(37, 29)
	for i := range img.Pixels {
		img.Pixels[i] = canvas.Color(rnd.Uint32())
		if i%5 == 0 {
			img.Pixels[i] = 0
		}
	}
	for _, at := range [][2]int{{10, 10}, {-20, 5}, {40, 40}, {-100, -100}, {55, -3}} {
		a, _ := canvas.NewCanvas(64, 64)
		for i := range a.Pixels {
			a.Pixels[i] = canvas.Color(rnd.Uint32()) | 0xFF000000
		}
		b, _ := canvas.NewCanvas(64, 64)
		copy(b.Pixels, a.Pixels)
		blendFadedSlow(a, Placed{Img: img, X: at[0], Y: at[1]})
		BlendFaded(b, Placed{Img: img, X: at[0], Y: at[1]})
		BlendFaded(b, Placed{Img: img, X: at[0], Y: at[1]}) // a second time, from the cache
		blendFadedSlow(a, Placed{Img: img, X: at[0], Y: at[1]})
		if !slices.Equal(a.Pixels, b.Pixels) {
			t.Fatalf("differs at %v", at)
		}
	}
}
