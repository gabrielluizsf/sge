package tilemap

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// diamond is a picture w by w/2 with an opaque diamond in it.
func diamond(w int) *canvas.Canvas {
	img, _ := canvas.NewCanvas(w, w/2)
	for y := range w / 2 {
		half := min(y, w/2-1-y) + 1
		for x := w/2 - 2*half; x < w/2+2*half; x++ {
			img.Put(x, y, canvas.RGB(uint8(x*3), uint8(y*5), 90))
		}
	}
	return img
}

// A picture of opaque runs is copied as its blend would draw it, clipped
// at every edge of the canvas.
func TestCopiedSpansDrawAsBlitDoes(t *testing.T) {
	tile := cut(diamond(32))
	if tile.spans == nil {
		t.Fatal("a diamond is not taken for runs")
	}
	for _, at := range [][2]int{{4, 3}, {-10, -5}, {30, 20}, {-40, 0}} {
		a, _ := canvas.NewCanvas(40, 24)
		b, _ := canvas.NewCanvas(40, 24)
		a.Clear(canvas.RGB(1, 2, 3))
		b.Clear(canvas.RGB(1, 2, 3))
		copySpans(a, tile, at[0], at[1])
		b.Blit(at[0], at[1], tile.img)
		for y := range a.Height {
			for x := range a.Width {
				if a.At(x, y) != b.At(x, y) {
					t.Fatalf("at %v: pixel %d,%d is %v, not %v", at, x, y, a.At(x, y), b.At(x, y))
				}
			}
		}
	}
}

// A picture with a hole, or a pixel half seen, is blended, not copied.
func TestPicturesToBlendHaveNoSpans(t *testing.T) {
	holed := diamond(32)
	holed.Put(16, 8, canvas.Color(0))
	if spansOf(holed) != nil {
		t.Error("a diamond with a hole is taken for runs")
	}
	faint := diamond(32)
	faint.Put(16, 8, canvas.RGBA(canvas.RGBAParams{R: 10, G: 10, B: 10, A: 128}))
	if spansOf(faint) != nil {
		t.Error("a diamond with a faint pixel is taken for runs")
	}
}
