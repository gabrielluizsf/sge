package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/game"
)

// pressing is a controller standing still, facing one way, pressing use.
type pressing string

func (p pressing) control(*farm, *game.Object, float64) intent {
	return intent{face: string(p), use: true}
}

// The can fills at the water's edge, the farmer standing in the shallows
// by the shore as they must to reach it, and at a town's well; full, it
// says so.
func TestTheCanFillsAtTheShoreAndTheWell(t *testing.T) {
	g, f := newGame(pressing("up"))
	o := g.Find("farmer")
	f.kit.slot = slotCan
	col, row := shallowCell(t, f)
	o.Transform.Position = f.ground.CellCenter(col, row)
	f.kit.used = 7
	for range 3 {
		g.Step(frame)
	}
	if f.kit.used != 0 {
		t.Fatalf("standing in the shallows, the can was not filled: %d used", f.kit.used)
	}
	if !floated(f, "Regador cheio") {
		t.Error("filling the can says nothing")
	}
	if col, row := f.hands.useCell(o); !f.marks(col, row) {
		t.Error("the cursor does not mark the water the can is filled at")
	}

	tn := eastTown(f)
	w := wellRect(f.art, tn)
	o.Transform.Position = standIn(f, w.col, w.row+max(w.rows, 1))
	f.hands.acting = 0
	f.kit.used = 7
	for range 3 {
		g.Step(frame)
	}
	if f.kit.used != 0 {
		t.Fatalf("at the well, the can was not filled: %d used", f.kit.used)
	}
}

// shallowCell is a cell of the farm's pond the farmer wades in, not deep.
func shallowCell(t *testing.T, f *farm) (int, int) {
	for row := 0; row < f.water.Rows; row++ {
		for col := 0; col < f.water.Cols; col++ {
			at := f.ground.CellCenter(col, row)
			if f.isWater(col, row) && !f.deepAt(at) {
				return col, row
			}
		}
	}
	t.Fatal("no shallow water on the farm")
	return 0, 0
}

// floated reports whether a word floating over the world says text.
func floated(f *farm, text string) bool {
	for _, w := range f.floats {
		if w.text == text {
			return true
		}
	}
	return false
}

// Holding the can, the cursor marks the water before the farmer, where it
// fills.
func TestTheCursorMarksWaterForTheCan(t *testing.T) {
	_, f := newGame(idle{})
	col, row := waterCell(f)
	f.kit.slot = slotCan
	if !f.marks(col, row) {
		t.Error("the cursor does not mark water with the can in hand")
	}
	f.kit.slot = slotHoe
	if f.marks(col, row) {
		t.Error("the cursor marks water with a hoe in hand")
	}
}
