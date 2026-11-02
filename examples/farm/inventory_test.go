package farm

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/backpack"
)

// The seeds on slot 1 are the one packet chosen in the inventory: pressing
// 1 again does not cycle, and running out does not sow another crop.
func TestSlotOneHoldsTheSeedChosenInTheInventory(t *testing.T) {
	_, f := newGame(idle{})
	for _, c := range []int{0, 3, 5} {
		f.bag.Add(seedsOf(c), 2)
	}
	if ev := f.pack.Bind(seedsOf(3), 0); ev != backpack.Bound {
		t.Fatalf("bind: %v", ev)
	}
	f.keepHand()
	for range 3 {
		f.hold(0) // 1, again and again
	}
	if f.slot != slotSeeds || f.seed != 3 {
		t.Fatalf("slot %d, seed %d", f.slot, f.seed)
	}
	f.bag.Take(seedsOf(3), 2)
	if _, ok := f.sowing(); ok {
		t.Error("sowed other seeds when the chosen ran out")
	}
}

// Moving a tool to another slot keeps it in hand by what it is, and the
// axe a player gets goes to the first free slot.
func TestToolsMoveBetweenSlotsAndTheAxeIsStocked(t *testing.T) {
	_, f := newGame(idle{})
	f.pack.Bind(itemCan, 0)
	f.keepHand()
	if f.slot != slotCan || f.pack.Hotbar[1] != seedsOf(0) {
		t.Fatalf("slot %d, hotbar %v", f.slot, f.pack.Hotbar)
	}
	if f.pack.Slot(itemAxe) >= 0 {
		t.Fatal("an axe before there is one")
	}
	f.axe = 0
	f.stock(itemAxe)
	if f.pack.Slot(itemAxe) != 2 {
		t.Errorf("hotbar %v", f.pack.Hotbar)
	}
}

// The inventory lists everything carried, with a tooltip that says what
// each is for, and a keyboard path: pick the seeds, go down to the bar and
// put them on slot 5.
func TestTheInventoryListsWhatIsCarriedAndTakesKeys(t *testing.T) {
	_, f := newGame(idle{})
	f.bag.Add(seedsOf(4), 3)
	f.bag.Add(cropOf(1), 2)
	f.bag.Add(wood, 1)
	items := carried{f}.Items()
	var names []string
	for _, it := range items {
		names = append(names, it.Name)
		if it.Name == "" || len(it.Info) == 0 {
			t.Errorf("%q has no tooltip", it.ID)
		}
	}
	if len(items) != 2+3 || !strings.Contains(strings.Join(names, ","), "Sementes de feijão") {
		t.Fatalf("items %v", names)
	}
	f.pack.Show()
	f.pack.Handle(backpack.Input{Dx: 2}) // the seeds, after the can and the hoe
	f.pack.Handle(backpack.Input{Number: 5})
	if f.pack.Hotbar[4] != seedsOf(4) {
		t.Errorf("hotbar %v", f.pack.Hotbar)
	}
}
