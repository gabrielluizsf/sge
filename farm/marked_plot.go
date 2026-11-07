package farm

import "github.com/gabrielluizsf/sge/game"

// markedPlot is the planted plot the cursor marks, the cell before the
// farmer o or the one the mouse is on, or nil when it marks none: indoors,
// swimming, or on a cell with nothing sown.
func (f *farm) markedPlot(o *game.Object) *plot {
	if f.room != nil || f.hands.wet {
		return nil
	}
	p := f.plotAt(f.hands.target(o))
	if p == nil || p.crop < 0 {
		return nil
	}
	return p
}
