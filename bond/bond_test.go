package bond

import (
	"errors"
	"testing"
)

var grown = Moment{Grown: true}

// Gifts and talk, a day at a time, make strangers friends; a second gift
// or chat on one day counts for nothing; a disliked gift sets back.
func TestGiftsAndTalkMakeFriends(t *testing.T) {
	var b Bond
	if b.Stage() != Strangers {
		t.Fatal(b.Stage())
	}
	if err := b.Give(90, 0); err != nil {
		t.Fatal(err)
	}
	if err := b.Give(90, 0); !errors.Is(err, ErrToday) || b.Affection != 90 {
		t.Errorf("a second gift today: %v, affection %d", err, b.Affection)
	}
	if b.Chat(0) != nil || b.Chat(0) == nil || b.Stage() != Acquainted {
		t.Errorf("chatting: affection %d, %v", b.Affection, b.Stage())
	}
	for day := 1; b.Stage() < Friends; day++ {
		b.Give(80, day)
	}
	b.Give(-500, 1000)
	if b.Affection != 0 {
		t.Errorf("a hated gift leaves %d", b.Affection)
	}
	b.Give(5000, 1001)
	if b.Affection != Most {
		t.Errorf("affection over the most: %d", b.Affection)
	}
}

// Only a friend fond enough is asked out; only one dated and fond enough
// asked to marry; they wed days on, and only then may have a child.
func TestCourtshipThenMarriageThenChild(t *testing.T) {
	b := Bond{Affection: DateAt - 1}
	if err := b.Propose(grown); !errors.Is(err, ErrNotYet) {
		t.Errorf("proposing before dating: %v", err)
	}
	if err := b.AskOut(grown); !errors.Is(err, ErrNotYet) {
		t.Errorf("asking out not fond enough: %v", err)
	}
	b.Affection = DateAt
	if err := b.AskOut(grown); err != nil || b.Stage() != Dating {
		t.Fatalf("asking out: %v %v", err, b.Stage())
	}
	if err := b.AskOut(grown); !errors.Is(err, ErrAlready) {
		t.Errorf("asking out again: %v", err)
	}
	if err := b.Propose(grown); !errors.Is(err, ErrNotYet) {
		t.Errorf("proposing not fond enough: %v", err)
	}
	b.Affection = ProposeAt
	if err := b.Propose(Moment{Day: 10, Grown: true}); err != nil || b.Stage() != Engaged {
		t.Fatalf("proposing: %v %v", err, b.Stage())
	}
	if b.MayHaveChild() || b.Wed(12) {
		t.Error("a child, or wed, before the wedding")
	}
	if !b.Wed(13) || !b.MayHaveChild() || b.Stage() != Married {
		t.Errorf("on the wedding day: %v", b.Stage())
	}
	if err := b.Propose(grown); !errors.Is(err, ErrAlready) {
		t.Errorf("proposing to one's spouse: %v", err)
	}
}

// No one courts a child, nor their kin, however fond.
func TestNoCourtingChildrenOrKin(t *testing.T) {
	for _, m := range []Moment{{Grown: false}, {Grown: true, Kin: true}} {
		b := Bond{Affection: Most}
		if err := b.AskOut(m); err == nil || b.Stage() == Dating {
			t.Errorf("%+v: asked out", m)
		}
		b.Courted = Dating
		if err := b.Propose(m); err == nil || b.Stage() == Engaged {
			t.Errorf("%+v: proposed", m)
		}
	}
}
