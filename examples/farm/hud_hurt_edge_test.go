package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// The red of a blow is made once for a size, is gone with the flash, and
// is on the edge of the screen, not its middle.
func TestTheRedOfABlowLiesOnTheEdge(t *testing.T) {
	_, f := build(SVGArt(), idle{})
	h := &hud{f: f}
	base := canvas.RGB(90, 140, 70)
	cv, _ := canvas.NewCanvas(64, 36)
	cv.Clear(base)
	h.hurtEdge(cv)
	if cv.At(0, 18) != base {
		t.Error("red with no blow")
	}
	f.flash = flashFor
	h.hurtEdge(cv)
	first := h.edge
	if p := cv.At(0, 18); p.R() <= base.R() || cv.At(32, 18) != base {
		t.Errorf("edge %v, middle %v", p, cv.At(32, 18))
	}
	h.hurtEdge(cv)
	if h.edge != first || first == nil {
		t.Error("the mask was made again")
	}
}
