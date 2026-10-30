package farm

import "testing"

// The player sleeps only at night, and wakes the next morning with the
// watered crops grown.
func TestThePlayerSleepsAtNight(t *testing.T) {
	_, f := newGame(idle{})
	p := f.plots[0]
	p.crop, p.watered = 0, true
	if f.sleep() {
		t.Fatal("slept in the morning")
	}
	f.clock.Hour = 21
	if !f.sleep() || f.clock.Day != 1 || f.clock.Hour != wakeAt || p.stage != 1 || p.watered {
		t.Fatalf("day %d hour %v, plot %+v", f.clock.Day, f.clock.Hour, p)
	}
}

// One still up a minute past midnight faints, loses all their money and
// wakes at home in the morning; one who went to bed loses nothing.
func TestStayingUpLosesTheMoney(t *testing.T) {
	g, f := newGame(idle{})
	f.money = 77
	f.clock.Hour = 23.99
	for range 60 * 3 {
		g.Step(frame)
	}
	if f.money != 0 || f.clock.Hour < wakeAt || f.room == nil || !f.room.own {
		t.Fatalf("money %d at %v, room %+v", f.money, f.clock.Hour, f.room)
	}
	f.money = 50
	f.clock.Hour = 23
	f.sleep()
	for range 60 * 3 {
		g.Step(frame)
	}
	if f.money != 50 {
		t.Errorf("lost the money after sleeping: %d", f.money)
	}
}

// Going into a house leaves the world outside waiting, with its people;
// the doormat leads back out to the doorstep.
func TestAHouseIsGoneIntoAndLeft(t *testing.T) {
	g, f := newGame(keyboard{})
	o := g.Find("farmer")
	d, ok := f.doorAt(homeDoor.Col, homeDoor.Row)
	if !ok || !d.own {
		t.Fatal("no door at home")
	}
	f.enter(d)
	g.Step(frame)
	if g.Find("wild") != nil || g.Find("ground") != nil || f.room == nil {
		t.Fatal("the world outside is still in the game")
	}
	o.Transform.Position = f.ground.CellCenter(matCol, roomRows-1)
	f.leave()
	g.Step(frame)
	if g.Find("wild") == nil || f.room != nil {
		t.Fatal("the world outside did not come back")
	}
	if c, r := f.cellOf(feet(o)); c != homeDoor.Col || r != homeDoor.Row {
		t.Errorf("came out at %d,%d", c, r)
	}
}

// A room is laid out anew each time it is gone into, with whoever is home
// then.
func TestARoomShowsWhoIsHomeNow(t *testing.T) {
	g, f := newGame(idle{})
	tn := eastTown(f)
	v := tradesman(f, tn, roleLumberjack)
	d, _ := f.doorAt(v.home.Door.Col, v.home.Door.Row)
	v.inside = true
	f.enter(d)
	g.Step(frame)
	if len(f.residents) != 1 {
		t.Fatalf("%d at home", len(f.residents))
	}
	f.leave()
	g.Step(frame)
	v.inside = false
	f.enter(d)
	g.Step(frame)
	if len(f.residents) != 0 || g.Find("villager") != nil {
		t.Errorf("%d at home after they went out", len(f.residents))
	}
}

// In their house, the player is offered sleep only beside the bed.
func TestSleepIsOfferedBesideTheBed(t *testing.T) {
	g, f := newGame(idle{})
	o := g.Find("farmer")
	f.enter(houseDoor{at: homeDoor, faceOut: "down", seed: seed, own: true})
	g.Step(frame)
	if words, _, _ := f.interaction(o); words == "E: dormir" {
		t.Fatal("offered sleep on the doormat")
	}
	for r := range roomRows {
		for c := range roomCols {
			if id := f.room.things.At(c, r); id == tileID["bed"] || id == tileID["bed_blue"] {
				o.Transform.Position = f.ground.CellCenter(c, r+1)
			}
		}
	}
	if words, _, ok := f.interaction(o); !ok || words != "E: dormir" {
		t.Errorf("beside the bed: %q", words)
	}
}
