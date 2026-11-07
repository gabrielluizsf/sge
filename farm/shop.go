package farm

import "github.com/gabrielluizsf/sge/town"

// shop is a stall's menu: buying any crop's seeds, or selling the
// harvest and wood; at is the stall's cell.
type shop struct {
	open   bool
	choice int
	at     town.Cell
	// kind is the shed it is the menu of, as sheds names them; none for a
	// stall.
	kind string
}

// The menu's lines are a crop's seeds, by the crop's index, then
// selling, which is the last.
var sellChoice = len(crops)

// choose does the menu's chosen line for k at the stall.
func (f *farm) choose(k *kit, s *shop) {
	if s.kind != "" {
		if _, lines := f.shedLines(s.kind); s.choice < len(lines) {
			lines[s.choice].do(f, k)
		}
		return
	}
	at := f.ground.CellCenter(s.at.Col, s.at.Row)
	if s.choice == sellChoice {
		f.sell(k, at)
		return
	}
	f.buy(k, s.choice, at)
}
