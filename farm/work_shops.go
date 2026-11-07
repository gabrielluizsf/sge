package farm

import (
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/pathfind"
)

// A baker burns a load of wood for each batch, which earns them breadPay;
// a butcher pays the dealer animalPay for an animal.
const (
	breadPay  = 8
	animalPay = 30
)

// bakerWood is a baker's seeing to their oven: baking while it burns, and
// when it is out, buying wood from a lumberjack, as the smith does, and not
// felling trees of their own. Nil with nothing to do but stand at the shop.
func (m *mind) bakerWood(f *farm) *job {
	v := m.v
	k := v.kit
	r := shedRect(f.art, v.town, 1)
	at := pathfind.Cell{Col: r.col + r.cols/2, Row: r.row + r.rows}
	switch {
	case v.fuel >= 1:
		return &job{to: at, face: "up", slot: -1, rest: 3, then: func(*farm) {
			v.fuel--
			k.money += breadPay
		}}
	case k.bag.Count(wood) > 0:
		return &job{to: at, face: "up", slot: -1, rest: 1, then: func(f *farm) {
			n := k.bag.Count(wood)
			k.bag.Take(wood, n)
			v.fuel += float64(n)
			f.float(f.center(at), i18n.T("+%d wood", n), 'y')
		}}
	}
	if lj := f.lumberjackWith(v.town, 2); lj != nil && k.money >= 2*woodPrice {
		col, row := f.cellOf(lj.pos)
		return &job{to: pathfind.Cell{Col: col, Row: row}, slot: -1, rest: 1, then: func(f *farm) { f.buyWood(v, lj) }}
	}
	return nil
}

// butcherBuys is a butcher's walk to the stable to buy an animal from its
// dealer, if there is one and they have the money. Its pen is stocked anew.
func (m *mind) butcherBuys(f *farm) *job {
	v := m.v
	if v.kit.money < animalPay {
		return nil
	}
	for _, d := range f.people {
		if d.town != v.town || d.role != roleDealer || d.inside {
			continue
		}
		r := stableRect(f.art, v.town)
		return &job{to: pathfind.Cell{Col: r.col + r.cols/2, Row: r.row + r.rows}, slot: -1, rest: 2, then: func(*farm) {
			if v.kit.money >= animalPay {
				v.kit.money -= animalPay
				d.kit.money += animalPay
			}
		}}
	}
	return nil
}
