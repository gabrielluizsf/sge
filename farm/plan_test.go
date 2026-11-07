package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/town"
)

// A person's plan differs from one day to the next and from one person to
// the next, and the same call always gives the same plan.
func TestAPlanDiffersByPersonAndDay(t *testing.T) {
	call := func(id string, day int) dayPlan {
		return planFor(planCall{id: id, role: roleWanderer, day: day, trip: 0.5})
	}
	same := func(a, b dayPlan) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	if !same(call("t1-0", 3), call("t1-0", 3)) {
		t.Fatal("the same call drew two plans")
	}
	days, people := 0, 0
	for d := range 20 {
		if !same(call("t1-0", d), call("t1-0", d+1)) {
			days++
		}
		if !same(call("t1-0", d), call("t1-1", d)) {
			people++
		}
	}
	if days < 15 || people < 15 {
		t.Errorf("of 20 draws, %d differ by day and %d by person", days, people)
	}
}

// A day has several stretches of work, each trade turns most to its own
// chore, and a trip to the next town is home well before bed.
func TestAPlanHasSeveralChoresAndTripsEndBeforeBed(t *testing.T) {
	counts := map[string][choreCount]int{}
	trips := 0
	for _, role := range []string{roleFarmer, roleLumberjack, roleSmith, roleWanderer} {
		for d := range 200 {
			p := planFor(planCall{id: "x" + role, role: role, day: d, trip: 0.4})
			if len(p) < 2 {
				t.Fatalf("%s day %d: %d spells", role, d, len(p))
			}
			c := counts[role]
			for i, s := range p {
				c[s.doing]++
				if s.to <= s.from || i > 0 && s.from < p[i-1].to {
					t.Fatalf("%s day %d: spells out of order: %v", role, d, p)
				}
				if s.doing == choreVisit {
					trips++
					if s.to > bedBy-1+1e-9 || s.to-s.from < 2*0.4+minStay {
						t.Fatalf("%s day %d: a visit from %.2f to %.2f", role, d, s.from, s.to)
					}
				}
			}
			counts[role] = c
		}
	}
	if trips == 0 {
		t.Fatal("no one ever travels")
	}
	for role, want := range map[string]chore{roleFarmer: chorePlant, roleLumberjack: choreChop, roleSmith: choreForge} {
		best := chore(0)
		for c := chore(0); c < choreVisit; c++ {
			if counts[role][c] > counts[role][best] {
				best = c
			}
		}
		if best != want {
			t.Errorf("%s turns most to chore %d, not %d: %v", role, best, want, counts[role])
		}
	}
	if p := planFor(planCall{id: "a", role: roleFarmer, day: 1}); len(p) == 0 {
		t.Error("no plan without a road")
	} else if _, ok := p.visit(); ok {
		t.Error("a visit with no road")
	}
}

// townWithRoad is a town of the farm's with a road to another.
func townWithRoad(t *testing.T, f *farm) *town.Town {
	t.Helper()
	for r := range 6 {
		for c := -3; c < 6; c++ {
			if tn := f.towns.Of(c, r); tn != nil && f.roadFrom(tn) != nil {
				return tn
			}
		}
	}
	t.Fatal("no town has a road")
	return nil
}

// A traveller is in one place at a time, by the clock: home, then on the
// road out, then there, then back on the road, then home, in that order,
// the road never crossed twice, and home well before bed.
func TestATravellerIsInOnePlaceAndHomeBeforeBed(t *testing.T) {
	_, f := newGame(idle{})
	tn := townWithRoad(t, f)
	f.settle(tn)
	var v *villager
	day := 0
	for d := range 60 {
		for _, p := range f.people {
			if p.town != tn || p.stage() < life.Adult {
				continue
			}
			f.clock.Day = d
			if _, ok := f.planOf(p).visit(); ok {
				v, day = p, d
			}
		}
		if v != nil {
			break
		}
	}
	if v == nil {
		t.Fatal("no one travels in 60 days")
	}
	f.clock.Day = day
	v.planDay = 0
	order := []stay{atHome, outbound, there, inbound, atHome}
	seen, k := []stay{atHome}, 0
	for h := 0.0; h < 24; h += 0.01 {
		f.clock.Hour = h
		w := f.whereabouts(v)
		if w.at < 0 || w.at > 1 {
			t.Fatalf("hour %.2f: %v along the road", h, w.at)
		}
		if w.stay != seen[len(seen)-1] {
			seen = append(seen, w.stay)
			if k++; k >= len(order) || seen[len(seen)-1] != order[k] {
				t.Fatalf("hour %.2f: stays %v out of order", h, seen)
			}
		}
		if h >= bedBy && w.stay != atHome {
			t.Fatalf("hour %.2f: not home by bed", h)
		}
	}
	if len(seen) != len(order) {
		t.Errorf("stays %v", seen)
	}
}

// Trades between villagers conserve money and goods, and need a buyer who
// can pay.
func TestATradeConservesMoneyAndGoods(t *testing.T) {
	a := &villager{kit: &kit{money: 10}}
	b := &villager{kit: &kit{money: 3}}
	a.kit.bag.Add(wood, 2)
	if _, ok := saleBetween(a, b); ok {
		t.Fatal("a buyer with $3 bought wood at $4")
	}
	b.kit.money = 20
	for range 2 {
		s, ok := saleBetween(a, b)
		if !ok {
			t.Fatal("no trade")
		}
		s.settle()
	}
	if a.kit.money+b.kit.money != 30 || a.kit.money != 18 || a.kit.bag.Count(wood) != 0 || b.kit.bag.Count(wood) != 2 {
		t.Errorf("a $%d %d wood, b $%d %d wood", a.kit.money, a.kit.bag.Count(wood), b.kit.money, b.kit.bag.Count(wood))
	}
}
