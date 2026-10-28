package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/inventory"
)

// talk is the conversation's panel, in the middle: who it is, what they
// want and what it pays, and the menu, the chosen line on a yellow bar.
func (h *hud) talk(cv *canvas.Canvas) {
	f := h.f
	v := f.talk.v
	var say []string
	switch {
	case v.smith() && !f.errandDone:
		say = append(say, "PRECISO DISTO PARA A FORJA:")
		say = append(say, wants(f, f.errand.Want)...)
		say = append(say, "GANHE: MACHADO DE PEDRA")
	case v.smith():
		say = append(say, "MACHADOS MELHORES CORTAM MAIS RAPIDO")
	case len(v.order.Want) > 0:
		say = append(say, "VOCE PODE ME TRAZER:")
		say = append(say, wants(f, v.order.Want)...)
		say = append(say, fmt.Sprint("PAGO $", v.order.Reward))
	default:
		say = append(say, "OI! HOJE NAO PRECISO DE NADA")
	}
	if v.role == roleLumberjack {
		say = append(say, fmt.Sprint("TENHO ", v.kit.bag.Count(wood), " MADEIRA"))
	}
	ls := f.lines(v)
	const w = 236
	ht := 12*(len(say)+len(ls)+3) + 4
	x, y := (cv.Width-w*scale)/2, (cv.Height-ht*scale)/2
	cv.Blit(x, y, panel(w, ht, false))
	row := func(i int, s string, ink rune) {
		if ink == 'y' {
			cv.Blit(x+4*scale, y+(4+i*12)*scale, panel(w-8, 12, true))
			ink = 'k'
		}
		cv.Blit(x+8*scale, y+(6+i*12)*scale, text(s, ink))
	}
	row(0, v.name+" - "+roleLabel(v), 'k')
	for i, s := range say {
		row(1+i, s, 'w')
	}
	for i, l := range ls {
		ink := 'w'
		if i == f.talk.choice {
			ink = 'y'
		}
		row(2+len(say)+i, l.text, ink)
	}
	row(2+len(say)+len(ls), "ESPACO: OK  E: SAIR", 'k')
}

// wants are the lines of a wish: each item, how many, and how many the
// player has.
func wants(f *farm, l inventory.List) []string {
	var out []string
	for _, s := range l {
		out = append(out, fmt.Sprintf("  %d %s - VOCE TEM %d", s.N, itemLabels[s.Item], f.bag.Count(s.Item)))
	}
	return out
}
