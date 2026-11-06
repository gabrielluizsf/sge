package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/backpack"
)

// The way out of a room is not written over a toast, nor over the hotbar:
// it was at the top right, where a long toast ran under it.
func TestTheWayOutIsClearOfTheToasts(t *testing.T) {
	_, f := build(SVGArt(), idle{})
	h := &hud{f: f}
	cv, _ := canvas.NewCanvas(1280, 720)
	f.pack.Place(backpack.View{W: cv.Width, H: cv.Height, Scale: f.pack.Scale})
	words := h.text("Leave by the door mat", 'w')
	x, y := h.wayOutAt(cv, words)
	if under := (toastTop + 3*(rowH+6)) * scale; y < under {
		t.Errorf("it is at y %d, among the toasts, which reach %d", y, under)
	}
	if bar := f.pack.Where(backpack.Spot{Area: backpack.Bar}); y+words.Height > bar.Y-iconPx*scale {
		t.Errorf("it ends at %d, on the hotbar in the middle of %d", y+words.Height, bar.Y)
	}
	if x < 0 || x+words.Width > cv.Width {
		t.Errorf("it runs from %d to %d, off the screen", x, x+words.Width)
	}
}
