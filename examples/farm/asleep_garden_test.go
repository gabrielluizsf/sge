package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/town"
)

// farTown is a town of f's world far from the player, so asleep.
func farTown(t *testing.T, f *farm) *town.Town {
	for r := 3; r < 40; r++ {
		if tw := f.keeper.Layout.Of(r, r); tw != nil && !f.keeper.Awake(tw) {
			return tw
		}
	}
	t.Fatal("no town far away")
	return nil
}

// A garden of a town the player is far from is tended by its people
// while no one sees them: it grows each morning and never withers.
func TestAnAsleepTownsGardenGrowsAndLives(t *testing.T) {
	_, f := newGame(idle{})
	tw := farTown(t, f)
	before := len(f.plots)
	f.settle(tw)
	if len(f.plots) == before {
		t.Fatal("the town has no garden")
	}
	p := f.plots[before]
	p.crop, p.days, p.stage = 0, 1, stageOf(crops[0], 1)
	for range 2 {
		f.growCrops()
	}
	if p.dead || p.days != 3 {
		t.Errorf("after two mornings asleep: %+v", p)
	}
}
