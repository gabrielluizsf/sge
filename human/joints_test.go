package human

import (
	"slices"
	"testing"
)

// No pose of any motion takes a joint past what a healthy adult's can do:
// the arm of the crawl does not turn full circle, nor a knee bend back.
func TestMotionsStayWithinHealthyJoints(t *testing.T) {
	for _, m := range slices.Concat(Motions, Crafts, Farming, Mining, Sitting, Eating, Drawing, Sleeping) {
		for k, p := range m.Poses {
			for _, b := range p.Beyond() {
				t.Errorf("%s, pose %d: %s", m.Name, k, b)
			}
		}
	}
}

// The check catches what it should: a spinning arm, a knee bent back.
func TestBeyondCatchesTheImpossible(t *testing.T) {
	if len(Pose{Raise: both(6.0)}.Beyond()) != 2 {
		t.Error("an arm turned full circle is healthy")
	}
	if len(Pose{Knee: [2]float64{-0.5, 0}}.Beyond()) != 1 {
		t.Error("a knee bent back is healthy")
	}
	if len(Pose{}.Beyond()) != 0 {
		t.Error("standing is not healthy")
	}
}

// A crawl's arms reach no further round than a shoulder can go, and do
// not come back the way they went: each lifts from hip to overhead and
// back, across less than a half turn.
func TestCrawlArmDoesNotSpin(t *testing.T) {
	lo, hi := 10.0, -10.0
	for _, p := range Swim.Poses {
		lo, hi = min(lo, p.Raise[0]), max(hi, p.Raise[0])
	}
	if hi-lo > 2.7 || hi > 3.0 {
		t.Errorf("an arm sweeps %v to %v", lo, hi)
	}
}
