package farm

import "testing"

// A morning's play with the Bot on the farm: the farmers of the towns
// about come and sow the player's soil, and keep it.
func TestVillagersPlantOnTheFarmThroughTheDay(t *testing.T) {
	g, f := newGame(Bot(2))
	for c := -2; c <= 2; c++ {
		for r := -2; r <= 2; r++ {
			if tn := f.towns.Of(c, r); tn != nil {
				f.settle(tn)
			}
		}
	}
	farmers := 0
	for _, v := range f.people {
		if v.role == roleFarmer && !v.ours {
			v.kit.bag.Add(seedsOf(cornCrop), 6)
			v.kit.seed = cornCrop
			farmers++
		}
	}
	if farmers == 0 {
		t.Skip("no farmers about")
	}
	f.clock.Hour = 6
	for range 60 * 240 {
		g.Step(frame)
	}
	if f.held() == 0 {
		t.Fatal("no villager planted on the farm")
	}
	for _, p := range f.farmPlots() {
		if p.owner != nil && p.crop < 0 {
			t.Errorf("a kept plot with nothing in it: %+v", p)
		}
	}
}
