package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/town"
)

// kin is the people of house h of town t.
func kin(f *farm, t *town.Town, h *town.House) []*villager {
	var out []*villager
	for _, v := range f.people {
		if v.town == t && v.home == h {
			out = append(out, v)
		}
	}
	return out
}

// Every house of a new town has a household: a man and a woman who are
// each other's partners, and a child that is theirs.
func TestEveryNewHouseHasAMenWomanAndChild(t *testing.T) {
	_, f := newGame(idle{})
	for _, rc := range [][2]int{{1, 0}, {2, 0}, {0, 1}, {1, 1}} {
		tn := f.towns.Of(rc[0], rc[1])
		if tn == nil {
			continue
		}
		f.settle(tn)
		for i := range tn.Houses {
			h := &tn.Houses[i]
			var men, women, young int
			for _, v := range kin(f, tn, h) {
				switch {
				case v.partner == nil:
					young++
				case v.woman:
					women++
				default:
					men++
				}
				if v.partner != nil && (v.partner.home != h || v.partner.woman == v.woman) {
					t.Errorf("house %d: %s has no partner of the other sex at home", i, v.id)
				}
			}
			if men != 1 || women != 1 || young != 1 {
				t.Errorf("town %v house %d: %d men, %d women, %d young", rc, i, men, women, young)
			}
		}
	}
}

// The child's parents in the family tree are the couple, and the couple is
// old enough to be them.
func TestTheChildIsTheCouplesInTheFamilyTree(t *testing.T) {
	_, f := newGame(idle{})
	tn := eastTown(f)
	h := &tn.Houses[0]
	var man, woman, child *villager
	for _, v := range kin(f, tn, h) {
		switch {
		case v.partner == nil:
			child = v
		case v.woman:
			woman = v
		default:
			man = v
		}
	}
	if man == nil || woman == nil || child == nil {
		t.Fatal("household incomplete")
	}
	for _, parent := range []*villager{man, woman} {
		if kids := f.family.Children(parent.id); len(kids) != 1 || kids[0] != child.id {
			t.Errorf("%s has children %v, want %s", parent.id, kids, child.id)
		}
	}
}

// A full household has no more children, so a town never grows past
// three to a house.
func TestAFullHouseholdHasNoMoreChildren(t *testing.T) {
	_, f := newGame(idle{})
	eastTown(f)
	for _, v := range f.people {
		if f.mayExpect(v, f.clock.Day) {
			t.Fatalf("%s may expect a child in a full house", v.id)
		}
	}
}

// A save made before households, with one person to a house, loads as it
// was: no one is added, killed or married off to make the rule hold, and
// waking the town again does not people it anew.
func TestAnOldSaveWithOnePersonToAHouseLoadsUnchanged(t *testing.T) {
	_, f := newGame(idle{})
	tn := eastTown(f)
	s := takeSnapshot(f)
	var old []villagerState
	for i, p := range s.People {
		if i%3 == 0 {
			p.Partner = -1
			old = append(old, p)
		}
	}
	s.People = old
	for i := range s.Plots {
		s.Plots[i].Owner = 0
	}
	_, g := newGame(idle{})
	if err := restore(g, &s); err != nil {
		t.Fatal(err)
	}
	if len(g.people) != len(old) {
		t.Fatalf("%d people, the save has %d", len(g.people), len(old))
	}
	g.settle(g.towns.Of(1, 0))
	if len(g.people) != len(old) || len(old) != len(tn.Houses) {
		t.Errorf("waking the town gave %d people for %d houses", len(g.people), len(tn.Houses))
	}
}

// Every town house has the couple's double bed and the child's single one,
// whatever the size of a bed in the look: a cell, or two across for the
// single and two by two for the double.
func TestEveryRoomHasABedForEachOfTheHousehold(t *testing.T) {
	for _, size := range []func(byte) (int, int){
		func(byte) (int, int) { return 1, 1 },
		func(p byte) (int, int) {
			switch p {
			case 'D':
				return 2, 2
			case 'b', 'B':
				return 2, 1
			}
			return 1, 1
		},
	} {
		short := 0
		for seed := range uint32(400) {
			_, things := roomLayout(seed, size)
			beds := bedsOf(bedsIn(things, size))
			if len(beds) != 2 || !beds[0].double || beds[1].double {
				short++
			}
		}
		if short > 0 {
			t.Errorf("%d of 400 rooms lack the double bed or the single one", short)
		}
	}
}
