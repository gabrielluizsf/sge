package farm

import "testing"

// What the farmer shows for a use is what the use did: a stroke of the
// axe on a tree, of the hammer on the anvil, of the hoe, a pour, a pick;
// filling the can and sowing are the act.
func TestUseClipByWhatTheUseDid(t *testing.T) {
	for _, c := range []struct {
		did     action
		clip    string
		seconds float64
	}{
		{actChop, "_chop", actFor},
		{actForge, "_forge", ForgeSeconds},
		{actHoe, "_hoe", HoeSeconds},
		{actWater, "_water", WaterSeconds},
		{actPick, "_pick", PickSeconds},
		{actFill, "_act", actFor},
		{actSow, "_act", actFor},
	} {
		clip, seconds := useClip(c.did)
		if clip != c.clip || seconds != c.seconds {
			t.Errorf("action %d: %s for %v seconds, not %s for %v", c.did, clip, seconds, c.clip, c.seconds)
		}
	}
}

// Picking a ripe crop is a pick whatever is in hand, the can that waters
// included; watering a thirsty plot is a pour; the hoe clearing a withered
// plant is a stroke of it; sowing is sowing.
func TestUseReportsWhatItDid(t *testing.T) {
	_, f := newGame(idle{})
	p := sown(f, 0)
	f.slot = slotSeeds
	if f.useAs(f.kit, p.col, p.row) != actNone {
		t.Error("sowed in a plot that is sown")
	}
	f.slot = slotCan
	if got := f.useAs(f.kit, p.col, p.row); got != actWater {
		t.Errorf("watering a thirsty plot is %d", got)
	}
	p.days, p.stage = crops[0].daysToRipe, ripe
	if got := f.useAs(f.kit, p.col, p.row); got != actPick {
		t.Errorf("a ripe crop picked with the can in hand is %d", got)
	}
	p = sown(f, 0)
	p.dead = true
	f.slot = slotHoe
	if got := f.useAs(f.kit, p.col, p.row); got != actHoe {
		t.Errorf("clearing a withered plant with the hoe is %d", got)
	}
	f.slot, f.seed = slotSeeds, 0
	f.bag.Add(seedsOf(0), 1)
	if got := f.useAs(f.kit, p.col, p.row); got != actSow {
		t.Errorf("sowing is %d", got)
	}
}
