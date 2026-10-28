package sprite

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/transform"
	"github.com/gabrielluizsf/sge/vector"
)

var red = canvas.RGB(255, 0, 0)

func blank(t *testing.T, w, h int) *canvas.Canvas {
	t.Helper()
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		t.Fatal(err)
	}
	cv.Clear(canvas.RGB(0, 0, 0))
	return cv
}

func TestDrawCentresOnThePosition(t *testing.T) {
	s, err := Rect(4, 4, red)
	if err != nil {
		t.Fatal(err)
	}
	cv := blank(t, 20, 20)
	tr := transform.New()
	tr.Position = vector.Vector3{X: 10, Y: 10}
	s.Draw(cv, tr)
	if got := cv.At(10, 10); got != red {
		t.Errorf("centre = %08x, want red", got)
	}
	if got := cv.At(8, 8); got != red {
		t.Errorf("top-left corner of the sprite = %08x, want red", got)
	}
	if got := cv.At(13, 13); got == red {
		t.Error("the sprite reaches past its own size")
	}
}

func TestDrawScales(t *testing.T) {
	s, _ := Rect(2, 2, red)
	cv := blank(t, 20, 20)
	tr := transform.New()
	tr.Position = vector.Vector3{X: 10, Y: 10}
	tr.Scale = vector.Vector3{X: 4, Y: 4, Z: 1}
	s.Draw(cv, tr)
	if got := cv.At(7, 7); got != red {
		t.Errorf("a pixel the scaled sprite covers = %08x, want red", got)
	}
}

func TestDrawRotates(t *testing.T) {
	// A wide bar turned a quarter stands upright.
	s, _ := Rect(10, 2, red)
	cv := blank(t, 30, 30)
	tr := transform.New()
	tr.Position = vector.Vector3{X: 15, Y: 15}
	tr.Rotate(vector.Vector3{Z: 1}, math.Pi/2)
	s.Draw(cv, tr)
	if got := cv.At(15, 11); got != red {
		t.Errorf("above the centre = %08x, want red", got)
	}
	if got := cv.At(11, 15); got == red {
		t.Error("left of the centre is still painted after the turn")
	}
}

func TestPivotTopLeft(t *testing.T) {
	s, _ := Rect(3, 3, red)
	s.PivotX, s.PivotY = 0, 0
	cv := blank(t, 10, 10)
	tr := transform.New()
	tr.Position = vector.Vector3{X: 2, Y: 2}
	s.Draw(cv, tr)
	if cv.At(2, 2) != red || cv.At(4, 4) != red || cv.At(1, 1) == red {
		t.Error("a top-left pivot does not put the corner on the position")
	}
}

func TestFromImageKeepsColours(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.NRGBA{10, 20, 30, 255})
	img.Set(1, 0, color.NRGBA{200, 100, 0, 128})
	s, err := FromImage(img)
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Image.At(0, 0); got != canvas.RGB(10, 20, 30) {
		t.Errorf("pixel 0 = %08x", got)
	}
	if got := s.Image.At(1, 0); got.A() != 128 || got.R() != 200 {
		t.Errorf("pixel 1 = %08x", got)
	}
}

func TestSheetCutsFramesInReadingOrder(t *testing.T) {
	img := blank(t, 4, 4)
	img.Pixels[0*img.Stride+2] = red // frame 1's top-left
	img.Pixels[2*img.Stride+0] = red // frame 2's top-left
	frames, err := Sheet(img, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 4 {
		t.Fatalf("%d frames, want 4", len(frames))
	}
	if frames[1].Image.At(0, 0) != red || frames[2].Image.At(0, 0) != red || frames[0].Image.At(0, 0) == red {
		t.Error("frames are not cut left to right, top to bottom")
	}
	if _, err := Sheet(img, 5, 1); err == nil {
		t.Error("a frame wider than the sheet is not an error")
	}
}

func TestDrawingNothingIsNothing(t *testing.T) {
	var s *Sprite
	s.Draw(blank(t, 2, 2), transform.New())
}
