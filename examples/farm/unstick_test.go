package farm

import "testing"

// A game saved before the roads had fences may have left the player where
// a fence now stands: it still loads, and the player is not stuck in it.
func TestASavedPlayerOnANewFenceIsMovedOut(t *testing.T) {
	_, f := newGame(idle{})
	col, row := -1, -1
	for c := 33; c < 50 && col < 0; c++ {
		for _, r := range []int{pathRow - 1, pathRow + 1} {
			if f.thingAt(c, r) == tileID["fence"] {
				col, row = c, r
				break
			}
		}
	}
	if col < 0 {
		t.Fatal("no fence beside the road to the east town")
	}
	s := takeSnapshot(f)
	at := f.ground.CellCenter(col, row)
	s.At = point{X: at.X, Y: at.Y}
	_, g := newGame(idle{})
	if err := restore(g, &s); err != nil {
		t.Fatalf("a saved player on a fence: %v", err)
	}
	c, r := g.cellOf(g.player.Transform.Position)
	if !g.open(c, r) || max(abs(c-col), abs(r-row)) > 2 {
		t.Errorf("the player stands at %d,%d, open %v, after the fence at %d,%d", c, r, g.open(c, r), col, row)
	}
}
