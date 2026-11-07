package farm

import (
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/town"
)

// settle peoples town t the first time it is woken: a plot in each
// house's garden, and in each house a household of a man, a woman and
// their child, the smith's family in the first house, the smith working
// the anvil; each of an age, the couple of one trade or the other, and
// the child young. Its seed decides who, so a town is the same people
// every game.
func (f *farm) settle(t *town.Town) {
	if f.settled[t] {
		return
	}
	f.settled[t] = true
	s := &settler{f: f, t: t, src: random.New(t.Seed)}
	// The first house is the smith's, the second a lumberjack's; the rest
	// are farmers', lumberjacks' and wanderers' as the seed falls.
	trades := []string{roleFarmer, roleFarmer, roleWanderer, roleLumberjack}
	var homes []household
	for i := range t.Houses {
		h := &t.Houses[i]
		f.plots = append(f.plots, &plot{col: h.Garden.Col, row: h.Garden.Row, crop: -1, town: t})
		role := roleLumberjack
		switch i {
		case 0:
			role = roleSmith
		case 1:
		default:
			role = trades[s.src.Intn(len(trades))]
		}
		mate := ""
		switch {
		case i > 0 && i == min(2, len(t.Houses)-1): // the dealer's mate bakes
			role, mate = roleDealer, roleBaker
		case i == 1: // the lumberjack's butchers
			mate = roleButcher
		}
		hh := s.household(h, role, mate)
		homes = append(homes, hh)
		f.people = append(f.people, hh.man, hh.woman, hh.child)
	}
	f.populate(t, homes)
}
