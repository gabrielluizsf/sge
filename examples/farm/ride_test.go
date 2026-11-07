package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

// aHorse is a game with a horse bought, the player beside it, and how to
// step it by seconds, driven by next.
func aHorse(t *testing.T, next *intent) (*farm, *creature, func(seconds float64)) {
	t.Helper()
	g, f := build(SVGArt(), press{next})
	f.clock.Hour = 10
	f.kit.money = 1000
	if !f.buyAnimal(f.kit, stableOffers[1]) {
		t.Fatal("no horse")
	}
	c := f.herd.crits[len(f.herd.crits)-1]
	f.setFeet(c.pos.Add(vector.Vector3{X: tile}))
	return f, c, func(seconds float64) {
		for range int(seconds / frame) {
			g.Step(frame)
		}
	}
}

// E or Space by the horse gets the player on it, by a motion that takes its
// time, which ends with them in the saddle, on its very spot.
func TestGettingOnAHorse(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	words, do, ok := f.interaction(f.player)
	if !ok || words == "" {
		t.Fatal("nothing offered by the horse")
	}
	do()
	if f.ride == nil || !f.ride.changing || f.riding(f.player) {
		t.Fatal("not getting on")
	}
	if got := f.player.Animator.Current(); got != f.hands.pose(f.player)+"_mount" {
		t.Errorf("playing %s", got)
	}
	step(MountSeconds + 0.2)
	if !f.riding(f.player) {
		t.Fatal("not in the saddle")
	}
	if d := feet(f.player).Sub(c.pos).Length(); d > 0.01 {
		t.Errorf("%v from the horse", d)
	}
}

// Space by the horse does the same.
func TestSpaceGetsOnAHorse(t *testing.T) {
	var next intent
	f, _, step := aHorse(t, &next)
	next = intent{use: true}
	step(3 * frame)
	if f.ride == nil {
		t.Error("not getting on")
	}
}

// Under a rider a horse goes where they go, faster than on foot, its legs
// going with it, and stays under them.
func TestAHorseGoesWithItsRider(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.getOn(c)
	step(MountSeconds + 0.2)
	f.setFeet(f.ground.CellCenter(25, pathRow)) // the farm's open path, running west
	step(2 * frame)
	start := feet(f.player)
	for range 30 {
		next = intent{dir: f.onGround(vector.Vector3{X: -1, Y: 0})}
		step(frame)
		if c.pos.Sub(feet(f.player)).Length() > float64(rideRun*tile)*frame*2 {
			t.Fatalf("the horse is at %v, the rider at %v", c.pos, feet(f.player))
		}
	}
	went := start.Sub(feet(f.player)).Length() / tile / (30 * frame)
	if went <= walkSpeed/tile || went > rideWalk*1.1 {
		t.Errorf("%v cells a second, on foot %v", went, float64(walkSpeed)/tile)
	}
	if !c.walking && c.pace() <= 1 {
		t.Errorf("pace %v", c.pace())
	}
}

// E gets the player down, by the same motion backward, beside the horse; the
// horse then walks back to its pen by itself, and lives there as before.
func TestGettingDownSendsTheHorseHome(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.getOn(c)
	step(MountSeconds + 0.2)
	col0, row0 := c.pen.r.col+c.pen.r.cols/2, c.pen.r.row+c.pen.r.rows+2
	f.setFeet(f.ground.CellCenter(col0, row0))
	step(2 * frame)
	f.getDown()
	if got := f.player.Animator.Current(); got != f.hands.pose(f.player)+"_dismount" {
		t.Errorf("playing %s", got)
	}
	step(MountSeconds + 0.2)
	if f.ride != nil || !c.homing || c.held {
		t.Fatalf("ride %+v homing %v held %v", f.ride, c.homing, c.held)
	}
	if d := feet(f.player).Sub(c.pos).Length(); d < 0.3*tile {
		t.Errorf("came down on the horse: %v", d)
	}
	step(30)
	col, row := f.cellOf(c.pos)
	if c.homing || !c.pen.r.inside().has(col, row) {
		t.Errorf("the horse at %d, %d, homing %v, pen %+v", col, row, c.homing, c.pen.r)
	}
}

// A horse with no one looking is in its pen at once.
func TestAHorseFarAwayIsHomeAtOnce(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	c.homing = true
	c.pos = c.pos.Add(vector.Vector3{X: 400 * tile})
	step(2 * frame)
	col, row := f.cellOf(c.pos)
	if c.homing || !c.pen.r.inside().has(col, row) {
		t.Errorf("at %d, %d, homing %v", col, row, c.homing)
	}
}

// A rider stops at the water's edge: a horse is not taken swimming.
func TestAHorseStopsAtTheWater(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.getOn(c)
	step(MountSeconds + 0.2)
	col, row := shoreSpot(t, f)
	f.setFeet(f.ground.CellCenter(col, row))
	step(frame)
	for range 60 {
		next = intent{dir: vector.Vector3{X: 1}}
		step(frame)
		if f.inWater(feet(f.player)) || f.hands.swims {
			t.Fatal("rode into the water")
		}
	}
}

// shoreSpot is a dry cell with water just east of it.
func shoreSpot(t *testing.T, f *farm) (int, int) {
	t.Helper()
	for row := 0; row < f.ground.Rows; row++ {
		for col := 0; col < f.ground.Cols-2; col++ {
			if f.open(col, row) && !f.isWater(col, row) && f.isWater(col+1, row) {
				return col, row
			}
		}
	}
	t.Fatal("no shore")
	return 0, 0
}

// Dying on a horse does not leave a rider behind, and a save made in the
// saddle has the horse in its pen, as no riding is kept.
func TestARiderDiesAndTheHorseGoesHome(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.getOn(c)
	step(MountSeconds + 0.2)
	f.fall()
	if f.ride != nil || !c.homing || c.held {
		t.Errorf("ride %+v homing %v held %v", f.ride, c.homing, c.held)
	}
}

// Only one's own horse is ridden, not those for sale at a stable.
func TestAStablesHorseIsNotRidden(t *testing.T) {
	var next intent
	f, c, _ := aHorse(t, &next)
	c.pen = &pen{kind: "horse", r: c.pen.r, town: f.towns.Of(1, 0)}
	if f.horseBy(f.player) != nil {
		t.Error("rode a horse for sale")
	}
}

// A pen given to a new animal has its gate free and a way to it from the
// player's door, and one moved may not turn its gate to a fence.
func TestAPenCanBeWalkedInto(t *testing.T) {
	_, f := newGame(idle{})
	r, ok := f.freePenSpot()
	if !ok {
		t.Fatal("no ground")
	}
	if !f.reachable(r) || !f.open(r.col+r.cols/2, r.row+r.rows) {
		t.Errorf("pen %+v can not be reached", r)
	}
	// Against the farm's north fence, the gate would open on it.
	shut := rect{20, -penRows, penCols, penRows}
	if f.penFits(shut, nil) {
		t.Errorf("a pen at %+v opens on the fence", shut)
	}
}

// A horse too far off to find its way still walks home, in sight.
func TestAHorseWithNoWayWalksStraightHome(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.setFeet(c.pos)
	start := c.pen.r.middle().Add(vector.Vector3{X: 6 * tile})
	c.pos, c.homing, c.route, c.to = start, true, nil, c.pen.r.middle()
	step(1)
	if d := c.pos.Sub(start).Length(); d <= 0 {
		t.Errorf("it stood still, or jumped: %v", d)
	}
	step(30)
	if c.homing {
		t.Errorf("not home: at %v", c.pos)
	}
}
