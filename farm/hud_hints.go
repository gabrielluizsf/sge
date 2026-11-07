package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/i18n"
)

// hints tells what E does where the player stands, over their head, and
// how to leave a room, over the hotbar. The other keys are in the pause
// menu's Commands.
func (h *hud) hints(cv *canvas.Canvas) {
	if words, _, ok := h.f.interaction(h.farmer); ok && words != "" {
		t := h.text(words, 'y')
		// Over the farmer's head: their position is where they stand.
		at := h.farmer.Game().ToScreen(h.farmer.Transform.Position)
		above := tile + scale*6
		if h.f.ride != nil {
			above = rideHint // over a rider's head, not across them
		}
		cv.Blit(int(at.X)-t.Width/2, int(at.Y)-above-t.Height, t)
	}
	if h.f.room != nil {
		t := h.text(i18n.T("Leave by the door mat"), 'w')
		x, y := h.wayOutAt(cv, t)
		cv.Blit(x, y, t)
	}
}

// rideHint is how far above where a rider's horse stands what they may do is
// written: over the head of one up on its back.
const rideHint = 3 * tile

// wayOutAbove is how far above the hotbar's slots' middle, in pixels of
// art, the way out of a room is written: over the slots' numbers.
const wayOutAbove = 24

// wayOutAt is where t, the way out of a room, goes: in the middle over the
// hotbar, where the toasts at the top of the screen cannot cover it.
func (h *hud) wayOutAt(cv *canvas.Canvas, t *canvas.Canvas) (x, y int) {
	bar := h.f.pack.Where(backpack.Spot{Area: backpack.Bar})
	return (cv.Width - t.Width) / 2, bar.Y - wayOutAbove*scale - t.Height
}
