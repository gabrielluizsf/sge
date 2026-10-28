package main

import "github.com/gabrielluizsf/sge/town"

// shop is a stall's menu: buying either crop's seeds, or selling the
// harvest and wood; at is the stall's cell.
type shop struct {
	open   bool
	choice int
	at     town.Cell
}

// The menu's lines are a crop's seeds, by the crop's index, then
// selling.
const sellChoice = 2

// choose does the menu's chosen line for k at the stall.
func (f *farm) choose(k *kit, s *shop) {
	at := f.ground.CellCenter(s.at.Col, s.at.Row)
	if s.choice == sellChoice {
		f.sell(k, at)
		return
	}
	f.buy(k, s.choice, at)
}
