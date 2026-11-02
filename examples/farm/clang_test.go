package farm

import "testing"

// The anvil rings as the hammer lands on the iron, not as it is raised:
// a stroke's clang comes forgeImpact of the way into it.
func TestTheAnvilRingsAsTheHammerLands(t *testing.T) {
	_, f := newGame(idle{})
	e := &ears{}
	f.speaker = e
	tn := eastTown(f)
	f.kit.slot = slotAxe
	if !f.use(f.kit, tn.Anvil.Col, tn.Anvil.Row) {
		t.Fatal("the axe did nothing at the anvil")
	}
	if e.count("clang") != 0 {
		t.Fatal("the anvil rang as the hammer was raised")
	}
	for range int(forgeImpact*ForgeSeconds/frame) + 2 {
		f.ring(frame) // the cues alone: the town's smith strikes too
	}
	if e.count("clang") != 1 {
		t.Errorf("the anvil rang %d times at the blow, not once", e.count("clang"))
	}
}
