package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/bond"
	"github.com/gabrielluizsf/sge/life"
)

// A gift from the hotbar makes one fonder as they like it, once a day:
// a second the same day is refused and kept; the next day it is taken.
func TestAGiftIsGivenOnceADay(t *testing.T) {
	_, f := newGame(idle{})
	v := single(t, f)
	love := gifts[favourite(v)]
	if f.gift(v, love) == nil {
		t.Fatal("gave what the player does not have")
	}
	f.bag.Add(love, 2)
	if err := f.gift(v, love); err != nil || v.bond.Affection != loved || f.bag.Count(love) != 1 {
		t.Fatalf("err %v, affection %d, left %d", err, v.bond.Affection, f.bag.Count(love))
	}
	if err := f.gift(v, love); err != bond.ErrToday || f.bag.Count(love) != 1 {
		t.Fatalf("a second gift the same day: %v, left %d", err, f.bag.Count(love))
	}
	f.passDays(1)
	if err := f.gift(v, love); err != nil || v.bond.Affection != 2*loved {
		t.Fatalf("the next day: %v, affection %d", err, v.bond.Affection)
	}
}

// Each villager loves one gift and dislikes another, the smith and the
// lumberjacks loving wood.
func TestLikesAndDislikes(t *testing.T) {
	_, f := newGame(idle{})
	single(t, f)
	for _, v := range f.people {
		got := map[int]int{}
		for _, item := range gifts {
			got[liking(v, item)]++
		}
		if got[loved] != 1 || got[disliked] != 1 || got[liked] != 1 {
			t.Errorf("%s likes %v", v.name, got)
		}
		if (v.smith() || v.role == roleLumberjack) && liking(v, wood) != loved {
			t.Errorf("a %s does not love wood", v.role)
		}
	}
}

// Chatting makes one a little fonder, once a day.
func TestChattingOnceADay(t *testing.T) {
	_, f := newGame(idle{})
	v := single(t, f)
	if f.chat(v) != nil || f.chat(v) != bond.ErrToday || v.bond.Affection != bond.Talk {
		t.Fatalf("affection %d", v.bond.Affection)
	}
}

// Asking out is refused one not fond enough, a child, and kin; and the
// menu offers it only to those it would not be refused.
func TestAskingOutIsRefused(t *testing.T) {
	_, f := newGame(idle{})
	v := single(t, f)
	v.bond.Affection = bond.DateAt - 1
	if err := f.askOut(v); err != bond.ErrNotYet {
		t.Errorf("not fond enough: %v", err)
	}
	if offers(f, v, "Pedir em namoro") {
		t.Error("offered before fond enough")
	}
	v.bond.Affection = bond.DateAt
	if !offers(f, v, "Pedir em namoro") {
		t.Error("not offered to a friend fond enough")
	}
	young := *v
	young.life = calendar.Draw(f.fate, 10, f.clock.Day)
	if err := f.askOut(&young); err != bond.ErrYoung || offers(f, &young, "Pedir em namoro") {
		t.Errorf("a child: %v", err)
	}
	wed(t, f, v)
	kid := haveChild(t, f)
	f.passDays(life.AdultFrom * daysPerYear)
	f.bury(v) // widowed, so free to court
	kid.bond.Affection = bond.Most
	if err := f.askOut(kid); err != bond.ErrKin || offers(f, kid, "Pedir em namoro") {
		t.Errorf("their own child: %v", err)
	}
}

// One dated and fond enough is asked to marry, weds some days on and
// moves into the player's house; while courting, the player asks no one
// else out.
func TestCourtingEndsInAWedding(t *testing.T) {
	_, f := newGame(idle{})
	v := single(t, f)
	v.bond.Affection = bond.ProposeAt
	if f.propose(v) != bond.ErrNotYet {
		t.Error("proposed to one not dated")
	}
	if err := f.askOut(v); err != nil || v.bond.Stage() != bond.Dating {
		t.Fatalf("asked out: %v", err)
	}
	if !offers(f, v, "Pedir em casamento") {
		t.Error("marriage not offered")
	}
	other := &villager{id: "other", life: v.life}
	other.bond.Affection = bond.Most
	if f.askOut(other) != errTaken {
		t.Error("asked another out while dating")
	}
	if err := f.propose(v); err != nil || v.bond.Stage() != bond.Engaged {
		t.Fatalf("proposed: %v", err)
	}
	f.passDays(bond.WeddingAfter - 1)
	if f.self.spouse != nil {
		t.Fatal("wed before the day")
	}
	f.passDays(1)
	if f.self.spouse != v || !v.ours || !livesIn(v, houseDoor{own: true}) {
		t.Fatal("not wed, or not living with the player")
	}
}

// offers reports whether talking to v offers the line text.
func offers(f *farm, v *villager, text string) bool {
	for _, l := range f.lines(v) {
		if l.text == text {
			return true
		}
	}
	return false
}
