package farm

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/town"
)

// chorePlace finds a grown villager of town tn and a day on which they have
// a spell of chore c with a place to be at, and the middle hour of it.
func chorePlace(f *farm, tn *town.Town, c chore) (v *villager, day int, hour float64) {
	for day := range 40 {
		for _, v := range f.people {
			if v.town != tn || v.ours || v.stage() < life.Adult {
				continue
			}
			for i, s := range f.planOn(v, day) {
				if s.doing == c && f.siteOf(v, c, spellKey(v, day, i)) != nil {
					return v, day, (s.from + s.to) / 2
				}
			}
		}
	}
	return nil, 0, 0
}

// sameSite reports whether two jobs are at the same place, doing the same.
func sameSite(a, b *job) bool {
	return a != nil && b != nil && a.to == b.to && a.face == b.face && a.slot == b.slot && a.uses == b.uses && a.target == b.target
}

// What a villager is doing at an hour, and where, is the same whichever way
// the hour is reached: asked from the morning or from the night, at their
// town's waking after a long time away, or just after the last.
func TestAnHoursPlaceIsTheSameHoweverItIsReached(t *testing.T) {
	for _, c := range []chore{choreChop, chorePlant, choreTrade} {
		_, f := newGame(idle{})
		tn := farTown(t, f)
		f.settle(tn)
		v, day, hour := chorePlace(f, tn, c)
		if v == nil {
			t.Fatalf("chore %d: nobody has it", c)
		}
		f.clock.Day, f.clock.Hour = day, 6
		morning := f.presenceAt(v, day, hour)
		f.clock.Day, f.clock.Hour = day+9, 23
		night := f.presenceAt(v, day, hour)
		if morning.job == nil || morning.chore != c || night.chore != c || !sameSite(morning.job, night.job) {
			t.Fatalf("chore %d: from the morning %+v, from the night %+v", c, morning, night)
		}
		// The town wakes at that hour after days away: they stand there, the
		// job in hand.
		f.clock.Day, f.clock.Hour = day, hour
		f.leftAt(tn)
		f.wake(tn)
		if want := f.center(morning.job.to); v.pos != want || v.mind.job == nil || !sameSite(v.mind.job, morning.job) {
			t.Errorf("chore %d: woke at %v with %+v, the place is %v", c, v.pos, v.mind.job, want)
		}
	}
}

// A villager away on their trip when the town wakes is on the road, not at a
// chore.
func TestAnAwayVillagerHasNoChoreAtTheirHour(t *testing.T) {
	_, f := newGame(idle{})
	tn := farTown(t, f)
	f.settle(tn)
	for day := range 40 {
		for _, v := range f.people {
			if v.town != tn || v.ours || v.stage() < life.Adult {
				continue
			}
			if s, ok := f.planOn(v, day).visit(); ok {
				if p := f.presenceAt(v, day, (s.from+s.to)/2); p.job != nil {
					t.Fatalf("%s is on a visit at %v and has the chore %+v", v.id, (s.from+s.to)/2, p)
				}
				return
			}
		}
	}
	t.Skip("no one visits in forty days")
}

// The sales between villagers at the market make no money and lose none:
// what the buyer pays the seller has.
func TestSalesBetweenVillagersKeepTheMoney(t *testing.T) {
	mk := func(money, logs int) *villager {
		v := &villager{kit: &kit{money: money}}
		v.kit.bag.Add(wood, logs)
		return v
	}
	a, b, c := mk(5, 9), mk(40, 0), mk(7, 1)
	total := func() int { return a.kit.money + b.kit.money + c.kit.money }
	before := total()
	f := &farm{}
	f.settleSwaps([]trading{{a, 9, 12}, {b, 9, 12}, {c, 10, 11}})
	if total() != before {
		t.Errorf("money %d before, %d after", before, total())
	}
	if got := a.kit.bag.Count(wood) + b.kit.bag.Count(wood) + c.kit.bag.Count(wood); got != 10 {
		t.Errorf("%d wood, there were 10", got)
	}
	if b.kit.bag.Count(wood) == 0 || b.kit.money >= 40 {
		t.Errorf("the rich one bought nothing: %d wood, $%d", b.kit.bag.Count(wood), b.kit.money)
	}
	if n := b.kit.bag.Count(wood); n > swapMost {
		t.Errorf("%d bought in one stretch, the most is %d", n, swapMost)
	}
}

// A town left alone for a time is settled from its people's plans, felled
// trees and all, and a time away longer than the most is settled as the
// most is, so the work of waking is the same for ten days as for a hundred.
func TestAwayTimeIsSettledAndBounded(t *testing.T) {
	run := func(away float64) (stumps, logs int) {
		_, f := newGame(idle{})
		tn := farTown(t, f)
		f.settle(tn)
		f.clock.Day, f.clock.Hour = 500, 18
		f.away = map[*town.Town]float64{tn: f.hours() - away}
		f.catchUp(tn)
		for _, v := range f.people {
			if v.town == tn {
				logs += v.kit.bag.Count(wood)
			}
		}
		return len(f.stumps), logs
	}
	s, _ := run(maxAway)
	if s == 0 {
		t.Fatal("three days away, and no tree was felled")
	}
	long, _ := run(100 * 24)
	if long != s {
		t.Errorf("%d stumps after the most, %d after a hundred days", s, long)
	}
	if none, _ := run(minAway / 2); none != 0 {
		t.Errorf("%d trees felled in a quarter of an hour", none)
	}
	if bound := 12 * 4 * 6 * (fellMost + mineMost); s > bound {
		t.Errorf("%d stumps, more than %d a day could make", s, bound)
	}
}

// A time away is saved with a strict range, and an old save without it loads.
func TestTimeAwayIsSavedAndOldSavesLoad(t *testing.T) {
	f := played(t)
	var tn *town.Town
	for _, p := range f.plots {
		if p.town != nil {
			tn = p.town
			break
		}
	}
	if tn == nil {
		t.Fatal("no town woken")
	}
	f.leftAt(tn)
	s := takeSnapshot(f)
	if len(s.Away) != 1 || math.Abs(s.Away[0].At-f.hours()) > 1e-9 {
		t.Fatalf("saved %+v", s.Away)
	}
	b, _ := json.Marshal(s)
	var back snapshot
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatal(err)
	}
	_, g := newGame(idle{})
	if err := restore(g, &back); err != nil || len(g.away) != 1 {
		t.Fatalf("restore: %v, %d away", err, len(g.away))
	}
	for name, bad := range map[string]func(*snapshot){
		"future":    func(s *snapshot) { s.Away[0].At = f.hours() + 5 },
		"negative":  func(s *snapshot) { s.Away[0].At = -1 },
		"not a num": func(s *snapshot) { s.Away[0].At = math.NaN() },
		"no town":   func(s *snapshot) { s.Away[0].Col = 9999 },
		"twice":     func(s *snapshot) { s.Away = append(s.Away, s.Away[0]) },
	} {
		var hostile snapshot
		json.Unmarshal(b, &hostile)
		bad(&hostile)
		if _, err := check(g, &hostile); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	var old snapshot
	json.Unmarshal(b, &old)
	old.Away = nil
	_, g2 := newGame(idle{})
	if err := restore(g2, &old); err != nil || len(g2.away) != 0 {
		t.Errorf("an old save: %v, %d away", err, len(g2.away))
	}
}
