package farm

import (
	"encoding/json"
	"reflect"
	"testing"
)

// furnished is a farm whose house has been arranged, with a chest of
// something in it, a chest out of doors, and furniture in the bag.
func furnished(t *testing.T) *farm {
	t.Helper()
	_, f := inHouse(t)
	f.home = append(f.home, newPiece("chair", 4, 2, 3), newPiece("chest", 6, 3, 1), newPiece("bed", 1, 3, 1))
	f.home[2].held[0] = chestSlot{item: wood, n: 7}
	f.home[2].held[9] = chestSlot{item: seedsOf(1), n: 2}
	f.leave()
	col, row := f.cellOf(feet(f.player))
	for dc := 1; dc < 4; dc++ {
		if p := newPiece("chest", col+dc, row, 2); f.canPlaceYard(p) {
			p.held[4] = chestSlot{item: stone, n: 3}
			f.putChest(p)
			break
		}
	}
	if len(f.yard) != 1 {
		t.Fatal("no chest out of doors")
	}
	return f
}

// The furniture, the bag's pieces and the chests with what they keep, in
// the house and out of doors, come back as they were.
func TestFurnitureAndChestsSurviveSaving(t *testing.T) {
	f := furnished(t)
	f.bag.Add("table", 2)
	a := takeSnapshot(f)
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var c snapshot
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatal(err)
	}
	_, g2 := build(longBed{SVGArt()}, idle{})
	if err := restore(g2, &c); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g2.home, f.home) {
		t.Errorf("home\n%+v\n%+v", g2.home, f.home)
	}
	if len(g2.yard) != 1 || !reflect.DeepEqual(g2.yard[0].p, f.yard[0].p) || g2.bag.Count("table") != 5 {
		t.Errorf("yard %+v", g2.yard)
	}
	if !reflect.DeepEqual(a, takeSnapshot(g2)) {
		t.Error("the snapshot changed in the round trip")
	}
	if g2.player.Game().Find("chest") == nil {
		t.Error("the chest out of doors is not in the world")
	}
}

// A save made before furniture existed has none of it: its house has only
// the bed, and it loads.
func TestAnOldSaveLoadsWithTheBedOnly(t *testing.T) {
	f := furnished(t)
	s := takeSnapshot(f)
	var raw map[string]any
	b, _ := json.Marshal(s)
	json.Unmarshal(b, &raw)
	delete(raw, "home")
	delete(raw, "yard")
	b, _ = json.Marshal(raw)
	var old snapshot
	if err := json.Unmarshal(b, &old); err != nil {
		t.Fatal(err)
	}
	_, g2 := build(longBed{SVGArt()}, idle{})
	if err := restore(g2, &old); err != nil {
		t.Fatal(err)
	}
	if len(g2.home) != 1 || g2.home[0].name != "bed" || len(g2.yard) != 0 {
		t.Errorf("home %+v yard %v", g2.home, g2.yard)
	}
}

// A hostile save of furniture is refused and leaves the game as it was.
func TestHostileFurnitureIsRefused(t *testing.T) {
	good := takeSnapshot(furnished(t))
	b, _ := json.Marshal(good)
	cases := map[string]func(*snapshot){
		"unknown piece":        func(s *snapshot) { s.Home[1].Name = "throne" },
		"heading":              func(s *snapshot) { s.Home[1].Heading = 4 },
		"negative heading":     func(s *snapshot) { s.Home[1].Heading = -1 },
		"outside the walls":    func(s *snapshot) { s.Home[1].Col = 0 },
		"far away":             func(s *snapshot) { s.Home[1].Row = 1 << 40 },
		"on the door":          func(s *snapshot) { s.Home[1].Col, s.Home[1].Row = entry[0], entry[1] },
		"on another piece":     func(s *snapshot) { s.Home[1].Col, s.Home[1].Row = s.Home[0].Col, s.Home[0].Row },
		"no bed":               func(s *snapshot) { s.Home = s.Home[1:2] },
		"too many pieces":      func(s *snapshot) { s.Home = append(s.Home, make([]pieceState, 100)...) },
		"a chair that keeps":   func(s *snapshot) { s.Home[1].Held = []slotState{{Item: wood, N: 1}} },
		"a chest of 11":        func(s *snapshot) { s.Home[2].Held = append(s.Home[2].Held, slotState{}) },
		"a chest of 3":         func(s *snapshot) { s.Home[2].Held = s.Home[2].Held[:3] },
		"unknown item":         func(s *snapshot) { s.Home[2].Held[1] = slotState{Item: "gold", N: 1} },
		"negative count":       func(s *snapshot) { s.Home[2].Held[0].N = -3 },
		"huge count":           func(s *snapshot) { s.Home[2].Held[0].N = 1 << 40 },
		"item with no count":   func(s *snapshot) { s.Home[2].Held[1] = slotState{Item: wood} },
		"the same item twice":  func(s *snapshot) { s.Home[2].Held[1] = slotState{Item: wood, N: 1} },
		"count with no item":   func(s *snapshot) { s.Home[2].Held[1] = slotState{N: 4} },
		"walled in":            func(s *snapshot) { s.Home = wallOf(s.Home) },
		"a table out of doors": func(s *snapshot) { s.Yard[0].Name = "table" },
		"two on one cell":      func(s *snapshot) { s.Yard = append(s.Yard, s.Yard[0]) },
		"a chest far away":     func(s *snapshot) { s.Yard[0].Col = 1 << 40 },
		"too many chests":      func(s *snapshot) { s.Yard = append(s.Yard, make([]pieceState, 600)...) },
		"a name with a null":   func(s *snapshot) { s.Yard[0].Name = "chest\x00" },
		"a bag of thrones":     func(s *snapshot) { s.Player.Bag = []stackState{{Item: "throne", N: 1}} },
	}
	for name, edit := range cases {
		var c snapshot
		json.Unmarshal(b, &c)
		edit(&c)
		_, g := build(longBed{SVGArt()}, idle{})
		before := takeSnapshot(g)
		if err := restore(g, &c); err == nil {
			t.Errorf("%s: accepted", name)
		} else if !reflect.DeepEqual(before, takeSnapshot(g)) {
			t.Errorf("%s: the game changed though the save was refused", name)
		}
	}
}

// wallOf is the bed and a row of chairs across the house, walling it in.
func wallOf(home []pieceState) []pieceState {
	out := []pieceState{home[0]}
	for c := 1; c <= roomCols-2; c++ {
		out = append(out, pieceState{Name: "chair", Col: c, Row: 3, Held: []slotState{}})
	}
	return out
}
