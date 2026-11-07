package farm

import "github.com/gabrielluizsf/sge/pathfind"

// area is the cells within cols and rows of a center.
type area struct {
	center     pathfind.Cell
	cols, rows int
}

// quarry is a job of striking what stands on the cells of an area, from a
// cell beside it: what counts, the tool, how many strokes it takes, and
// the cells that may be touched. Only cells in sight count.
type quarry struct {
	area
	is    func(thing int) bool
	slot  int
	uses  int
	where func(col, row int) bool
}

// strike is a job striking one of the quarry's things, picked at random,
// or nil when there is none.
func (m *mind) strike(f *farm, q quarry) *job {
	near := f.sitesOf(q, true)
	if len(near) == 0 {
		return nil
	}
	return near[m.v.luck.Intn(len(near))]
}

// sitesOf are the jobs of striking what q is after, in reading order: each
// the cell to stand on beside it, and what it is. With seen, only those in
// sight; without, all, as where they stand is wanted for a time they are
// not watched.
func (f *farm) sitesOf(q quarry, seen bool) []*job {
	var near []*job
	a := q.area
	for r := a.center.Row - a.rows; r <= a.center.Row+a.rows; r++ {
		for col := a.center.Col - a.cols; col <= a.center.Col+a.cols; col++ {
			if !q.is(f.thingAt(col, r)) || !q.where(col, r) {
				continue
			}
			for _, b := range banks {
				if sc, sr := col-b.dc, r-b.dr; f.open(sc, sr) && q.where(sc, sr) && (!seen || f.inSight(f.center(pathfind.Cell{Col: sc, Row: sr}))) {
					near = append(near, &job{to: pathfind.Cell{Col: sc, Row: sr}, face: b.face, slot: q.slot, uses: q.uses, rest: 1, target: pathfind.Cell{Col: col, Row: r}})
					break
				}
			}
		}
	}
	return near
}

// chopIn is a job felling one of the trees in a, picked at random, from a
// cell beside it; only trees and cells for which ok is true count.
func (m *mind) chopIn(f *farm, a area, ok func(col, row int) bool) *job {
	return m.strike(f, quarry{area: a, is: isTree, slot: slotAxe, uses: axes[m.v.kit.axe].hits, where: ok})
}
