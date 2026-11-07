package farm

import (
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/town"
)

// place is the town v is in now: the one they visit, or their own.
func (v *villager) place() *town.Town {
	if v.guest != nil {
		return v.guest
	}
	return v.town
}

// pathCell is town cell c as a cell of the pathfinder's.
func pathCell(c town.Cell) pathfind.Cell { return pathfind.Cell(c) }
