package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/toast"
)

// toastTop is where the first toast rests, in pixels of art from the top:
// under the navbar.
const toastTop = navBottom + 4

// toastSlide is how far, in pixels of art, a toast comes from above.
const toastSlide = 10

// notices are the toasts, top centre, the oldest highest, each softly
// panelled and as opaque as its moment says.
func (h *hud) notices(cv *canvas.Canvas) {
	for _, n := range h.f.notices {
		t := h.text(n.Text, 'k')
		w := t.Width/scale + 16
		pitch := float64(rowH+6) * scale
		x := (cv.Width - w*scale) / 2
		y := toastTop*scale + int((n.Row*pitch)+(n.Shift()*toastSlide*scale))
		o := n.Opacity()
		toast.BlitAlpha(cv, x, y, h.panel(w, rowH+2, false), o)
		toast.BlitAlpha(cv, x+8*scale, y+((rowH+2)*scale-t.Height)/2, t, o)
	}
}
