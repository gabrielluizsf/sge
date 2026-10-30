package farm

import "testing"

// A tree falls after as many strokes as the axe takes, gives wood, leaves
// a stump, and grows back; with no axe nothing happens.
func TestATreeIsFelledAndGrowsBack(t *testing.T) {
	g, f := newGame(idle{})
	col, row := 3, 2 // the farm's tree at T
	if f.thingAt(col, row) != tileID["tree0"] {
		t.Fatal("no tree at 3,2")
	}
	f.slot = slotAxe
	if f.use(f.kit, col, row) {
		t.Fatal("struck without an axe")
	}
	f.axe = 1 // iron: three strokes
	for i := range axes[1].hits {
		if !f.use(f.kit, col, row) {
			t.Fatalf("stroke %d missed", i)
		}
	}
	if f.thingAt(col, row) != tileID["stump"] || f.bag.Count(wood) < 1 || f.bag.Count(wood) > axes[1].most {
		t.Fatalf("not felled: %d, wood %d", f.thingAt(col, row), f.bag.Count(wood))
	}
	if f.use(f.kit, col, row) {
		t.Error("struck a stump")
	}
	for range int(regrowIn/frame) + 2 {
		g.Step(frame)
	}
	if f.thingAt(col, row) != tileID["tree0"] {
		t.Error("the tree did not grow back")
	}
}

// Trees of the wild fall too, and stay felled when their chunk is made
// again.
func TestAWildTreeIsFelled(t *testing.T) {
	_, f := newGame(idle{})
	col, row := -1, 0
	for c := 31; c < 60 && col < 0; c++ {
		for r := -10; r < 30; r++ {
			if f.thingAt(c, r) == tileID["tree0"] {
				col, row = c, r
				break
			}
		}
	}
	f.slot, f.axe = slotAxe, 3 // diamond: one stroke
	if col < 0 || !f.use(f.kit, col, row) || f.thingAt(col, row) != tileID["stump"] {
		t.Fatalf("the wild's tree at %d,%d did not fall", col, row)
	}
}
