package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// shop is the stall's menu, in a panel in the middle: the seeds for sale
// and selling the harvest, the chosen line on a gold bar.
func (h *hud) shop(cv *canvas.Canvas) {
	lines := []string{"Loja"}
	for _, c := range crops {
		lines = append(lines, fmt.Sprintf("Semente de %s $%d", c.label, c.seedPrice))
	}
	lines = append(lines, fmt.Sprint("Vender colheita e madeira $", worth(h.f.kit)), "Espaço: escolher   E: sair")
	menu(cv, lines, 1+h.f.shop.choice)
}

// menu is a panel in the middle of the screen with a line of text a row,
// the first its title, the row chosen on a gold bar.
func menu(cv *canvas.Canvas, lines []string, chosen int) {
	w := 0
	for _, l := range lines {
		w = max(w, text(l, 'k').Width/scale+20)
	}
	ht := rowH*len(lines) + 8
	x, y := (cv.Width-w*scale)/2, (cv.Height-ht*scale)/2
	cv.Blit(x, y, panel(w, ht, false))
	for i, l := range lines {
		ry := y + (4+i*rowH)*scale
		if i == chosen {
			cv.Blit(x+4*scale, ry, panel(w-8, rowH, true))
		}
		t := text(l, 'k')
		cv.Blit(x+10*scale, ry+(rowH*scale-t.Height)/2, t)
	}
}
