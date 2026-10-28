package main

// talk is a conversation with a villager: the choice on its menu.
type talk struct {
	v      *villager
	choice int
}

// line is a choice in a conversation: its words, and what choosing it
// does.
type line struct {
	text string
	do   func(f *farm)
}

// The prices a lumberjack trades wood at with the player: they sell dear
// and buy cheap.
const (
	lumberSells = woodPrice + 1
	lumberBuys  = woodPrice - 1
)

// lines are what the player may say to v: deliver what they ask for,
// trade wood with a lumberjack, do the smith's errand or buy their axes,
// and say goodbye.
func (f *farm) lines(v *villager) []line {
	var ls []line
	if !v.smith() && len(v.order.Want) > 0 {
		ls = append(ls, line{text: "ENTREGAR PEDIDO", do: func(f *farm) { f.deliver(v) }})
	}
	switch v.role {
	case roleLumberjack:
		ls = append(ls,
			line{text: "COMPRAR MADEIRA $" + itoa(lumberSells), do: func(f *farm) { f.tradeWood(v.kit, f.kit, lumberSells, v.pos) }},
			line{text: "VENDER MADEIRA $" + itoa(lumberBuys), do: func(f *farm) { f.tradeWood(f.kit, v.kit, lumberBuys, v.pos) }})
	case roleSmith:
		if !f.errandDone {
			ls = append(ls, line{text: "ENTREGAR ENCOMENDA", do: func(f *farm) { f.doErrand(v) }})
			break
		}
		for i := 1; i < len(axes); i++ {
			ls = append(ls, line{text: "MACHADO DE " + axes[i].label + " $" + itoa(axes[i].price), do: func(f *farm) { f.buyAxe(i, v.pos) }})
		}
		ls = append(ls, line{text: "VENDER LENHA $" + itoa(woodPrice), do: func(f *farm) { f.tradeWood(f.kit, v.kit, woodPrice, v.pos) }})
	}
	return append(ls, line{text: "TCHAU", do: func(f *farm) { f.talk = nil; f.sound("ui_close", nil) }})
}
