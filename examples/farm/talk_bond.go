package farm

import "github.com/gabrielluizsf/sge/bond"

// bondLines are what the player may say to v as their bond grows: chat,
// give each gift they carry, ask a friend fond enough, grown up, not kin
// and free out, and ask the one they date, fond enough, to marry.
func (f *farm) bondLines(v *villager) []line {
	ls := []line{{text: "Conversar", do: func(f *farm) { f.chat(v) }}}
	for _, item := range gifts {
		if n := f.bag.Count(item); n > 0 {
			ls = append(ls, line{text: "Dar " + itemLabels[item] + " (" + itoa(n) + ")", do: func(f *farm) { f.gift(v, item) }})
		}
	}
	if v.ours {
		return ls
	}
	m := f.moment(v)
	free := v.partner == nil && !f.courting(v) && m.Grown && !m.Kin
	switch {
	case free && v.bond.Stage() == bond.Friends && v.bond.Affection >= bond.DateAt:
		ls = append(ls, line{text: "Pedir em namoro", do: func(f *farm) { f.askOut(v) }})
	case v.bond.Courted == bond.Dating && v.bond.Affection >= bond.ProposeAt:
		ls = append(ls, line{text: "Pedir em casamento", do: func(f *farm) { f.propose(v) }})
	}
	return ls
}
