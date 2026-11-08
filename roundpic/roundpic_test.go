package roundpic

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func square(w, h int, c canvas.Color) *canvas.Canvas {
	cv, _ := canvas.NewCanvas(w, h)
	cv.FillRect(canvas.Area{Width: w, Height: h}, c)
	return cv
}

func TestCropIsACircleWithATransparentCornerAndARing(t *testing.T) {
	ring := canvas.Color(0xFF112233)
	got := Crop(square(40, 60, 0xFFCC8844), 24, ring)
	if got.Width != 24 || got.Height != 24 {
		t.Fatalf("size %dx%d", got.Width, got.Height)
	}
	for _, p := range [][2]int{{0, 0}, {23, 0}, {0, 23}, {23, 23}} {
		if a := got.At(p[0], p[1]).A(); a != 0 {
			t.Errorf("corner %v alpha %d", p, a)
		}
	}
	if c := got.At(12, 12); c != 0xFFCC8844 {
		t.Errorf("centre %08x", uint32(c))
	}
	if c := got.At(12, 1); c.A() < 128 || c.R() != 0x11 {
		t.Errorf("top edge is not the ring: %08x", uint32(c))
	}
}

func TestNeutralIsACircleToo(t *testing.T) {
	got := Neutral(20, 0xFF888888, 0xFF222222)
	if got.At(0, 0).A() != 0 || got.At(10, 10).A() != 255 {
		t.Errorf("corner %08x centre %08x", uint32(got.At(0, 0)), uint32(got.At(10, 10)))
	}
}

func TestACropOfNilIsNeutral(t *testing.T) {
	if got := Crop(nil, 16, 0xFF000000); got == nil || got.At(8, 8).A() != 255 {
		t.Error("nil source")
	}
}
