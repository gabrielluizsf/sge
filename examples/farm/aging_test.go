package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/lineage"
)

// The towns' people are of every age, children, grown-ups and elders;
// the player starts a young adult; and a year of days on, everyone is a
// year older.
func TestEveryoneAgesWithTheYears(t *testing.T) {
	_, f := newGame(idle{})
	for r := range 3 {
		for c := -3; c < 4; c++ {
			if tn := f.towns.Of(c, r); tn != nil {
				f.settle(tn)
			}
		}
	}
	stages := map[life.Stage]int{}
	years := map[*villager]int{}
	for _, v := range f.people {
		stages[v.stage()]++
		years[v] = v.who.Years
		if v.smith() && v.stage() != life.Adult {
			t.Errorf("a smith of %d", v.who.Years)
		}
	}
	if stages[life.Child] == 0 || stages[life.Adult] == 0 || stages[life.Elder] == 0 {
		t.Errorf("ages: %v", stages)
	}
	if y := f.self.who(f.clock.Day).Years; y != founderYears {
		t.Errorf("the player starts %d years old", y)
	}
	f.passDays(daysPerYear)
	if y := f.self.who(f.clock.Day).Years; y != founderYears+1 {
		t.Errorf("a year on, the player is %d", y)
	}
	for _, v := range f.people {
		if was, ok := years[v]; ok && v.who.Years != was+1 {
			t.Errorf("%s was %d, a year on is %d", v.name, was, v.who.Years)
		}
	}
}

// A child grows up: drawn anew, bigger, as they become a teen.
func TestAChildGrowsUp(t *testing.T) {
	_, f := newGame(idle{})
	v := single(t, f)
	v.life = calendar.Draw(f.fate, life.TeenFrom-0.01, f.clock.Day)
	f.age(v, f.clock.Day)
	f.personOf(v)
	child := height(v.body)
	if v.stage() != life.Child {
		t.Fatalf("a %v", v.stage())
	}
	f.passDays(1)
	if v.stage() != life.Teen || height(v.body) <= child {
		t.Errorf("a day on, a %v, %d pixels tall, was %d", v.stage(), height(v.body), child)
	}
}

// height is how tall b stands, in pixels.
func height(b Body) int { return b.Clip(AnimStand, "down").Frames[0].Make().Height }

// One whose span runs out dies of old age: gone from the towns and the
// game, dead in the family tree, their house standing.
func TestVillagersDieOfOldAge(t *testing.T) {
	_, f := newGame(idle{})
	v := single(t, f)
	houses := len(v.town.Houses)
	v.life.Span = yearsOf(v.life, f.clock.Day) + 0.5/daysPerYear
	f.passDays(1)
	if f.personByID(v.id) != nil {
		t.Fatal("still among the living")
	}
	if len(v.town.Houses) != houses {
		t.Error("their house is gone")
	}
	if _, ok := f.family.Heir(v.id); ok {
		t.Error("an heir for one with no children")
	}
}

// The towns' couples have children, born when due at home, babies in
// the family tree; none once a town is crowded.
func TestCouplesInTownsHaveChildren(t *testing.T) {
	_, f := newGame(idle{})
	single(t, f)
	var mum *villager
	for _, v := range f.people {
		if v.woman && v.partner != nil && f.mayExpect(v, f.clock.Day) {
			mum = v
		}
	}
	if mum == nil {
		t.Fatal("no couple who may have a child")
	}
	e := lineage.Expect([2]string{mum.id, mum.partner.id}, f.clock.Day)
	mum.expecting = &e
	n := f.townSize(mum.town)
	f.passDays(lineage.Gestation)
	if f.townSize(mum.town) < n+1 || len(f.family.Children(mum.id)) == 0 {
		t.Fatalf("no child born: %d in town, was %d", f.townSize(mum.town), n)
	}
	b := f.personByID(f.family.Children(mum.id)[0])
	if b.stage() != life.Baby || b.home != mum.home || b.town != mum.town {
		t.Errorf("the baby: %+v", b.who)
	}
	for range crowd * len(mum.town.Houses) {
		f.people = append(f.people, &villager{town: mum.town})
	}
	if f.mayExpect(mum, f.clock.Day) {
		t.Error("a crowded town has more children")
	}
}
