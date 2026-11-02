package farm

import "testing"

// What the farmer shows for a use: the axe on a tree is a stroke of the
// axe, on the anvil a smith's stroke of the hammer, which takes longer;
// anything else is the act.
func TestUseClipByWhatIsUsed(t *testing.T) {
	for _, c := range []struct {
		slot    int
		anvil   bool
		clip    string
		seconds float64
	}{
		{slotAxe, false, "_chop", actFor},
		{slotAxe, true, "_forge", ForgeSeconds},
		{slotHoe, true, "_act", actFor},
		{slotHoe, false, "_act", actFor},
	} {
		clip, seconds := useClip(c.slot, c.anvil)
		if clip != c.clip || seconds != c.seconds {
			t.Errorf("slot %d, anvil %v: %s for %v seconds, not %s for %v", c.slot, c.anvil, clip, seconds, c.clip, c.seconds)
		}
	}
}
