package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/backpack"
)

// The Craft tab lists every recipe with the materials wanted and owned,
// and making one takes the materials and gives the item; lacking a
// material refuses and takes nothing.
func TestCraftingConsumesMaterialsAndRefusesWithoutThem(t *testing.T) {
	_, f := newGame(idle{})
	book := craftBook{f}
	if len(book.Recipes()) != len(recipes) {
		t.Fatalf("%d recipes", len(book.Recipes()))
	}
	if book.Make("chair") || f.bag.Count("chair") != 0 {
		t.Fatal("made a chair from nothing")
	}
	f.bag.Add(wood, 4)
	if !book.Make("chair") || f.bag.Count("chair") != 1 || f.bag.Count(wood) != 1 {
		t.Fatalf("chair %d, wood %d", f.bag.Count("chair"), f.bag.Count(wood))
	}
	if book.Make("chair") || f.bag.Count(wood) != 1 {
		t.Error("a second chair from 1 wood")
	}
	if book.Make("stove") || f.bag.Count(wood) != 1 {
		t.Error("a stove without stone took wood")
	}
	f.bag.Add(wood, 4)
	f.bag.Add(stone, 10)
	if !book.Make("stove") || f.bag.Count(stone) != 0 {
		t.Error("no stove from the wood and stone")
	}
	if book.Make("sofa") {
		t.Error("made what has no recipe")
	}
}

// Every recipe makes a piece of furniture that is an item with an icon, a
// name and a tooltip, and the list says what each lacks.
func TestEveryRecipeIsAnItemWithATooltip(t *testing.T) {
	_, f := newGame(idle{})
	f.bag.Add(wood, 3)
	for _, r := range (craftBook{f}).Recipes() {
		if !isFurnishing(r.ID) && r.ID != itemPick || r.Name == "" || r.Icon == "" || len(r.Info) == 0 || len(r.Needs) == 0 {
			t.Errorf("recipe %+v", r)
		}
		for _, n := range r.Needs {
			if n.Name == "" || n.Icon == "" || n.Want <= 0 {
				t.Errorf("%s: need %+v", r.ID, n)
			}
		}
	}
	if r := (craftBook{f}).Recipes()[1]; r.ID != "chair" || !r.Can() || r.Needs[0].Have != 3 {
		t.Errorf("chair: %+v", r)
	}
}

// The inventory screen has the Craft tab, and Enter on it crafts.
func TestTheInventoryCraftsFromItsTab(t *testing.T) {
	_, f := newGame(idle{})
	f.bag.Add(wood, 3)
	f.pack.Show()
	f.pack.Handle(backpack.Input{Tab: 1})
	f.pack.Handle(backpackDown) // past the pickaxe, to the chair
	if ev := f.pack.Handle(backpack.Input{Select: true}); ev != backpack.Crafted || f.bag.Count("chair") != 1 {
		t.Fatalf("event %v, chairs %d", ev, f.bag.Count("chair"))
	}
	f.pack.Handle(backpack.Input{Tab: 1})
	found := false
	for _, it := range (carried{f}).Items() {
		found = found || it.ID == "chair" && it.Count == 1
	}
	if !found {
		t.Error("the chair is not in the backpack")
	}
}

// The player's own house has only the bed to begin with.
func TestTheHouseStartsWithOnlyTheBed(t *testing.T) {
	g, f := newGame(idle{})
	if len(f.home) != 1 || f.home[0].name != "bed" {
		t.Fatalf("home %v", f.home)
	}
	f.enter(houseDoor{at: homeDoor, faceOut: "down", seed: seed, own: true})
	g.Step(frame)
	furniture := 0
	for r := range roomRows {
		for c := range roomCols {
			for _, n := range append(append([]string{}, furnishings...), "bed_blue") {
				if id := f.room.things.At(c, r); id != -1 && id == tileID[pieceTileName(n, 0, 0, 0)] {
					furniture++
				}
			}
		}
	}
	if furniture != 1 {
		t.Errorf("%d pieces in the room", furniture)
	}
}
