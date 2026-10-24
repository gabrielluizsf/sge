package main

// shop is the stall's menu: buying either crop's seeds, or selling the
// harvest.
type shop struct {
	open   bool
	choice int
}

// The menu's lines are a crop's seeds, by the crop's index, then
// selling.
const sellChoice = 2

// choose does the menu's chosen line for the farmer at the stall.
func (f *farm) choose(s *shop) {
	at := f.ground.CellCenter(stallCol, stallRow)
	if s.choice == sellChoice {
		f.sell(at)
		return
	}
	f.buy(s.choice, at)
}
