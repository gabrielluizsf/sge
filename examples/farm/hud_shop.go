package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// shop is the stall's menu, in a panel in the middle: the seeds for sale
// and selling the harvest, the chosen line on a gold bar.
func (h *hud) shop(cv *canvas.Canvas) {
	choice := h.f.shop.choice
	first := min(max(0, choice-shopRows/2), len(crops)-shopRows)
	lines := []string{"Loja"}
	for _, c := range crops[first : first+shopRows] {
		lines = append(lines, fmt.Sprintf("Semente de %s $%d (%d dias)", c.label, c.seedPrice, c.daysToRipe))
	}
	lines = append(lines, fmt.Sprint("Vender colheita e madeira $", worth(h.f.kit)), "Espaço: escolher   E: sair")
	chosen := 1 + choice - first
	if choice == sellChoice {
		chosen = 1 + shopRows
	}
	h.menu(cv, menuOf{lines: lines, chosen: chosen})
}

// shopRows is how many crops' seeds the shop lists at once, scrolling to
// the one chosen.
const shopRows = 8

// menuOf is a menu's lines of text, the first its title, which of them
// is chosen, and how wide its panel is at least, in pixels of art.
type menuOf struct {
	lines  []string
	chosen int
	minW   int
}

// menu is a panel in the middle of the screen with a line of text a row,
// the row chosen on a gold bar; it is where the panel's corner is.
func (h *hud) menu(cv *canvas.Canvas, m menuOf) (int, int) {
	w := m.minW
	for _, l := range m.lines {
		w = max(w, h.text(l, 'k').Width/scale+20)
	}
	ht := rowH*len(m.lines) + 8
	x, y := (cv.Width-w*scale)/2, (cv.Height-ht*scale)/2
	cv.Blit(x, y, h.panel(w, ht, false))
	for i, l := range m.lines {
		ry := y + (4+i*rowH)*scale
		if i == m.chosen {
			cv.Blit(x+4*scale, ry, h.panel(w-8, rowH, true))
		}
		t := h.text(l, 'k')
		cv.Blit(x+10*scale, ry+(rowH*scale-t.Height)/2, t)
	}
	return x, y
}
