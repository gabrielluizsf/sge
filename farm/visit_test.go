package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// visitor is a grown villager of town a who goes to the next town on some
// day, with the day set, and the hour of the middle of their stay.
func visitor(t *testing.T, f *farm, a *town.Town) (*villager, float64) {
	t.Helper()
	f.settle(a)
	for d := range 80 {
		f.clock.Day = d
		for _, p := range f.people {
			if p.town == a && p.stage() >= life.Adult {
				if s, ok := f.planOf(p).visit(); ok {
					return p, (s.from + s.to) / 2
				}
			}
		}
	}
	t.Fatal("no one visits in 80 days")
	return nil, 0
}

// A visitor is in the next town while they stay, to be met there, and is
// not in their own; and they are back in their own at the end of the day.
func TestAVisitorIsOnlyInTheTownTheyVisit(t *testing.T) {
	g, f := build(SVGArt(), still{})
	a := townWithRoad(t, f)
	b := f.roadFrom(a).to
	v, mid := visitor(t, f, a)
	f.settle(b)
	f.clock.Hour = mid
	f.player.Transform.Position = f.center(pathfind.Cell(b.Center))
	cv, err := canvas.NewCanvas(1280, 720)
	if err != nil {
		t.Fatal(err)
	}
	g.Draw(cv)
	for range 40 {
		g.Step(frame)
	}
	if !f.keeper.Awake(b) {
		t.Fatal("the next town is not awake")
	}
	if v.obj == nil || v.guest != b || v.place() != b {
		t.Fatalf("not in the town they visit: obj %v guest %v", v.obj != nil, v.guest != nil)
	}
	if f.keeper.Awake(a) {
		t.Fatal("their own town is awake too; the test is not far enough")
	}
	f.clock.Hour = bedBy - 0.5
	for range 40 {
		g.Step(frame)
	}
	if v.guest != nil || v.obj != nil && f.keeper.Awake(b) {
		t.Errorf("late in the evening still visiting: guest %v", v.guest != nil)
	}
}

// Two villagers who meet stop and face each other, and a ware changes
// hands for its price with a note over each; the money of the two is the
// same after.
func TestTwoVillagersWhoMeetTrade(t *testing.T) {
	f, tn, step := seenTown(t, SVGArt())
	var pair []*villager
	for _, v := range f.people {
		if v.town == tn && v.stage() >= life.Adult && v.obj != nil && len(pair) < 2 {
			pair = append(pair, v)
		}
	}
	if len(pair) < 2 {
		t.Fatal("no two villagers about")
	}
	a, b := pair[0], pair[1]
	a.inside, b.inside = false, false
	a.kit.bag.Add(wood, 1)
	b.kit.money = 20
	before := a.kit.money + b.kit.money
	b.obj.Transform.Position = a.obj.Transform.Position.Add(vector.Vector3{X: tile})
	floats := len(f.floats)
	f.meet(a, b)
	if a.mind.hold <= 0 || b.mind.hold <= 0 || a.mind.facing == "" || a.mind.facing == b.mind.facing {
		t.Errorf("they did not stop and face each other: %v %v", a.mind.facing, b.mind.facing)
	}
	if a.kit.money+b.kit.money != before || a.kit.money != 15+woodPrice || b.kit.bag.Count(wood) != 1 {
		t.Errorf("a $%d, b $%d, %d wood", a.kit.money, b.kit.money, b.kit.bag.Count(wood))
	}
	if len(f.floats) != floats+2 {
		t.Errorf("%d notes, want 2", len(f.floats)-floats)
	}
	step(1)
}
