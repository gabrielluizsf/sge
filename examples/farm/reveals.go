package farm

import (
	"github.com/gabrielluizsf/sge/game"
)

// reveals are what the farm keeps in sight: the farmer, and the cell they
// work, so a tree or a roof drawn over either fades to show them.
func (f *farm) reveals() []game.Reveal {
	p := f.player
	if p == nil || p.Hidden {
		return nil
	}
	height := 2 * tile
	if s := p.Look(); s != nil && s.Image != nil {
		height = s.Image.Height
	}
	keep := []game.Reveal{{At: feet(p), Height: height}}
	if f.room == nil && f.hands != nil {
		col, row := f.hands.target(p)
		keep = append(keep, game.Reveal{At: f.ground.CellCenter(col, row), Height: tile})
	}
	return keep
}
