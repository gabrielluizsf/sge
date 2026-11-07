package farm

import (
	"math"

	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// road is the way on foot from a town to the next: where it turns, its
// length in cells, and how many hours it takes to walk, at a pace a
// villager keeps up with.
type road struct {
	to    *town.Town
	pts   []town.Cell
	cells float64
	hours float64
}

// paceCells is how many cells a villager walks in an hour of the game's
// clock, a fifth less than they can, so they are never left behind.
const paceCells = walkSpeed / tile * (dayLength / 24) * 0.8

// roadFrom is the road from town t to the next, nil when it has none; it
// is worked out once.
func (f *farm) roadFrom(t *town.Town) *road {
	if t == nil {
		return nil
	}
	if r, ok := f.roads[t]; ok {
		return r
	}
	var r *road
	if to, pts := f.towns.Route(t); to != nil {
		r = &road{to: to, pts: pts}
		for i := 1; i < len(pts); i++ {
			r.cells += math.Abs(float64(pts[i].Col-pts[i-1].Col)) + math.Abs(float64(pts[i].Row-pts[i-1].Row))
		}
		r.hours = r.cells / paceCells
	}
	f.roads[t] = r
	return r
}

// point is the world point a fraction p of the way along road r, from its
// first town to its second.
func (f *farm) point(r *road, p float64) vector.Vector3 {
	left := min(max(p, 0), 1) * r.cells
	for i := 1; i < len(r.pts); i++ {
		a, b := r.pts[i-1], r.pts[i]
		d := math.Abs(float64(b.Col-a.Col)) + math.Abs(float64(b.Row-a.Row))
		if left > d && i < len(r.pts)-1 {
			left -= d
			continue
		}
		k := 0.0
		if d > 0 {
			k = min(left/d, 1)
		}
		c := float64(a.Col) + float64(b.Col-a.Col)*k
		w := float64(a.Row) + float64(b.Row-a.Row)*k
		ci, ri := int(math.Floor(c)), int(math.Floor(w))
		at := f.ground.CellCenter(ci, ri)
		return at.Add(vector.Vector3{X: (c - float64(ci)) * tile, Y: (w - float64(ri)) * tile})
	}
	return f.ground.CellCenter(r.pts[0].Col, r.pts[0].Row)
}
