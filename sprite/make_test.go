package sprite

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/transform"
)

// A sprite with Make makes its image when first drawn, and once only.
func TestASpriteIsMadeWhenFirstDrawn(t *testing.T) {
	made := 0
	s := &Sprite{Make: func() *canvas.Canvas {
		made++
		img, _ := canvas.NewCanvas(2, 2)
		img.Clear(canvas.RGB(255, 0, 0))
		return img
	}}
	if made != 0 {
		t.Fatal("made before it was drawn")
	}
	dst, _ := canvas.NewCanvas(4, 4)
	tr := transform.New()
	tr.Position.X, tr.Position.Y = 1, 1
	s.Draw(dst, tr)
	s.Draw(dst, tr)
	if made != 1 || dst.At(1, 1) != canvas.RGB(255, 0, 0) {
		t.Errorf("made %d times, pixel %v", made, dst.At(1, 1))
	}
}
