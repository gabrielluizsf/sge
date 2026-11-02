package town

import "testing"

// bigLayout has houses of several sizes, up to 8 by 7 cells.
func bigLayout() *Layout {
	sizes := []Cell{{8, 6}, {7, 5}, {7, 5}, {7, 6}, {7, 6}, {8, 7}}
	return &Layout{Seed: 11, Region: 40, Chance: 0.6, Kinds: len(sizes), Gap: 2,
		Footprint: func(k int) Cell { return sizes[k] },
		Fixed:     map[Cell]Cell{{1, 0}: {58, 4}},
		Skip:      func(rx, ry int) bool { return rx == 0 && ry == 0 }}
}

// Houses stand on their real footprint, none touching another with a
// cell between them at least, none on the street, the plaza or the
// stall, anvil and forge; each one's door is on its front, in the yard.
func TestBigHousesDoNotOverlapOrTouch(t *testing.T) {
	l := bigLayout()
	for ry := -4; ry < 4; ry++ {
		for rx := -4; rx < 4; rx++ {
			tn := l.Of(rx, ry)
			if tn == nil {
				continue
			}
			for i, a := range tn.Houses {
				if a.Width < 7 || a.Depth < 5 {
					t.Fatalf("house %+v lost its footprint", a)
				}
				for x := a.Col - 1; x <= a.Col+a.Width; x++ {
					for y := a.Row - 1; y <= a.Row+a.Depth; y++ {
						if o := tn.House(x, y); o != nil && o != &tn.Houses[i] {
							t.Fatalf("houses %+v and %+v touch", a, *o)
						}
					}
				}
				for x := a.Col; x < a.Col+a.Width; x++ {
					for y := a.Row; y < a.Row+a.Depth; y++ {
						if tn.Road(x, y) || !tn.Grounds(x, y, 0) || l.At(x, y) != tn {
							t.Fatalf("house %+v strays at %d,%d", a, x, y)
						}
						if (Cell{x, y}) == tn.Stall || (Cell{x, y}) == tn.Anvil || (Cell{x, y}) == tn.Forge {
							t.Fatalf("house %+v on the stall's ground", a)
						}
					}
				}
				if a.Door.Col < a.Col || a.Door.Col >= a.Col+a.Width || tn.House(a.Door.Col, a.Door.Row) != nil || tn.House(a.Garden.Col, a.Garden.Row) != nil {
					t.Fatalf("house %+v: its door is not on its front", a)
				}
			}
		}
	}
}

// A town stays within its region however large its houses.
func TestBigTownsStayInTheirRegion(t *testing.T) {
	l := bigLayout()
	if l.Lots() != 1 || l.Reach() != defaultLead+8 {
		t.Fatalf("lots %d, reach %d", l.Lots(), l.Reach())
	}
	for ry := -3; ry < 3; ry++ {
		for rx := -3; rx < 3; rx++ {
			tn := l.Of(rx, ry)
			if tn == nil {
				continue
			}
			for _, h := range tn.Houses {
				for _, c := range []Cell{{h.Col - 1, h.Row - 1}, {h.Col + h.Width, h.Row + h.Depth}} {
					if l.At(c.Col, c.Row) != tn {
						t.Fatalf("house %+v leaves region %v", h, tn.Region)
					}
				}
			}
		}
	}
}
