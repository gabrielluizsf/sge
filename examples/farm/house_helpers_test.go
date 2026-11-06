package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/game"
)

// longBed is the stylized art that says its bed covers two cells, as the
// realistic art's does, so the tests can arrange a piece larger than a cell.
type longBed struct{ Art }

func (longBed) Furniture(name string) (int, int) {
	if name == "bed" {
		return 2, 1
	}
	return 1, 1
}

// inHouse is a farm with the player inside their own house, the bed alone
// in it, and a bag of furniture to put down.
func inHouse(t *testing.T) (*game.Game, *farm) {
	t.Helper()
	g, f := build(longBed{SVGArt()}, idle{})
	f.enter(houseDoor{at: homeDoor, faceOut: "down", seed: seed, own: true})
	g.Step(frame)
	f.player.Transform.Position = f.ground.CellCenter(entry[0], entry[1]) // by the door
	for _, n := range furnishings {
		f.bag.Add(n, 3)
	}
	return g, f
}

// arrange puts the cursor on cell col, row, to put down or lift there.
func arrange(f *farm, col, row int) {
	if f.edit == nil {
		f.startEdit("")
	}
	f.edit.col, f.edit.row = col, row
}

// placed reports whether a piece called name stands on cell col, row.
func placed(f *farm, name string, col, row int) bool {
	i := f.pieceAt(col, row)
	return i >= 0 && f.home[i].name == name
}
