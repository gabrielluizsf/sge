package farm

import "testing"

// Filling the can at a town's well is shown by the well's motion, and at
// the water of the farm, a river or a lake by the crouch at the bank;
// the farmer stays where they are for as long as it lasts.
func TestFillingIsShownByWhereItIsDone(t *testing.T) {
	g, f := newGame(pressing("up"))
	o := g.Find("farmer")
	f.kit.slot = slotCan
	col, row := shallowCell(t, f)
	f.kit.used = 5
	if got := f.useAs(f.kit, col, row); got != actFillRiver {
		t.Errorf("filling at the pond is %d, not the river's", got)
	}
	w := wellRect(f.art, eastTown(f))
	f.kit.used = 5
	if got := f.useAs(f.kit, w.col, w.row); got != actFillWell {
		t.Errorf("filling at the well is %d, not the well's", got)
	}
	o.Transform.Position = f.ground.CellCenter(col, row)
	f.kit.used = 5
	g.Step(frame)
	if f.hands.acting < DipSeconds/2 {
		t.Errorf("the farmer is held %v s at the bank, not the length of the dip", f.hands.acting)
	}
}

// Eating is held for as long as the eating lasts, so it is seen whole.
func TestEatingHoldsThePlayerForItsLength(t *testing.T) {
	g, f := newGame(pressing("up"))
	f.bag.Add(cropOf(2), 2)
	f.pack.Hotbar[4] = cropOf(2)
	f.hold(4)
	f.vit.Food = 30
	g.Step(frame)
	if f.hands.acting < EatSeconds/2 {
		t.Errorf("the farmer is held %v s eating", f.hands.acting)
	}
}
