package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/tilemap"
)

// rockCell is the cell of a rock in the wild, found by looking over the
// land beside the farm.
func rockCell(t *testing.T, f *farm) (col, row int) {
	t.Helper()
	for c := -120; c < 160; c++ {
		for r := -120; r < 160; r++ {
			if isRock(f.thingAt(c, r)) {
				return c, r
			}
		}
	}
	t.Fatal("no rock in the wild")
	return 0, 0
}

// The pickaxe is made in the Craft tab from wood alone, once, and put on
// the hotbar.
func TestThePickaxeIsCraftedFromWoodAndStocked(t *testing.T) {
	_, f := newGame(idle{})
	book := craftBook{f}
	if book.Make(itemPick) {
		t.Fatal("made a pickaxe from nothing")
	}
	f.bag.Add(wood, 20)
	if !book.Make(itemPick) || f.bag.Count(itemPick) != 1 || f.pack.Slot(itemPick) < 0 {
		t.Fatalf("pickaxe %d, on slot %d", f.bag.Count(itemPick), f.pack.Slot(itemPick))
	}
	if book.Make(itemPick) || f.bag.Count(itemPick) != 1 {
		t.Error("a second pickaxe")
	}
	if !f.hold(f.pack.Slot(itemPick)) || f.slot != slotPick {
		t.Error("the pickaxe is not taken in hand")
	}
}

// A rock breaks after its strokes, gives 2 to 4 stone and XP, is solid
// until then and grows back.
func TestARockIsMinedAndGrowsBack(t *testing.T) {
	g, f := newGame(idle{})
	col, row := rockCell(t, f)
	rock := f.thingAt(col, row)
	if !f.tiles.IsSolid(rock) {
		t.Error("a rock is not solid")
	}
	f.slot = slotPick
	if f.use(f.kit, col, row) {
		t.Fatal("struck without a pickaxe")
	}
	f.bag.Add(itemPick, 1)
	xp := f.xpValue()
	for i := range rockHits {
		if !f.use(f.kit, col, row) {
			t.Fatalf("stroke %d missed", i)
		}
	}
	if f.thingAt(col, row) != tilemap.Empty || f.bag.Count(stone) < stoneLeast || f.bag.Count(stone) > stoneMost {
		t.Fatalf("not broken: tile %d, stone %d", f.thingAt(col, row), f.bag.Count(stone))
	}
	if f.xpValue() <= xp {
		t.Error("mining gave no XP")
	}
	if f.use(f.kit, col, row) {
		t.Error("struck a rock that is gone")
	}
	for range int(rockRegrowIn/frame) + 2 {
		g.Step(frame)
	}
	if f.thingAt(col, row) != rock {
		t.Error("the rock did not grow back")
	}
}

// The axe does not break rocks, and E shows "mine" only with the pickaxe
// in hand by a rock.
func TestOnlyThePickaxeBreaksRocksAndEShowsIt(t *testing.T) {
	_, f := newGame(idle{})
	col, row := rockCell(t, f)
	f.axe, f.slot = 3, slotAxe
	if f.use(f.kit, col, row) || isRock(f.thingAt(col, row)) == false {
		t.Error("the axe struck a rock")
	}
	if f.marks(col, row) {
		t.Error("the axe marks a rock")
	}
	f.bag.Add(itemPick, 1)
	f.slot = slotPick
	if !f.marks(col, row) {
		t.Error("the pickaxe does not mark a rock")
	}
}

// Rocks lie only where nothing else does: not on the farm, in a town, on
// a road, in the water or on a tree.
func TestRocksAreKeptOffPathsHousesAndWater(t *testing.T) {
	_, f := newGame(idle{})
	n := 0
	for c := -60; c < 200; c++ {
		for r := -60; r < 200; r++ {
			if !isRock(f.thingAt(c, r)) {
				continue
			}
			n++
			if f.onFarm(c, r) || f.towns.Grounds(c, r, 0) || f.wild.TileAt(wildWater, c, r) != tilemap.Empty || f.wild.TileAt(wildGround, c, r) == tileID["path"] {
				t.Fatalf("a rock at %d,%d is where it should not be", c, r)
			}
		}
	}
	if n == 0 {
		t.Error("no rocks about")
	}
}

// A save with a rock mined and growing back loads back with the rock
// gone; a save that says a rock was mined where there is none, or being
// struck where there is nothing, is refused.
func TestMinedRocksAreSavedAndHostileOnesRefused(t *testing.T) {
	_, f := newGame(idle{})
	col, row := rockCell(t, f)
	f.bag.Add(itemPick, 1)
	f.slot = slotPick
	for range rockHits {
		f.use(f.kit, col, row)
	}
	snap := takeSnapshot(f)
	_, g := newGame(idle{})
	if err := restore(g, &snap); err != nil {
		t.Fatal(err)
	}
	if g.thingAt(col, row) != tilemap.Empty || len(g.stumps) != 1 {
		t.Errorf("the mined rock came back at once: tile %d, %d stumps", g.thingAt(col, row), len(g.stumps))
	}
	rock := f.stumps[0].tree
	hostile(t, "rock mined in the open", func(s *snapshot) { s.Stumps = []stumpState{{Col: 0, Row: 0, Tree: rock, Left: 1}} })
	hostile(t, "rock time", func(s *snapshot) { s.Stumps = []stumpState{{Col: col, Row: row, Tree: rock, Left: 1e9}} })
	hostile(t, "rock struck in the open", func(s *snapshot) { s.Hits = []hitState{{Col: 0, Row: 0, N: 2}} })
	hostile(t, "rock struck too often", func(s *snapshot) { s.Hits = []hitState{{Col: col, Row: row, N: rockHits}} })
	hostile(t, "pickaxe on the hotbar without one", func(s *snapshot) {
		s.Player.Bag = nil
		s.Hotbar[5] = itemPick
	})
}
