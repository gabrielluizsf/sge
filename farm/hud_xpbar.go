package farm

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/i18n"
)

// xpBarW is the bar's width, in pixels of art.
const xpBarW = 96

// The bar's fill, and the gold of a level up's glow.
const (
	fillInk canvas.Color = 0xFF5BBF4A
	goldInk canvas.Color = 0xFFFFD34D
)

// levelW is the width, in pixels of art, of the level badge.
func levelW(h *hud) int { return h.text(i18n.T("Lv %d", h.f.prog.Level), 'k').Width/scale + 10 }

// level is the badge "Lv n", gold while the level up glows.
func (h *hud) level(cv *canvas.Canvas, x, y int) {
	t := h.text(i18n.T("Lv %d", h.f.prog.Level), 'k')
	cv.Blit(x, y, h.panel(levelW(h), 12, h.f.glow > 0))
	cv.Blit(x+5*scale, y+(12*scale-t.Height)/2, t)
}

// xpBar is the XP bar at x, y: a track, and a fill that follows the
// XP smoothly, and starts again from empty at a new level.
func (h *hud) xpBar(cv *canvas.Canvas, x, y int) {
	cv.Blit(x, y, h.panel(xpBarW, 9, h.f.glow > 0))
	_, frac := math.Modf(h.f.shown)
	w := int(frac * float64((xpBarW-4)*scale))
	if w > 0 {
		cv.FillRect(canvas.Area{X: x + 2*scale, Y: y + 2*scale, Width: w, Height: 5 * scale}, fillInk)
	}
}

// glowRing is a gold halo around the portrait at cx, cy, fading as the
// glow runs out.
func (h *hud) glowRing(cv *canvas.Canvas, cx, cy int) {
	if h.f.glow <= 0 {
		return
	}
	a := int(150 * h.f.glow / glowFor)
	cv.FillCircle(canvas.Circle{X: cx, Y: cy, R: (faceD/2 + 3) * scale}, canvas.Fade(goldInk, a))
}
