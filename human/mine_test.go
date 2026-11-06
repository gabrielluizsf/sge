package human

import "testing"

// Mining is the felling cut's swing under a name of its own, so the
// pickaxe shows in it and the axe does not; it has its own grip.
func TestMiningSwingsAsTheChopDoes(t *testing.T) {
	if Mine.Name == Chop.Name || len(Mine.Poses) != len(Chop.Poses) || Mine.Length != Chop.Length {
		t.Fatalf("mine %q: %d poses of %v s", Mine.Name, len(Mine.Poses), Mine.Length)
	}
	for k := range Mine.Poses {
		if Mine.Poses[k] != Chop.Poses[k] {
			t.Fatalf("pose %d differs from the chop's", k)
		}
	}
	if len(MineGrip.During) != 1 || MineGrip.During[0] != Mine.Name {
		t.Errorf("the pickaxe is held in %v", MineGrip.During)
	}
}
