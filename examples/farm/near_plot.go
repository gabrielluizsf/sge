package farm

import "github.com/gabrielluizsf/sge/game"

// nearReach is how far, in tiles, a planted plot may be from the farmer's
// feet for its note to show.
const nearReach = 1.5

// nearPlot is the planted plot closest to the farmer within nearReach, or
// nil. It looks only at the cells around the farmer.
func (f *farm) nearPlot(o *game.Object) *plot {
	at := feet(o)
	col, row := f.cellOf(at)
	var best *plot
	bestD := nearReach * nearReach * tile * tile
	for dr := -2; dr <= 2; dr++ {
		for dc := -2; dc <= 2; dc++ {
			p := f.plotAt(col+dc, row+dr)
			if p == nil || p.crop < 0 {
				continue
			}
			c := f.ground.CellCenter(p.col, p.row)
			dx, dy := c.X-at.X, c.Y-at.Y
			if d := dx*dx + dy*dy; d <= bestD {
				best, bestD = p, d
			}
		}
	}
	return best
}
