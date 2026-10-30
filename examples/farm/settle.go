package farm

import (
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/town"
)

// settle peoples town t the first time it is woken: a plot in each
// house's garden, a villager of some trade in each house, a woman or a
// man of any skin, the smith in the first, working the anvil. Its seed decides who, so a town is the same
// people every game.
func (f *farm) settle(t *town.Town) {
	if f.settled[t] {
		return
	}
	f.settled[t] = true
	src := random.New(t.Seed)
	person := func(home *town.House, role string, stand town.Cell) *villager {
		woman := src.Float() < 0.5
		names := menNames
		if woman {
			names = womenNames
		}
		v := &villager{
			name: names[src.Intn(len(names))], woman: woman, role: role,
			who:  Who{Woman: woman, Role: role, Genes: [5]int{src.Intn(64), src.Intn(64), src.Intn(64), src.Intn(64), src.Intn(64)}},
			kit:  &kit{money: 15, axe: 0, slot: slotCan},
			town: t, home: home,
			pos:  f.ground.CellCenter(stand.Col, stand.Row),
			luck: random.New(uint32(src.Intn(1 << 30))),
			wait: float64(src.Intn(orderPause)),
		}
		v.kit.bag.Add(seedsOf(src.Intn(len(crops))), 2)
		v.mind = &mind{v: v}
		return v
	}
	// The first house is the smith's, the second a lumberjack's; the rest
	// are farmers', lumberjacks' and wanderers' as the seed falls.
	var smithHome *town.House
	trades := []string{roleFarmer, roleFarmer, roleWanderer, roleLumberjack}
	for i := range t.Houses {
		h := &t.Houses[i]
		f.plots = append(f.plots, &plot{col: h.Garden.Col, row: h.Garden.Row, crop: -1})
		switch i {
		case 0:
			smithHome = h
			continue
		case 1:
			f.people = append(f.people, person(h, roleLumberjack, h.Door))
		default:
			f.people = append(f.people, person(h, trades[src.Intn(len(trades))], h.Door))
		}
	}
	smith := person(smithHome, roleSmith, town.Cell{Col: t.Anvil.Col, Row: t.Anvil.Row + 1})
	smith.kit.money, smith.fuel = 40, 2
	f.people = append(f.people, smith)
}
