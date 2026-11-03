package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

func down() vector.Vector3 { return vector.Vector3{Y: 1} }

// Holding Shift is faster than walking, by about four fifths, on land;
// in shallow water they wade and in deep they swim, however they mean to
// go, so there is no running there.
func TestRunIsFasterOnLandAndNowhereElse(t *testing.T) {
	if r := speedOf(true, false, false) / speedOf(false, false, false); r < 1.6 || r > 2 {
		t.Errorf("running is %.2f times walking", r)
	}
	if speedOf(true, true, false) != wadeSpeed || speedOf(false, true, false) != wadeSpeed {
		t.Error("they run in the shallows")
	}
	if speedOf(true, true, true) != swimSpeed || speedOf(true, false, true) != swimSpeed {
		t.Error("they run swimming")
	}
}

// What is shown while moving is the run clip when running, the walk when
// not, and the swim in deep water.
func TestRunPlaysTheRunClip(t *testing.T) {
	for _, c := range []struct {
		running, swims bool
		clip           string
	}{{true, false, "_run"}, {false, false, "_walk"}, {true, true, "_swim"}} {
		if got := moveClip(c.running, c.swims); got != c.clip {
			t.Errorf("running %v, swimming %v: %s, not %s", c.running, c.swims, got, c.clip)
		}
	}
}

// A person runs only if told to and moving, on land: a villager's intent
// never says so, and a keyboard's Shift does.
func TestOnlyAnIntentToRunRuns(t *testing.T) {
	if !running(intent{dir: down(), run: true}, false, false) {
		t.Error("did not run on land with Shift")
	}
	if running(intent{dir: down()}, false, false) || running(intent{run: true}, false, false) {
		t.Error("ran without Shift, or standing")
	}
	if running(intent{dir: down(), run: true}, true, false) || running(intent{dir: down(), run: true}, false, true) {
		t.Error("ran in the water")
	}
}
