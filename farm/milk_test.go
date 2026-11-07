package farm

import (
	"encoding/json"
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

// aCow is a game with the player beside the first cow.
func aCow(t *testing.T) (*farm, *creature, func(seconds float64)) {
	t.Helper()
	g, f := newGame(idle{})
	if f.herd == nil || len(f.herd.crits) == 0 || f.herd.crits[0].kind != "cow" {
		t.Fatal("the farm has no cow")
	}
	c := f.herd.crits[0]
	f.player.Transform.Position.X, f.player.Transform.Position.Y = c.pos.X+tile, c.pos.Y
	step := func(seconds float64) {
		for range int(seconds / frame) {
			g.Step(frame)
		}
	}
	return f, c, step
}

// The cow gives a bottle a day: at her side the player sits, milks, gets
// up with one, and that is all until the next day.
func TestTheCowGivesOneBottleADay(t *testing.T) {
	f, c, step := aCow(t)
	words, do, ok := f.cowNear(f.player)
	if !ok || words == "" {
		t.Fatal("nothing to do by the cow")
	}
	do()
	if f.milking == nil || f.hands.acting <= 0 || !c.held {
		t.Fatalf("not milking: %+v", f.milking)
	}
	if got := f.player.Animator.Current(); got != f.milking.dir+"_milkdown" {
		t.Errorf("playing %s", got)
	}
	step(MilkDownSeconds + 0.5)
	if got := f.player.Animator.Current(); f.milking == nil || got != f.milking.dir+"_milk" {
		t.Errorf("not at it: playing %s", got)
	}
	if f.bag.Count(milk) != 0 {
		t.Error("the bottle came before the milking was done")
	}
	step(milkTotal)
	if f.milking != nil || c.held || f.bag.Count(milk) != 1 {
		t.Fatalf("milking %+v, held %v, %d bottles", f.milking, c.held, f.bag.Count(milk))
	}
	if f.startMilking(c) || f.useAs(f.kit, 0, 0) == actMilk {
		t.Error("milked twice in a day")
	}
	if _, _, ok := f.cowNear(f.player); !ok {
		t.Error("nothing says she was milked today")
	}
	f.clock.Day++
	if !f.startMilking(c) {
		t.Error("no milk the next day")
	}
}

// Space by the cow milks her, with whatever is in hand.
func TestSpaceMilksTheCow(t *testing.T) {
	f, _, _ := aCow(t)
	col, row := f.cellOf(feet(f.player))
	if got := f.useAs(f.kit, col, row); got != actMilk {
		t.Errorf("did %v", got)
	}
}

// Far from any cow nothing happens.
func TestNoMilkingFarFromACow(t *testing.T) {
	f, c, _ := aCow(t)
	f.player.Transform.Position.X = c.pos.X + 10*tile
	if _, _, ok := f.cowNear(f.player); ok {
		t.Error("milking from far off")
	}
}

// A blow breaks the milking off with no milk, and she may be milked again.
func TestABlowBreaksTheMilkingOff(t *testing.T) {
	f, c, step := aCow(t)
	f.milkNow(c)
	step(1)
	f.hurtPlayer(5, feet(f.player).Add(vector.Vector3{X: tile}))
	if f.milking != nil || c.held || f.bag.Count(milk) != 0 || f.milkedToday(c) {
		t.Errorf("milking %+v, held %v, bottles %d", f.milking, c.held, f.bag.Count(milk))
	}
}

// Wherever she stands in the pen, the place to milk her is inside it.
func TestTheMilkersPlaceIsInThePen(t *testing.T) {
	f, c, _ := aCow(t)
	r := c.pen.r.inside()
	for _, at := range [][2]float64{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0.5, 0.5}} {
		c.pos = vector.Vector3{X: (float64(r.col) + 0.5 + float64(r.cols-1)*at[0]) * tile, Y: (float64(r.row) + 0.5 + float64(r.rows-1)*at[1]) * tile}
		start := c.pos
		_, seat := f.herd.milkSeat(c)
		col, row := f.cellOf(seat)
		if !r.has(col, row) {
			t.Errorf("cow at %v: the seat %v is off the field", at, seat)
		}
		if c.pos.Sub(start).Length() > 0.6*tile {
			t.Errorf("cow at %v moved %v", at, c.pos.Sub(start).Length())
		}
	}
}

// The day she was milked and the herd itself are kept in a save.
func TestTheHerdIsSaved(t *testing.T) {
	f, c, step := aCow(t)
	f.milkNow(c)
	step(milkTotal + 0.1)
	data, err := json.Marshal(takeSnapshot(f))
	if err != nil {
		t.Fatal(err)
	}
	var s snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	_, g := newGame(idle{})
	if err := restore(g, &s); err != nil {
		t.Fatal(err)
	}
	if n := len(g.herd.crits); n != 1 || !g.milkedToday(g.herd.crits[0]) || g.herd.crits[0].pos != c.pos {
		t.Errorf("%d animals: %+v", n, g.herd.crits[0])
	}
}

// A save of a kind unknown, or an impossible coat or milking, is refused.
func TestABadHerdIsRefused(t *testing.T) {
	f, _, _ := aCow(t)
	for _, bad := range []creatureState{{Kind: "dragon"}, {Kind: "cow", Coat: 9}, {Kind: "cow", Milked: 1 << 30}, {Kind: "cow", Coat: -1}} {
		s := takeSnapshot(f)
		s.Herd = []creatureState{bad}
		_, g := newGame(idle{})
		if err := restore(g, &s); err == nil {
			t.Errorf("%+v was let in", bad)
		}
	}
}
