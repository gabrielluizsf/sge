package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// drawnRows is the first and last row, within columns x0 to x1, where a
// frame with the change differs from one without.
func changedRows(a, b *canvas.Canvas, x0, x1 int) (int, int) {
	top, bottom := -1, -1
	for y := 0; y < a.Height; y++ {
		for x := x0; x < x1; x++ {
			if a.At(x, y) != b.At(x, y) {
				if top < 0 {
					top = y
				}
				bottom = y
			}
		}
	}
	return top, bottom
}

func TestTheNavbarIsAtTheTopAndToastsStandUnderIt(t *testing.T) {
	_, f := newGame(idle{})
	with := drawn(t, f, true)
	without := drawn(t, f, false)
	top, bottom := changedRows(with, without, 400, 880)
	if top < navBottom*scale {
		t.Errorf("a toast starts at row %d, over the navbar that ends at %d", top, navBottom*scale)
	}
	if bottom > 160 {
		t.Errorf("a toast reaches row %d", bottom)
	}
}

func TestTheNavbarHoldsTheLevelAndTheXPNumbers(t *testing.T) {
	_, f := newGame(idle{})
	f.earn(30)
	f.shown, f.floats = f.xpValue(), nil
	a := drawn(t, f, false)
	f.earn(40)
	f.shown, f.floats = f.xpValue(), nil
	b := drawn(t, f, false)
	top, bottom := changedRows(a, b, 0, 500)
	if top < 0 || bottom >= navBottom*scale {
		t.Errorf("XP changed rows %d..%d, not inside the bar of %d", top, bottom, navBottom*scale)
	}
}

func TestThePortraitIsMadeOnceAndIsACircle(t *testing.T) {
	_, f := newGame(idle{})
	h := &hud{f: f}
	a := h.portrait(40)
	if b := h.portrait(40); a != b {
		t.Error("masked again")
	}
	if a.At(0, 0).A() != 0 || a.At(20, 20).A() == 0 {
		t.Error("not a circle")
	}
}
