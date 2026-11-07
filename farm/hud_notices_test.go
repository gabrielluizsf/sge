package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// toastRows is the first and last pixel row of the screen's middle band
// where something other than the world drew, versus a view without it.
func drawn(t *testing.T, f *farm, with bool) *canvas.Canvas {
	t.Helper()
	cv, _ := canvas.NewCanvas(1280, 720)
	f.notices = nil
	if with {
		f.notify("Dia 3: bom dia")
		for range 30 {
			f.fade(1.0 / 60)
		}
	}
	f.player.Game().Draw(cv)
	return cv
}

func TestAToastIsDrawnNearTheTopNotAtAThirdOfTheScreen(t *testing.T) {
	g, f := newGame(idle{})
	_ = g
	without := drawn(t, f, false)
	with := drawn(t, f, true)
	top, bottom := -1, -1
	for y := 0; y < 720; y++ {
		for x := 400; x < 880; x++ {
			if with.At(x, y) != without.At(x, y) {
				if top < 0 {
					top = y
				}
				bottom = y
			}
		}
	}
	if top < navBottom*scale || bottom > 160 {
		t.Errorf("toast spans rows %d..%d, want near the top", top, bottom)
	}
}
