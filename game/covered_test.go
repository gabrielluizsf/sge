package game

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A world that covers the whole canvas is not cleared first: what it
// draws is everything, and a pass over every pixel is spared.
func TestACoveredWorldIsNotCleared(t *testing.T) {
	g := New()
	g.Background = canvas.RGB(255, 0, 0)
	g.Covered = true
	cv, _ := canvas.NewCanvas(4, 4)
	cv.Clear(canvas.RGB(0, 0, 255))
	g.Draw(cv)
	if cv.At(1, 1) != canvas.RGB(0, 0, 255) {
		t.Errorf("cleared to %v", cv.At(1, 1))
	}
}
