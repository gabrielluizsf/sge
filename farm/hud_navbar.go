package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/i18n"
)

// The navbar, in pixels of art from the top left: two cards along the top,
// the player's on the left (portrait, level, name, XP) and the day's and
// the purse's on the right. It ends at navBottom; toasts stand under it.
const (
	navTop    = 4
	navH      = 34
	navBottom = navTop + navH
	navSide   = 6
	faceD     = 26 // the portrait's width, round
)

// navbar draws the two cards.
func (h *hud) navbar(cv *canvas.Canvas) {
	h.vitalsCard(cv, h.playerCard(cv))
	h.dayCard(cv)
}

// playerCard is the portrait with its ring, and beside it the level, the
// name and age, and the XP bar with its numbers. It is w art pixels wide.
func (h *hud) playerCard(cv *canvas.Canvas) (width int) {
	f := h.f
	name := h.text(f.self.name+", "+i18n.T("%d years", f.self.who(f.clock.Day).Years), 'k')
	numbers := h.text(i18n.T("%d / %d XP", f.prog.XP, xpNeed(f.prog.Level)), 'k')
	left := faceD + 14 // art pixels from the card's edge to its texts
	w := left + max(name.Width/scale+8+levelW(h), xpBarW+6+numbers.Width/scale) + 6
	x, y := navSide*scale, navTop*scale
	cv.Blit(x, y, h.panel(w, navH, false))
	h.glowRing(cv, x+(4+faceD/2)*scale, y+(navH/2)*scale)
	cv.Blit(x+4*scale, y+(navH-faceD)/2*scale, h.portrait(faceD*scale))
	h.level(cv, x+left*scale, y+5*scale)
	cv.Blit(x+(left+levelW(h)+6)*scale, y+5*scale+(12*scale-name.Height)/2, name)
	h.xpBar(cv, x+left*scale, y+19*scale)
	cv.Blit(x+(left+xpBarW+6)*scale, y+19*scale+(9*scale-numbers.Height)/2, numbers)
	return w
}

// dayCard is the day and the hour, and under them the purse, at the right.
func (h *hud) dayCard(cv *canvas.Canvas) {
	f, c := h.f, h.f.clock
	ink := 'k'
	if c.Between(22, wakeAt) {
		ink = 'r'
	}
	day := h.text(i18n.T("Day %d  %02d:%02d", c.Day+1, int(c.Hour), int(c.Hour*60)%60), rune(ink))
	money := h.text("$"+itoa(f.money), 'k')
	w := max(day.Width, 20*scale+money.Width)/scale + 14
	x, y := cv.Width-(w+navSide)*scale, navTop*scale
	cv.Blit(x, y, h.panel(w, navH, false))
	cv.Blit(x+7*scale, y+5*scale+(12*scale-day.Height)/2, day)
	cv.Blit(x+3*scale, y+16*scale, h.icon("icon_coin"))
	cv.Blit(x+21*scale, y+16*scale+(16*scale-money.Height)/2, money)
}
