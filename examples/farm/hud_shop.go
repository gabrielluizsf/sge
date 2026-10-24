package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// shop is the stall's menu, in a panel in the middle: the seeds for sale
// and selling the harvest, the chosen line on a yellow bar.
func (h *hud) shop(cv *canvas.Canvas) {
	const w, ht = 150, 80
	x, y := (cv.Width-w*scale)/2, (cv.Height-ht*scale)/2
	cv.Blit(x, y, panel(w, ht, false))
	line := func(i int, s string, ink rune) {
		if ink == 'y' {
			cv.Blit(x+4*scale, y+(4+i*12)*scale, panel(w-8, 12, true))
			ink = 'k'
		}
		cv.Blit(x+8*scale, y+(6+i*12)*scale, text(s, ink))
	}
	line(0, "LOJA", 'k')
	for i, c := range crops {
		line(1+i, fmt.Sprintf("SEMENTE %s $%d", c.label, c.seedPrice), h.ink(i))
	}
	worth := 0
	for c, n := range h.f.harvest {
		worth += n * crops[c].sellPrice
	}
	line(1+sellChoice, fmt.Sprint("VENDER COLHEITA $", worth), h.ink(sellChoice))
	line(5, "ESPACO: OK  E: SAIR", 'k')
}

// ink is the colour of the shop's line i: yellow, for the bar, when
// chosen.
func (h *hud) ink(i int) rune {
	if i == h.f.shop.choice {
		return 'y'
	}
	return 'w'
}
