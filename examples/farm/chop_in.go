package farm

import "github.com/gabrielluizsf/sge/pathfind"

// area is the cells within cols and rows of a center.
type area struct {
	center     pathfind.Cell
	cols, rows int
}

// chopIn is a job felling one of the trees in a, picked at random, from a
// cell beside it; only trees and cells for which ok is true count.
func (m *mind) chopIn(f *farm, a area, ok func(col, row int) bool) *job {
	v := m.v
	var near []*job
	for r := a.center.Row - a.rows; r <= a.center.Row+a.rows; r++ {
		for col := a.center.Col - a.cols; col <= a.center.Col+a.cols; col++ {
			if !isTree(f.thingAt(col, r)) || !ok(col, r) {
				continue
			}
			for _, b := range banks {
				if sc, sr := col-b.dc, r-b.dr; f.open(sc, sr) && ok(sc, sr) {
					near = append(near, &job{to: pathfind.Cell{Col: sc, Row: sr}, face: b.face, slot: slotAxe, uses: axes[v.kit.axe].hits, rest: 1})
					break
				}
			}
		}
	}
	if len(near) == 0 {
		return nil
	}
	return near[v.luck.Intn(len(near))]
}
