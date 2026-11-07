package farm

import (
	"github.com/gabrielluizsf/sge/bond"
	"github.com/gabrielluizsf/sge/i18n"
)

// bondLines are what the player may say to v as their bond grows: chat,
// give each gift they carry, ask a friend fond enough, grown up, not kin
// and free out, and ask the one they date, fond enough, to marry.
func (f *farm) bondLines(v *villager) []line {
	ls := []line{{text: i18n.T("Chat"), do: func(f *farm) { f.chat(v) }}}
	for _, item := range gifts {
		if n := f.bag.Count(item); n > 0 {
			ls = append(ls, line{text: i18n.T("Give %s (%d)", i18n.T(itemLabels[item]), n), do: func(f *farm) { f.gift(v, item) }})
		}
	}
	if v.ours {
		return ls
	}
	m := f.moment(v)
	free := v.partner == nil && !f.courting(v) && m.Grown && !m.Kin
	switch {
	case free && v.bond.Stage() == bond.Friends && v.bond.Affection >= bond.DateAt:
		ls = append(ls, line{text: i18n.T("Ask out"), do: func(f *farm) { f.askOut(v) }})
	case v.bond.Courted == bond.Dating && v.bond.Affection >= bond.ProposeAt:
		ls = append(ls, line{text: i18n.T("Propose"), do: func(f *farm) { f.propose(v) }})
	}
	return ls
}
