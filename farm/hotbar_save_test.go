package farm

import (
	"encoding/json"
	"testing"
)

// handItems are all the items the inventory lets the player take in hand.
func handItems(f *farm) []string {
	ids := []string{itemCan, itemHoe, itemAxe, itemPick}
	for i := range crops {
		ids = append(ids, seedsOf(i), cropOf(i))
	}
	for _, g := range goods {
		ids = append(ids, g.item)
	}
	ids = append(ids, furnishings...)
	var out []string
	for _, id := range ids {
		if it, ok := (carried{f}).Find(id); ok && it.Hand {
			out = append(out, id)
		}
	}
	return out
}

// Whatever the player may put on the hotbar is kept in a save and loads:
// a save refused for what the game itself let them do is a game lost.
func TestEveryItemOnTheHotbarSavesAndLoads(t *testing.T) {
	_, f := newGame(idle{})
	f.axe = 0
	f.bag.Add(itemPick, 1)
	for _, id := range handItems(f) {
		if id != itemCan && id != itemHoe && id != itemAxe && id != itemPick {
			f.bag.Add(id, 1)
		}
		if ev := f.pack.Bind(id, 4); ev != 0 && f.pack.Hotbar[4] != id {
			t.Fatalf("%s did not go on the hotbar", id)
		}
		f.hold(4)
		data, err := json.Marshal(takeSnapshot(f))
		if err != nil {
			t.Fatal(err)
		}
		var s snapshot
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		_, g := newGame(idle{})
		if err := restore(g, &s); err != nil {
			t.Errorf("%s on the hotbar: %v", id, err)
		}
	}
}

// What no game could have put on the hotbar is still refused.
func TestAHotbarOfNonsenseIsRefused(t *testing.T) {
	_, f := newGame(idle{})
	for _, id := range []string{"dragon", wood, stone, wool, " milk"} {
		s := takeSnapshot(f)
		s.Hotbar[4] = id
		_, g := newGame(idle{})
		if err := restore(g, &s); err == nil {
			t.Errorf("%q was let on the hotbar", id)
		}
	}
}
