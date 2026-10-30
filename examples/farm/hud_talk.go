package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/inventory"
)

// talk is the conversation's panel, in the middle: who it is, what they
// want and what it pays, and the menu, the chosen line on a gold bar.
func (h *hud) talk(cv *canvas.Canvas) {
	f := h.f
	v := f.talk.v
	var say []string
	switch {
	case v.smith() && !f.errandDone:
		say = append(say, "Preciso disto para a forja:")
		say = append(say, wants(f, f.errand.Want)...)
		say = append(say, "Em troca: machado de pedra")
	case v.smith():
		say = append(say, "Machados melhores cortam mais rápido.")
	case len(v.order.Want) > 0:
		say = append(say, "Você pode me trazer:")
		say = append(say, wants(f, v.order.Want)...)
		say = append(say, fmt.Sprint("Pago $", v.order.Reward))
	default:
		say = append(say, "Oi! Hoje não preciso de nada.")
	}
	if v.role == roleLumberjack {
		say = append(say, fmt.Sprint("Tenho ", v.kit.bag.Count(wood), " de madeira"))
	}
	ls := f.lines(v)
	lines := append([]string{v.name + " - " + roleLabel(v)}, say...)
	lines = append(lines, "")
	for _, l := range ls {
		lines = append(lines, l.text)
	}
	lines = append(lines, "Espaço: escolher   E: sair")
	menu(cv, lines, 2+len(say)+f.talk.choice)
}

// wants are the lines of a wish: each item, how many, and how many the
// player has.
func wants(f *farm, l inventory.List) []string {
	var out []string
	for _, s := range l {
		out = append(out, fmt.Sprintf("   %d %s  (você tem %d)", s.N, itemLabels[s.Item], f.bag.Count(s.Item)))
	}
	return out
}
