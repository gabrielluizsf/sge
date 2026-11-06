package backpack

import "testing"

// shop is a Crafter of two recipes made from 3 wood: a chair takes 2, a
// table 5.
type shop struct{ wood, made int }

func (s *shop) Recipes() []Recipe {
	return []Recipe{
		{ID: "chair", Name: "Chair", Icon: "chair", Needs: []Need{{Name: "Wood", Icon: "wood", Want: 2, Have: s.wood}}},
		{ID: "table", Name: "Table", Icon: "table", Needs: []Need{{Name: "Wood", Icon: "wood", Want: 5, Have: s.wood}}},
	}
}

func (s *shop) Make(id string) bool {
	need := map[string]int{"chair": 2, "table": 5}[id]
	if need == 0 || s.wood < need {
		return false
	}
	s.wood -= need
	s.made++
	return true
}

func tabbed(wood int) (*Pack, *shop) {
	sh := &shop{wood: wood}
	p := opened()
	p.Tabs, p.Crafts, p.Rows = []string{"Backpack", "Craft"}, sh, 4
	p.Place(View{W: 1280, H: 720, Scale: 2})
	return p, sh
}

func TestTabKeysSwitchBetweenBackpackAndCraft(t *testing.T) {
	p, _ := tabbed(3)
	if p.Crafting() {
		t.Fatal("opens on the craft tab")
	}
	if ev := p.Handle(Input{Tab: 1}); ev != Tabbed || !p.Crafting() {
		t.Fatalf("event %v, crafting %v", ev, p.Crafting())
	}
	p.Handle(Input{Tab: 1}) // round to the first
	if p.Crafting() {
		t.Error("did not wrap round")
	}
}

func TestCraftingMakesWhatThereIsMaterialFor(t *testing.T) {
	p, sh := tabbed(3)
	p.Handle(Input{Tab: 1})
	if ev := p.Handle(Input{Select: true}); ev != Crafted || p.Made() != "chair" || sh.wood != 1 {
		t.Fatalf("chair: event %v made %q wood %d", ev, p.Made(), sh.wood)
	}
	if ev := p.Handle(Input{Select: true}); ev != Refused || sh.made != 1 {
		t.Errorf("a second chair from 1 wood: event %v, made %d", ev, sh.made)
	}
	p.Handle(Input{Dy: 1})
	if p.CraftCursor() != 1 {
		t.Fatalf("cursor %d", p.CraftCursor())
	}
	if ev := p.Handle(Input{Select: true}); ev != Refused || sh.wood != 1 {
		t.Errorf("table from 1 wood: %v wood %d", ev, sh.wood)
	}
}

func TestCraftingByMouseAndTheTabsByClick(t *testing.T) {
	p, sh := tabbed(9)
	d := painter{p: p, cv: blank(1280, 720), art: flat{}}
	p.Paint(d.cv, flat{})
	tab := p.rects.tabs[1]
	if ev := p.Handle(Input{Press: true, Mouse: Point{tab.X + 2, tab.Y + 2}}); ev != Tabbed || !p.Crafting() {
		t.Fatalf("click on the tab: %v", ev)
	}
	p.Paint(d.cv, flat{})
	row := p.rects.rows[1]
	if ev := p.Handle(Input{Press: true, Mouse: Point{row.X + 5, row.Y + 5}}); ev != Crafted || sh.wood != 4 {
		t.Errorf("click on the table row: %v wood %d", ev, sh.wood)
	}
}

func TestACraftListLongerThanTheTabScrolls(t *testing.T) {
	p, _ := tabbed(9)
	p.Rows = 1
	p.Handle(Input{Tab: 1})
	p.Handle(Input{Dy: 1})
	if p.craft.top != 1 || p.craft.cursor != 1 {
		t.Errorf("top %d cursor %d", p.craft.top, p.craft.cursor)
	}
}
