package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/tilemap"
)

// The farm starts with a cow in a fenced pen, which wanders inside it,
// and the pen has a gate the player can go in by.
func TestTheFarmStartsWithACowInAPen(t *testing.T) {
	g, f := newGame(idle{})
	f.player.Transform.Position = f.ground.CellCenter(penGate.Col, penGate.Row+2)
	g.Step(frame)
	cow := g.Find("animal")
	if cow == nil {
		t.Fatal("no animal on the farm")
	}
	field := penRect.inside()
	lo, hi := float64(field.col)*tile, float64(field.col+field.cols)*tile
	top, bottom := float64(field.row)*tile, float64(field.row+field.rows)*tile
	moved := false
	start := cow.Transform.Position
	for range 60 * 60 {
		g.Step(frame)
		p := cow.Transform.Position
		if p.X < lo || p.X > hi || p.Y < top || p.Y > bottom {
			t.Fatalf("the cow left its pen: %v", p)
		}
		moved = moved || p != start
	}
	if !moved {
		t.Error("the cow never moved")
	}
	for col := penRect.col; col < penRect.col+penRect.cols; col++ {
		for _, row := range []int{penRect.row, penRect.row + penRect.rows - 1} {
			solid := f.things.At(col, row) != tilemap.Empty && isFence(f.things.At(col, row))
			if gate := col == penGate.Col && row == penGate.Row; solid == gate {
				t.Fatalf("cell %d,%d: fence %v, gate %v", col, row, solid, gate)
			}
		}
	}
}

// Two farms with the same seed keep the same cow.
func TestTheCowWandersTheSameEveryTime(t *testing.T) {
	var at [2]float64
	for i := range at {
		g, f := newGame(idle{})
		f.player.Transform.Position = f.ground.CellCenter(penGate.Col, penGate.Row+2)
		for range 600 {
			g.Step(frame)
		}
		at[i] = g.Find("animal").Transform.Position.X
	}
	if at[0] != at[1] {
		t.Fatalf("the cow was at %v, then %v", at[0], at[1])
	}
}
