package town

import (
	"testing"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

func layout() *Layout {
	return &Layout{Seed: 11, Region: 40, Chance: 0.5, Kinds: 4,
		Fixed: map[Cell]Cell{{1, 0}: {58, 4}},
		Skip:  func(rx, ry int) bool { return rx == 0 && ry == 0 }}
}

// The same seed lays out the same towns, whatever is asked first; some
// regions have one and some none; skipped regions none, fixed ones theirs.
func TestTownsAreLaidOutTheSameEveryTime(t *testing.T) {
	a, b := layout(), layout()
	b.Of(5, 5)
	towns := 0
	for ry := -6; ry < 6; ry++ {
		for rx := -6; rx < 6; rx++ {
			ta, tb := a.Of(rx, ry), b.Of(rx, ry)
			if (ta == nil) != (tb == nil) || ta != nil && (ta.Center != tb.Center || len(ta.Houses) != len(tb.Houses)) {
				t.Fatalf("region %d,%d differs", rx, ry)
			}
			if ta != nil {
				towns++
				if a.Of(rx, ry) != ta {
					t.Fatal("a town was planned twice")
				}
			}
		}
	}
	if towns < 30 || towns > 110 {
		t.Errorf("%d towns in 144 regions", towns)
	}
	if a.Of(0, 0) != nil {
		t.Error("a town in a skipped region")
	}
	if f := a.Of(1, 0); f == nil || f.Center != (Cell{58, 4}) {
		t.Errorf("the fixed town is %+v", f)
	}
}

// Every town stays in its region; its houses stand off the road, on its
// grounds, none on another; each house's doorstep and garden lie in front
// of it, between it and the street; the stall is beside the plaza.
func TestATownsPlanHoldsTogether(t *testing.T) {
	l := layout()
	for ry := -4; ry < 4; ry++ {
		for rx := -4; rx < 4; rx++ {
			tn := l.Of(rx, ry)
			if tn == nil {
				continue
			}
			for dy := -6; dy <= 6; dy++ {
				for dx := -12; dx <= 12; dx++ {
					c, r := tn.Center.Col+dx, tn.Center.Row+dy
					if tn.Grounds(c, r, 1) && l.At(c, r) != tn {
						t.Fatalf("town %v reaches out of its region at %d,%d", tn.Region, c, r)
					}
					if tn.Road(c, r) && tn.House(c, r) != nil {
						t.Fatalf("a house on the road at %d,%d", c, r)
					}
				}
			}
			if len(tn.Houses) == 0 {
				t.Errorf("town %v has no houses", tn.Region)
			}
			for _, h := range tn.Houses {
				for _, c := range []Cell{h.Door, h.Garden} {
					if tn.House(c.Col, c.Row) != nil || tn.Road(c.Col, c.Row) || !tn.Grounds(c.Col, c.Row, 0) {
						t.Fatalf("house %+v: its yard cell %v is not clear", h, c)
					}
					if !tn.Road(c.Col, tn.Center.Row) {
						t.Fatalf("house %+v is not on the street", h)
					}
				}
				if h.Kind < 0 || h.Kind >= 4 {
					t.Fatalf("house kind %d", h.Kind)
				}
			}
			for _, s := range []Cell{tn.Stall, tn.Anvil, tn.Forge} {
				if tn.Road(s.Col, s.Row) || tn.House(s.Col, s.Row) != nil || tn.House(s.Col, s.Row+1) != nil || !tn.Grounds(s.Col, s.Row, 0) {
					t.Errorf("%v is not beside the plaza", s)
				}
			}
			if tn.Stall == tn.Anvil || tn.Anvil == tn.Forge || !tn.Road(tn.Stall.Col, tn.Stall.Row+1) || !tn.Road(tn.Anvil.Col, tn.Anvil.Row+1) {
				t.Error("the stall, anvil and forge are not each beside the plaza")
			}
		}
	}
}

// Near finds the towns around a cell, across regions.
func TestNearFindsTheTownsAround(t *testing.T) {
	l := layout()
	near := l.Near(50, 4, 10)
	if len(near) != 1 || near[0].Center != (Cell{58, 4}) {
		t.Errorf("near 50,4: %v", near)
	}
	if len(l.Near(20, 4, 5)) != 0 {
		t.Error("a town near the skipped region's middle")
	}
}

// The keeper wakes a town as the camera nears it and lets its objects go
// when the camera leaves, waking it anew when it comes back.
func TestTheKeeperWakesTownsNearTheCamera(t *testing.T) {
	g := game.New()
	woken := 0
	k := &Keeper{Layout: layout(), Range: 20, Cell: 10, Wake: func(tn *Town) []*game.Object {
		woken++
		return []*game.Object{game.NewObject("villager")}
	}}
	o := game.NewObject("keeper")
	o.AddScript(k)
	g.Add(o)
	g.Camera = vector.Vector3{X: 300, Y: 40} // cell 30, 4: 28 from the town
	g.Step(1.0 / 60)
	if woken != 0 || g.Find("villager") != nil {
		t.Fatal("woke a town out of range")
	}
	g.Camera.X = 450 // cell 45: 13 from it
	g.Step(1.0 / 60)
	g.Step(1.0 / 60)
	if woken != 1 || g.Find("villager") == nil || !k.Awake(k.Layout.Of(1, 0)) {
		t.Fatalf("woke %d times", woken)
	}
	g.Camera.X = 360 // 22 cells: within half as far again
	g.Step(1.0 / 60)
	if g.Find("villager") == nil {
		t.Fatal("slept too soon")
	}
	g.Camera.X = 0
	g.Step(1.0 / 60)
	g.Step(1.0 / 60)
	if g.Find("villager") != nil {
		t.Fatal("the town did not sleep")
	}
	g.Camera.X = 560
	g.Step(1.0 / 60)
	if woken != 2 {
		t.Errorf("woke %d times", woken)
	}
}
