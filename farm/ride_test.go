package farm

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

// walkUp is more than a walk up to a horse, from the farthest one gets on at.
const walkUp = 1.0

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
	step(walkUp)
	if got := f.player.Animator.Current(); got != f.hands.pose(f.player)+"_mount" {
		t.Errorf("playing %s", got)
	}
	step(MountSeconds + walkUp)
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
	step(MountSeconds + walkUp)
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

// E gets the player down, by jumping off its side, beside the horse; the
// horse then walks back to its pen by itself, and lives there as before.
func TestGettingDownSendsTheHorseHome(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.getOn(c)
	step(MountSeconds + walkUp)
	col0, row0 := c.pen.r.col+c.pen.r.cols/2, c.pen.r.row+c.pen.r.rows+2
	f.setFeet(f.ground.CellCenter(col0, row0))
	step(2 * frame)
	f.getDown()
	if got := f.player.Animator.Current(); got != f.hands.pose(f.player)+"_dismount" {
		t.Errorf("playing %s", got)
	}
	step(DismountSeconds + 2*frame)
	if f.ride != nil || !c.homing || c.held {
		t.Fatalf("ride %+v homing %v held %v", f.ride, c.homing, c.held)
	}
	if d := feet(f.player).Sub(c.pos).Length(); math.Abs(d-MountAside*tile) > 0.5*tile {
		t.Errorf("came down %v from the horse", d)
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
	step(MountSeconds + walkUp)
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
	step(MountSeconds + walkUp)
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

// One walks up to a horse before getting on, and is drawn on its spot from then
// on: the motion's picture has them beside it, where they walked to.
func TestWalkingUpToAHorse(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	start := feet(f.player)
	f.getOn(c)
	spot := f.besideHorse(c)
	step(2 * frame)
	if d := feet(f.player).Sub(spot).Length(); d >= start.Sub(spot).Length() || feet(f.player) == c.pos {
		t.Errorf("not walking up: %v from the spot", d)
	}
	if got := f.player.Animator.Current(); got != f.hands.pose(f.player)+"_walk" {
		t.Errorf("playing %s", got)
	}
	step(walkUp)
	if d := feet(f.player).Sub(c.pos).Length(); d > 0.01 {
		t.Errorf("drawn %v from the horse getting on", d)
	}
}

// A horse with its side for getting on against a wall turns about.
func TestAHorseTurnsWhenItsSideIsBlocked(t *testing.T) {
	var next intent
	f, c, _ := aHorse(t, &next)
	col, row := f.cellOf(f.besideHorse(c))
	rockCol, rockRow := rockCell(t, f)
	f.things.Set(col, row, f.thingAt(rockCol, rockRow))
	if !f.roomBeside(c) {
		t.Fatal("no room on either side")
	}
	if !f.openAt(f.besideHorse(c)) {
		t.Error("still blocked")
	}
}

// The spot beside a horse is at the same distance whichever way it faces, and a
// quarter turn from its heading in the world.
func TestTheSpotBesideAHorse(t *testing.T) {
	var next intent
	f, c, _ := aHorse(t, &next)
	for h := range 8 {
		c.heading = h
		side := f.besideHorse(c).Sub(c.pos)
		if d := side.Length(); math.Abs(d-MountAside*tile) > 0.01 {
			t.Errorf("heading %d: %v away", h, d)
		}
		ahead := f.onGround(headingWay(h))
		if d := side.X*ahead.X + side.Y*ahead.Y; math.Abs(d) > 0.01 {
			t.Errorf("heading %d: not square to the way it goes (%v)", h, d)
		}
	}
}

// Getting down starts the moment it is asked for, is soon over, and leaves the
// player free at once.
func TestGettingDownIsQuick(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.getOn(c)
	step(MountSeconds + walkUp)
	f.getDown()
	if f.ride == nil || !f.ride.changing || !f.ride.down {
		t.Fatal("not getting down at once")
	}
	if got := f.player.Animator.Current(); got != f.hands.pose(f.player)+"_dismount" {
		t.Errorf("playing %s", got)
	}
	step(DismountSeconds + 3*frame)
	if f.ride != nil {
		t.Error("still getting down")
	}
	if f.hands.acting > 0 {
		t.Errorf("not free: %v seconds left", f.hands.acting)
	}
	if DismountSeconds > 0.7 {
		t.Errorf("getting down takes %v seconds", DismountSeconds)
	}
}

// Nothing is written over a rider's head, and E still gets them down.
func TestNoWordsOverARider(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.getOn(c)
	step(MountSeconds + walkUp)
	words, do, ok := f.interaction(f.player)
	if words != "" {
		t.Errorf("%q written over the rider", words)
	}
	if !ok {
		t.Fatal("E does nothing on the horse")
	}
	do()
	if f.ride == nil || !f.ride.down {
		t.Error("E does not get them down")
	}
}

// One jumps down with one's back to the horse: facing away from it, from the
// spot beside it, as the motion's last picture has them.
func TestOneLandsFacingAwayFromTheHorse(t *testing.T) {
	var next intent
	f, c, step := aHorse(t, &next)
	f.getOn(c)
	step(MountSeconds + walkUp)
	f.getDown()
	spot := f.besideHorse(c)
	from := c.pos
	step(DismountSeconds + 2*frame)
	away := unitOf(spot.Sub(from), vector.Vector3{})
	if d := away.X*f.hands.facing.X + away.Y*f.hands.facing.Y; d < 0.99 {
		t.Errorf("facing %v, not away from the horse along %v", f.hands.facing, away)
	}
}
