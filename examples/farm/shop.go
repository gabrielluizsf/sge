package farm

import "github.com/gabrielluizsf/sge/town"

// shop is a stall's menu: buying any crop's seeds, or selling the
// harvest and wood; at is the stall's cell.
type shop struct {
	open   bool
	choice int
	at     town.Cell
	// animals says it is a stable's menu, of animals for sale.
	animals bool
}

// The menu's lines are a crop's seeds, by the crop's index, then
// selling, which is the last.
var sellChoice = len(crops)

// choose does the menu's chosen line for k at the stall.
func (f *farm) choose(k *kit, s *shop) {
	if s.animals {
		f.buyAnimal(k, stableOffers[s.choice])
		return
	}
	at := f.ground.CellCenter(s.at.Col, s.at.Row)
	if s.choice == sellChoice {
		f.sell(k, at)
		return
	}
	f.buy(k, s.choice, at)
}
