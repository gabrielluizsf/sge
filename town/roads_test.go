package town

import "testing"

// roomy is a layout of the realistic look: large houses, wide gaps.
func roomy() *Layout {
	l := &Layout{Seed: 77, Region: 40, Chance: 0.45, Kinds: 4, Gap: 3, Lead: 7, Margin: 3,
		Footprint: func(kind int) Cell { return Cell{Col: 4 + kind%2, Row: 3} },
		Fixed:     map[Cell]Cell{{1, 0}: {58, 4}, {-1, 0}: {-22, 4}},
		Skip:      func(rx, ry int) bool { return rx == 0 && ry == 0 }}
	l.Roads = []Rect{{30, 4, 58 - l.Reach(), 4}, {-22 + l.Reach(), 4, -1, 4}}
	l.Keep = []Rect{{-16, -16, 46, 36}}
	return l
}

func paved(l *Layout, c, r int) bool {
	if l.OnRoad(c, r, 0) {
		return true
	}
	t := l.At(c, r)
	return t != nil && t.Road(c, r)
}

// townAt is the town whose plaza's middle is cell c, r, or nil.
func townAt(l *Layout, c, r int) *Town {
	t := l.At(c, r)
	if t != nil && t.Center == (Cell{c, r}) {
		return t
	}
	return nil
}

// reaches walks the paved cells from town t's middle and reports whether
// another town's middle is on the way.
func reaches(l *Layout, t *Town) bool {
	seen := map[Cell]bool{t.Center: true}
	todo := []Cell{t.Center}
	for len(todo) > 0 {
		c := todo[0]
		todo = todo[1:]
		if o := townAt(l, c.Col, c.Row); o != nil && o != t {
			return true
		}
		for _, d := range []Cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			n := Cell{c.Col + d.Col, c.Row + d.Row}
			if !seen[n] && paved(l, n.Col, n.Row) && abs(n.Col-t.Center.Col) < 400 && abs(n.Row-t.Center.Row) < 400 {
				seen[n] = true
				todo = append(todo, n)
			}
		}
	}
	return false
}

// Following the road out of any town ends in another town.
func TestEveryTownsRoadLeadsToAnotherTown(t *testing.T) {
	l := roomy()
	n := 0
	for ry := -7; ry <= 7; ry++ {
		for rx := -7; rx <= 7; rx++ {
			if tw := l.Of(rx, ry); tw != nil {
				n++
				if !reaches(l, tw) {
					t.Errorf("the road out of the town at %v leads nowhere", tw.Center)
				}
			}
		}
	}
	if n < 60 {
		t.Fatalf("only %d towns", n)
	}
}

// No road runs over a house, nor into the farm.
func TestRoadsKeepOffHousesAndFarm(t *testing.T) {
	l := roomy()
	for r := -300; r < 300; r++ {
		for c := -300; c < 300; c++ {
			if !l.OnRoad(c, r, 0) && !l.Fence(c, r) {
				continue
			}
			if c >= 0 && c < 30 && r >= 0 && r < 20 {
				t.Fatalf("road or fence in the farm at %d,%d", c, r)
			}
			if tw := l.At(c, r); tw != nil && tw.House(c, r) != nil {
				t.Fatalf("road or fence on a house at %d,%d", c, r)
			}
		}
	}
}

// A road is open to the wild on both sides every few steps: from any
// cell of it, a short walk over free cells reaches grass on either side.
func TestFencesLeaveTheWildReachableFromEveryRoadCell(t *testing.T) {
	l := roomy()
	fences, cells := 0, 0
	for ry := -4; ry <= 4; ry++ {
		for rx := -4; rx <= 4; rx++ {
			for _, r := range l.roadsOf(rx, ry) {
				for row := r.R0; row <= r.R1; row++ {
					for col := r.C0; col <= r.C1; col++ {
						cells++
						for _, side := range []int{-1, 1} {
							if !opens(l, r, col, row, side) {
								t.Fatalf("road cell %d,%d is walled in on side %d", col, row, side)
							}
						}
						if l.Fence(col, row+1) || l.Fence(col+1, row) {
							fences++
						}
					}
				}
			}
		}
	}
	if fences < cells/10 || fences > cells {
		t.Errorf("%d fenced of %d road cells", fences, cells)
	}
}

// opens walks over free cells from road cell col, row to a cell two or
// more cells off the roads, on the given side of r, in at most 20 steps.
func opens(l *Layout, r Rect, col, row, side int) bool {
	seen := map[Cell]bool{{col, row}: true}
	todo := []Cell{{col, row}}
	for step := 0; step < 20 && len(todo) > 0; step++ {
		var next []Cell
		for _, c := range todo {
			for _, d := range []Cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
				n := Cell{c.Col + d.Col, c.Row + d.Row}
				if seen[n] || l.Fence(n.Col, n.Row) {
					continue
				}
				seen[n] = true
				if r.R0 == r.R1 && (n.Row-r.R0)*side <= 0 && n.Row != r.R0 || r.C0 == r.C1 && r.R0 != r.R1 && (n.Col-r.C0)*side <= 0 && n.Col != r.C0 {
					continue // the wrong side
				}
				if !l.OnRoad(n.Col, n.Row, 2) {
					return true
				}
				next = append(next, n)
			}
		}
		todo = next
	}
	return false
}

// Fences come in runs, and not everywhere.
func TestFencesRunInStretches(t *testing.T) {
	l := roomy()
	r := Rect{C0: 100, R0: 100, C1: 300, R1: 100}
	runs, fenced := 0, 0
	for c := r.C0; c <= r.C1; c++ {
		if inRun(l.Seed, r, c, 99) && beside(r, c, 99) {
			fenced++
			if c == r.C0 || !(inRun(l.Seed, r, c-1, 99) && beside(r, c-1, 99)) {
				runs++
			}
		}
	}
	if runs < 5 || fenced < 40 || fenced > 150 {
		t.Errorf("%d runs, %d fenced of 201", runs, fenced)
	}
}

// The road query is the same whatever is asked first.
func TestRoadsDoNotDependOnTheOrderOfQuestions(t *testing.T) {
	a, b := roomy(), roomy()
	b.OnRoad(200, 200, 0)
	b.OnRoad(-150, 90, 0)
	for r := 60; r < 120; r++ {
		for c := 100; c < 300; c++ {
			if a.OnRoad(c, r, 2) != b.OnRoad(c, r, 2) || a.Fence(c, r) != b.Fence(c, r) {
				t.Fatalf("differs at %d,%d", c, r)
			}
		}
	}
}

func BenchmarkOnRoad(b *testing.B) {
	l := roomy()
	l.OnRoad(0, 0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.OnRoad(i%300, (i/300)%300, 2)
	}
}

func BenchmarkFence(b *testing.B) {
	l := roomy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Fence(i%300, (i/300)%300)
	}
}

// BenchmarkNewRegion is the cost of the first road question about a
// region, which plans the towns around it and their roads.
func BenchmarkNewRegion(b *testing.B) {
	l := roomy()
	for i := 0; i < b.N; i++ {
		l.OnRoad(40*(i+10), 40*((i*7)%1000+10), 0)
	}
}

// A fence runs the way the road beside it does: down the rows beside a road
// down the rows, across the columns beside one across them. A fence square
// to its road stood out from it like a row of bars.
func TestAFenceRunsAlongItsRoad(t *testing.T) {
	down := Rect{C0: 10, R0: 0, C1: 10, R1: 36}
	across := Rect{C0: 0, R0: 10, C1: 36, R1: 10}
	for _, c := range []struct {
		road Rect
		want bool
	}{{down, true}, {across, false}} {
		l := &Layout{Seed: 7, Region: 40, Roads: []Rect{c.road}}
		seen := 0
		for r := range 40 {
			for col := range 40 {
				if l.Fence(col, r) {
					seen++
					if got := l.FenceDown(col, r); got != c.want {
						t.Errorf("road %+v: fence at %d,%d runs down: %v", c.road, col, r, got)
					}
				}
			}
		}
		if seen == 0 {
			t.Errorf("road %+v has no fence to look at", c.road)
		}
	}
	if (&Layout{Seed: 7, Region: 40}).FenceDown(1000, 1000) {
		t.Error("a fence that is not there runs down")
	}
}
