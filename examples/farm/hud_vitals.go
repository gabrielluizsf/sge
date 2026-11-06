package farm

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/vitals"
)

// The vitals card: a bar of health and one of food, each beside its icon,
// next to the player's card. Their fills; vitalW the bar's width in pixels
// of art; low the fraction under which a bar pulses.
const (
	vitalW                 = 52
	cardW                  = 16 + vitalW + 8
	low                    = 0.25
	healthInk canvas.Color = 0xFFD23B3C
	foodInk   canvas.Color = 0xFFE0A030
)

// vitalsCard draws the two bars on a card beside the player's card, which
// is left art pixels wide.
func (h *hud) vitalsCard(cv *canvas.Canvas, left int) {
	x, y := (navSide+left+4)*scale, navTop*scale
	cv.Blit(x, y, h.panel(cardW, navH, false))
	h.vitalBar(cv, x+3*scale, y+1*scale, vitalRow{"icon_heart", h.f.vit.Health, healthInk})
	h.vitalBar(cv, x+3*scale, y+17*scale, vitalRow{"icon_food", h.f.vit.Food, foodInk})
}

// vitalRow is one bar: its icon, its value from 0 to vitals.Max, its colour.
type vitalRow struct {
	icon string
	of   float64
	ink  canvas.Color
}

// vitalBar draws a row at x, y: the icon and the bar, which pulses when it
// is about to run out.
func (h *hud) vitalBar(cv *canvas.Canvas, x, y int, r vitalRow) {
	cv.Blit(x, y, h.icon(r.icon))
	bx, by := x+16*scale, y+(16-9)/2*scale
	cv.Blit(bx, by, h.panel(vitalW, 9, false))
	frac := r.of / vitals.Max
	w := int(frac * float64((vitalW-4)*scale))
	if w <= 0 {
		return
	}
	ink := r.ink
	if frac < low {
		ink = canvas.Fade(ink, 150+int(105*math.Sin(h.t*9)))
	}
	cv.FillRect(canvas.Area{X: bx + 2*scale, Y: by + 2*scale, Width: w, Height: 5 * scale}, ink)
}
