package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/daytime"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// seenTown is the east town, settled, with the player standing in it and
// the view drawn, so what is in sight is known.
func seenTown(t *testing.T, a Art) (f *farm, tn *town.Town, step func(n int)) {
	t.Helper()
	g, f := build(a, still{})
	tn = eastTown(f)
	f.player.Transform.Position = f.center(pathfind.Cell(tn.Center))
	cv, err := canvas.NewCanvas(1280, 720)
	if err != nil {
		t.Fatal(err)
	}
	g.Draw(cv)
	for range 3 {
		g.Step(frame) // the keeper wakes the town
	}
	return f, tn, func(n int) {
		for range n {
			g.Step(frame)
		}
	}
}

// Everyone has a bed to be in by ten, and is not at work at that hour: the
// schedules of every trade say sleep at 22, and a day to rise to.
func TestEveryTradeSleepsAt22AndRises(t *testing.T) {
	for role, s := range schedules {
		for _, h := range []float64{22, 23.5, 3, 5.9} {
			if got := s.At(clockAt(h)); got != doSleep {
				t.Errorf("%s at %v: %q, not sleep", role, h, got)
			}
		}
		if got := s.At(clockAt(12)); got == doSleep || got == "" {
			t.Errorf("%s at noon: %q", role, got)
		}
		if got := s.At(clockAt(bedBy - 0.1)); got == doSleep {
			t.Errorf("%s is in bed before %v", role, bedBy)
		}
	}
}

// Villagers in sight are given work in sight: a job that leads them out of
// the screen leaves them frozen at its edge, standing there all day. This
// is the bug that left the realistic towns' streets empty.
func TestWorkGivenToThoseInSightIsInSight(t *testing.T) {
	for name, a := range map[string]Art{"stylized": SVGArt(), "roomy": roomyArt{SVGArt()}} {
		f, tn, _ := seenTown(t, a)
		f.clock.Hour = 9
		for _, v := range f.people {
			if v.town != tn || v.stage() < life.Adult {
				continue
			}
			for range 40 {
				v.mind.job, v.mind.route = nil, nil
				if j := v.mind.next(f); j != nil && !f.inSight(f.center(j.to)) {
					t.Fatalf("%s: a %s is sent to %v, out of sight", name, v.role, j.to)
				}
			}
		}
	}
}

// Someone who goes out of sight when it is their night is home in bed, as
// if they had walked there unseen, and not frozen in the street.
func TestThoseOutOfSightAtNightAreHome(t *testing.T) {
	f, tn, _ := seenTown(t, SVGArt())
	v := tradesman(f, tn, roleFarmer)
	f.clock.Hour = 23
	f.player.Transform.Position = f.center(pathfind.Cell{Col: tn.Center.Col + 400, Row: tn.Center.Row})
	v.mind.bedUnseen(f, v.obj)
	if !v.inside || !v.obj.Hidden || v.mind.job != nil {
		t.Fatalf("inside %v, hidden %v, job %v", v.inside, v.obj.Hidden, v.mind.job)
	}
	if c, r := f.cellOf(v.obj.Transform.Position); c != v.home.Door.Col || r != v.home.Door.Row {
		t.Errorf("at %d,%d, not their door", c, r)
	}
	v.inside = false
	f.clock.Hour = 12
	v.mind.bedUnseen(f, v.obj)
	if v.inside {
		t.Error("put to bed at noon")
	}
}

// A grown wanderer works, felling trees and selling the wood at the stall
// for money; a child does not.
func TestAGrownWandererWorksForMoney(t *testing.T) {
	f, tn, _ := seenTown(t, SVGArt())
	v := tradesman(f, tn, roleWanderer)
	if v == nil {
		v = tradesman(f, tn, roleFarmer)
		v.role = roleWanderer
	}
	v.who.Stage = life.Adult
	v.kit.bag.Add(wood, oddLoad)
	j := v.mind.oddJobs(f)
	if j == nil || j.then == nil {
		t.Fatal("no trip to the stall with a load of wood")
	}
	before := v.kit.money
	j.then(f)
	if v.kit.money != before+oddLoad*woodPrice || v.kit.bag.Count(wood) != 0 {
		t.Errorf("money %d to %d, %d wood left", before, v.kit.money, v.kit.bag.Count(wood))
	}
	v.who.Stage = life.Child
	if v.mind.oddJobs(f) != nil {
		t.Error("a child has work")
	}
}

// Every grown villager starts with some money.
func TestGrownVillagersStartWithMoney(t *testing.T) {
	f, tn, _ := seenTown(t, SVGArt())
	for _, v := range f.people {
		if v.town == tn && v.stage() >= life.Adult && v.kit.money <= 0 {
			t.Errorf("a %s with no money", v.role)
		}
	}
}

// What a villager at home in bed looks like when someone goes in: lying
// in the bed's cells, asleep, once, and the bed's own tile away when the
// art draws it with them.
func TestASleeperLiesInTheBed(t *testing.T) {
	f, tn, step := seenTown(t, SVGArt())
	v := tradesman(f, tn, roleLumberjack)
	v.who.Stage = life.Adult
	v.partner = nil // a couple lie in the double bed: see the next test
	f.clock.Hour = 23
	v.inside = true
	d, _ := f.doorAt(v.home.Door.Col, v.home.Door.Row)
	f.enter(d)
	step(2)
	i, used := usedBed(f.room)
	if len(f.residents) != 1 || !used {
		t.Fatalf("%d residents, bed used %v", len(f.residents), used)
	}
	o := f.residents[0].obj
	if want := f.room.beds[i].center(f); o.Transform.Position != want {
		t.Errorf("at %v, the bed's middle is %v", o.Transform.Position, want)
	}
	if o.Animator == nil || o.Animator.Current() != AnimSleep {
		t.Errorf("playing %v", o.Animator)
	}
	if b := f.room.beds[i]; f.room.things.At(b.c0, b.r0) == tilemap.Empty {
		t.Error("the bed's tile is gone with a stylized sleeper")
	}
}

// Awake in the evening, they are not put in bed; and a sleeper in a look
// whose bed is baked with them takes the bed's tile away, but only once
// their frames are baked.
func TestTheBedGoesOnlyWithABakedSleeper(t *testing.T) {
	for _, baked := range []bool{false, true} {
		f, tn, step := seenTown(t, bedArt{SVGArt(), baked})
		v := tradesman(f, tn, roleLumberjack)
		v.who.Stage = life.Adult
		v.partner = nil // a double bed keeps its tile, its two sides drawn on it
		f.clock.Hour = 20
		v.inside = true
		d, _ := f.doorAt(v.home.Door.Col, v.home.Door.Row)
		f.enter(d)
		step(2)
		if _, used := usedBed(f.room); used {
			t.Fatal("in bed at eight")
		}
		f.leave()
		f.clock.Hour = 23
		f.enter(d)
		step(2)
		i, used := usedBed(f.room)
		b := f.room.beds[i]
		gone := f.room.things.At(b.c0, b.r0) == tilemap.Empty
		if used != baked || gone != baked {
			t.Errorf("baked %v: in bed %v, tile gone %v", baked, used, gone)
		}
	}
}

// bedArt is an art that draws beds with their sleepers, whose people's
// sleep is baked or not.
type bedArt struct {
	Art
	baked bool
}

func (bedArt) SleepBaked() bool { return true }

func (a bedArt) Person(w Who) Body { return sleepyBody{a.Art.Person(w), a.baked} }

type sleepyBody struct {
	Body
	baked bool
}

func (b sleepyBody) Sleeps() bool { return b.baked }

// clockAt is a clock at hour h.
func clockAt(h float64) *daytime.Clock { return &daytime.Clock{Hour: h} }

// usedBed is the index of the first bed of r that someone lies in.
func usedBed(r *room) (int, bool) {
	for i, b := range r.beds {
		if b.used {
			return i, true
		}
	}
	return 0, false
}
