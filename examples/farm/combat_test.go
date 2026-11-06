package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

// fighting is a session at a fight: the smith struck, going after the
// player, whose health they have hurt by the end.
func fighting(t *testing.T) (*Session, *villager) {
	t.Helper()
	s := New(SVGArt(), Bot(10))
	if _, err := s.Open("fight", Demo{Frames: 200, Width: 640, Height: 360}); err != nil {
		t.Fatal(err)
	}
	for _, v := range s.f.people {
		if v.wounds > 0 {
			return s, v
		}
	}
	t.Fatal("no one was struck")
	return nil, nil
}

// A struck villager takes out an axe and strikes back, and blood is spilt.
func TestAStruckVillagerStrikesBack(t *testing.T) {
	s, v := fighting(t)
	f := s.f
	if v.angry <= 0 || v.kit.slot != slotAxe || v.kit.axe == noAxe {
		t.Errorf("angry %v slot %d axe %d", v.angry, v.kit.slot, v.kit.axe)
	}
	if f.vit.Health >= 62 || f.blood.Live() == 0 {
		t.Errorf("health %v, blood %d", f.vit.Health, f.blood.Live())
	}
}

// Far from the player, an angry villager gives up.
func TestAnAngryVillagerGivesUpWhenThePlayerIsFar(t *testing.T) {
	s, v := fighting(t)
	f := s.f
	f.hands.ctl = still{}
	f.player.Transform.Position = vector.Vector3{X: v.obj.Transform.Position.X + 2*npcLeash, Y: v.obj.Transform.Position.Y}
	v.mind.control(f, v.obj, 0.1)
	if v.angry != 0 {
		t.Errorf("still angry: %v", v.angry)
	}
}

// A villager struck dead is gone like one dead of old age: out of the
// people, dead in the family tree, their partner free.
func TestAKilledVillagerIsGoneAndTheirPartnerFree(t *testing.T) {
	s, v := fighting(t)
	f := s.f
	mate := f.people[0]
	if mate == v {
		mate = f.people[1]
	}
	v.partner, mate.partner = mate, v
	f.wound(v, npcHealth, f.player.Transform.Position)
	if f.personByID(v.id) != nil || v.obj != nil || mate.partner != nil {
		t.Error("the dead are still about")
	}
	for _, e := range f.family.Entries() {
		if e.ID == v.id && !e.Dead {
			t.Error("alive in the family tree")
		}
	}
}

// The axe hurts more than bare hands, and the better it is the more.
func TestAnAxeHurtsMoreThanFists(t *testing.T) {
	k := &kit{axe: noAxe, slot: slotAxe}
	fist := k.blow()
	k.axe = 0
	stone := k.blow()
	k.axe = 3
	if !(fist < stone && stone < k.blow()) {
		t.Errorf("%v %v %v", fist, stone, k.blow())
	}
	k.slot = slotHoe
	if k.blow() != fist {
		t.Error("a hoe in hand is bare hands")
	}
}

// Being hit ends a talk and stands a seated player up.
func TestBeingHitStandsThePlayerUpAndEndsTheTalk(t *testing.T) {
	_, f := newGame(idle{})
	from := f.player.Transform.Position
	f.seat = &seat{from: from}
	f.talk = &talk{}
	f.hurtPlayer(10, from.Add(vector.Vector3{X: tile}))
	if f.seat != nil || f.talk != nil || f.vit.Health != 90 || f.flash <= 0 {
		t.Errorf("seat %v talk %v health %v flash %v", f.seat, f.talk, f.vit.Health, f.flash)
	}
}

// With no health left the player dies: with no heir the line ends, with
// the next one at their best.
func TestDyingOfABlowEndsTheLineWithNoHeir(t *testing.T) {
	_, f := newGame(idle{})
	f.hurtPlayer(1000, f.player.Transform.Position)
	if f.question == nil || f.vit.Health != 100 {
		t.Errorf("question %v health %v", f.question, f.vit.Health)
	}
}

// Eating a crop in hand feeds the player, one at a time, and a full
// player eats nothing.
func TestEatingFeedsThePlayer(t *testing.T) {
	_, f := newGame(idle{})
	f.bag.Add(cropOf(2), 2)
	f.pack.Hotbar[4] = cropOf(2)
	if !f.hold(4) || f.slot != slotFood {
		t.Fatalf("the crop was not taken in hand: slot %d", f.slot)
	}
	f.vit.Food = 30
	col, row := f.cellOf(f.player.Transform.Position)
	if f.useAs(f.kit, col, row) != actEat || f.vit.Food != 30+crops[2].nutrition || f.bag.Count(cropOf(2)) != 1 {
		t.Errorf("food %v bag %d", f.vit.Food, f.bag.Count(cropOf(2)))
	}
	f.vit.Food = 100
	if f.useAs(f.kit, col, row) != actNone || f.bag.Count(cropOf(2)) != 1 {
		t.Error("ate when full")
	}
}

// Sleeping mends the player.
func TestSleepingMendsThePlayer(t *testing.T) {
	_, f := newGame(idle{})
	f.vit.Health, f.vit.Food = 30, 20
	f.rest()
	if f.vit.Health <= 30 {
		t.Errorf("health %v", f.vit.Health)
	}
}
