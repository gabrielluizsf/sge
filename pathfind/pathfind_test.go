package pathfind

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

// grid is a map of '#' walls in '.' floor, walls all round it.
func grid(rows ...string) Grid {
	return Grid{Limit: 1000, Open: func(c, r int) bool {
		return r >= 0 && r < len(rows) && c >= 0 && c < len(rows[r]) && rows[r][c] != '#'
	}}
}

// The way goes round a wall, the shortest way, and cuts no corner.
func TestAPathGoesRoundAWall(t *testing.T) {
	g := grid(
		".....",
		".###.",
		".#...",
		".#.#.",
		".....",
	)
	way := g.Path(Cell{0, 2}, Cell{2, 2})
	if len(way) == 0 || way[len(way)-1] != (Cell{2, 2}) {
		t.Fatalf("no way: %v", way)
	}
	prev := Cell{0, 2}
	for _, c := range way {
		if !g.Open(c.Col, c.Row) {
			t.Fatalf("walked through a wall at %v", c)
		}
		dc, dr := c.Col-prev.Col, c.Row-prev.Row
		if dc*dc > 1 || dr*dr > 1 || dc != 0 && dr != 0 && !(g.Open(prev.Col+dc, prev.Row) && g.Open(prev.Col, prev.Row+dr)) {
			t.Fatalf("jumped or cut a corner from %v to %v", prev, c)
		}
		prev = c
	}
	// Down the side, across the bottom and up: no diagonal clears the walls.
	if len(way) != 6 {
		t.Errorf("%d steps: %v", len(way), way)
	}
}

// No way is found into a walled-off room, nor past the limit, nor to a
// wall; from a cell to itself the way is empty.
func TestNoPath(t *testing.T) {
	g := grid(
		"..#..",
		"..#..",
		"..#..",
	)
	if way := g.Path(Cell{0, 0}, Cell{4, 0}); way != nil {
		t.Errorf("found %v through a wall", way)
	}
	if way := g.Path(Cell{0, 0}, Cell{2, 0}); way != nil {
		t.Errorf("found %v to a wall", way)
	}
	if way := g.Path(Cell{1, 1}, Cell{1, 1}); way == nil || len(way) != 0 {
		t.Errorf("to itself: %v", way)
	}
	open := Grid{Limit: 10, Open: func(int, int) bool { return true }}
	if way := open.Path(Cell{0, 0}, Cell{100, 3}); way != nil {
		t.Error("searched past the limit")
	}
}

// A route is walked cell by cell, passing each without stopping, to the
// end.
func TestARouteIsSteeredAlong(t *testing.T) {
	r := &Route{Cells: []Cell{{1, 0}, {1, 1}}}
	center := func(c Cell) vector.Vector3 { return vector.Vector3{X: float64(c.Col) * 10, Y: float64(c.Row) * 10} }
	at := vector.Vector3{}
	steps := 0
	for {
		way, done := r.Steer(at, center)
		if done {
			break
		}
		if steps++; steps > 100 {
			t.Fatal("never got there")
		}
		if l := way.Length(); l > 3 {
			way = way.Scale(3 / l)
		}
		at = at.Add(way)
	}
	if at.Sub(vector.Vector3{X: 10, Y: 10}).Length() > arrived {
		t.Errorf("stopped at %v", at)
	}
	if steps != 8 {
		t.Errorf("took %d steps", steps)
	}
}

func TestQueuePopsInRank(t *testing.T) {
	var q queue
	for _, r := range []float64{5, 1, 4, 2, 3} {
		q.push(node{rank: r})
	}
	var got []string
	for len(q) > 0 {
		got = append(got, string(rune('0'+int(q.pop().rank))))
	}
	if s := strings.Join(got, ""); s != "12345" {
		t.Errorf("popped %s", s)
	}
}
