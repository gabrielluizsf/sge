package sprite

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/transform"
	"github.com/gabrielluizsf/sge/vector"
)

// A sprite only moved, not turned or scaled, is copied rather than mapped
// pixel by pixel, landing on exactly the pixels the mapping picks, at any
// fractional position and pivot.
func TestDrawUnturnedIsTheMappingsPixels(t *testing.T) {
	img, _ := canvas.NewCanvas(5, 3)
	for y := 0; y < 3; y++ {
		for x := 0; x < 5; x++ {
			img.Pixel(x, y, canvas.RGBA(canvas.RGBAParams{R: uint8(40 * x), G: uint8(80 * y), B: 200, A: uint8(120 + 20*x)}))
		}
	}
	s := &Sprite{Image: img, PivotX: 0.5, PivotY: 0.5}
	for _, at := range [][2]float64{{10, 10}, {10.5, 10.5}, {10.49, 9.51}, {-1.5, 3.2}, {7.999, 0.001}} {
		tr := transform.New()
		tr.Position = vector.Vector3{X: at[0], Y: at[1]}
		fast, _ := canvas.NewCanvas(20, 20)
		slow, _ := canvas.NewCanvas(20, 20)
		fast.Clear(canvas.RGB(9, 9, 9))
		slow.Clear(canvas.RGB(9, 9, 9))
		s.Draw(fast, tr)
		slow.BlitMatrix(img, s.Matrix(tr), canvas.Area{Width: 5, Height: 3})
		for i := range fast.Pixels {
			if fast.Pixels[i] != slow.Pixels[i] {
				t.Fatalf("at %v: pixel %d,%d is %08x, the mapping's %08x", at, i%fast.Stride, i/fast.Stride, fast.Pixels[i], slow.Pixels[i])
			}
		}
	}
}
