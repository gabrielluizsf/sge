package farm

import "github.com/gabrielluizsf/sge/game"

// useCell is the cell the farmer o uses what they hold on: the one the
// cursor marks; but wading, can in hand, before no water, the water they
// stand in, which the can is filled from.
func (p *farmer) useCell(o *game.Object) (col, row int) {
	col, row = p.target(o)
	if p.wet && p.kit.slot == slotCan && !p.f.isWater(col, row) {
		return p.f.cellOf(feet(o))
	}
	return col, row
}
