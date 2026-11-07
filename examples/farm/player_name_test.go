package farm

import (
	"strings"
	"testing"
)

type namedArt struct {
	Art
	name string
}

func (a namedArt) PlayerName() string { return a.name }

func TestTheFarmerHasTheNameTheArtGives(t *testing.T) {
	_, f := build(namedArt{SVGArt(), "Thandi"}, idle{})
	if f.self.name != "Thandi" {
		t.Errorf("the farmer is %q", f.self.name)
	}
}

func TestABadOrMissingNameLeavesTheFoundersName(t *testing.T) {
	for _, n := range []string{"", " x", "A\nB", "<i>", strings.Repeat("a", 21), "\xff"} {
		if got := playerName(namedArt{SVGArt(), n}); got != founderName {
			t.Errorf("%q gave %q", n, got)
		}
	}
	if playerName(SVGArt()) != founderName {
		t.Error("an art that names nobody changed the name")
	}
}

func TestRandomNamesFitAndFollowTheSex(t *testing.T) {
	for seed := range uint32(40) {
		for _, w := range []bool{true, false} {
			n := RandomName(w, seed)
			if playerName(namedArt{SVGArt(), n}) != n {
				t.Errorf("%q is not a name the farm takes", n)
			}
			if n != RandomName(w, seed) {
				t.Error("not repeatable")
			}
		}
	}
}

// Starting again after a line ends keeps the name the player chose: the new
// farmer is not given the founder's.
func TestRestartingKeepsThePlayersName(t *testing.T) {
	_, f := build(namedArt{SVGArt(), "Gabriel"}, idle{})
	if f.self.name != "Gabriel" {
		t.Fatalf("named %s", f.self.name)
	}
	f.restart()
	if f.self.name != "Gabriel" {
		t.Errorf("after restarting, named %s", f.self.name)
	}
}
