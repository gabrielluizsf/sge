package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/town"
)

// The ground behind a town's northern houses, for its shops: wide enough
// for the pen and all of them side by side, deep enough for the deepest.
const (
	backCols  = 27
	backRows  = penRows
	shedPitch = 6
)

// sheds are the shops that stand there, west to east after the pen, by the
// names their pictures have.
var sheds = []string{"stable", "bakery", "butcher"}

// townPen is the pen of town t's stable, at the west of its back ground.
func townPen(t *town.Town) rect {
	if t.Back == (town.Rect{}) {
		return rect{}
	}
	return rect{t.Back.C0, t.Back.R1 - penRows + 1, penCols, penRows}
}

// shedRect is the cells of the i-th of town t's sheds, a's footprint for
// it, its front to the row free to walk along; none when t keeps no ground.
func shedRect(a Art, t *town.Town, i int) rect {
	p := townPen(t)
	if p.cols == 0 {
		return rect{}
	}
	cols, rows := a.Footprint(sheds[i])
	return rect{p.col + p.cols + 1 + i*shedPitch, t.Back.R1 - rows + 1, cols, rows}
}

// stableRect is the cells of town t's stable.
func stableRect(a Art, t *town.Town) rect { return shedRect(a, t, 0) }

// onShed reports whether cell col, row is under one of town t's sheds.
func onShed(a Art, t *town.Town, col, row int) bool {
	for i := range sheds {
		if shedRect(a, t, i).has(col, row) {
			return true
		}
	}
	return false
}

// shedBy is the town and the name of the shed o stands in front of.
func (f *farm) shedBy(o *game.Object) (*town.Town, string) {
	col, row := f.cellOf(feet(o))
	for _, t := range f.towns.Near(col, row, f.towns.Reach()+backRows+8) {
		for i, name := range sheds {
			r := shedRect(f.art, t, i)
			if row == r.row+r.rows && col >= r.col-1 && col <= r.col+r.cols {
				return t, name
			}
		}
	}
	return nil, ""
}

// openShed opens the menu of the shed name of town t.
func (f *farm) openShed(t *town.Town, name string) {
	f.shop.open, f.shop.kind, f.shop.choice, f.shop.at = true, name, 0, town.Cell{Col: t.Back.C0, Row: t.Back.R1}
	f.sound("ui_open", nil)
}

// shedOf is the index in sheds of the shop a trade keeps, or -1.
func shedOf(role string) int {
	switch role {
	case roleDealer:
		return 0
	case roleBaker:
		return 1
	case roleButcher:
		return 2
	}
	return -1
}

// shopPost is where one who keeps a shop stands at a moment of their day,
// by key: at its front or to a side of it, and a dealer in the pen too.
func shopPost(f *farm, v *villager, key uint32) *job {
	i := shedOf(v.role)
	p := townPen(v.town)
	if i < 0 || p.cols == 0 {
		return nil
	}
	r := shedRect(f.art, v.town, i)
	front := r.row + r.rows
	spots := []pathfind.Cell{{Col: r.col + r.cols/2, Row: front}, {Col: r.col, Row: front}, {Col: r.col + r.cols - 1, Row: front}}
	if v.role == roleDealer {
		spots = append(spots, pathfind.Cell{Col: p.col + p.cols/2, Row: p.row + p.rows}, pathfind.Cell{Col: p.col + p.cols/2, Row: p.row + p.rows - 2},
			pathfind.Cell{Col: p.col + 1, Row: p.row + 1}, pathfind.Cell{Col: p.col + p.cols - 2, Row: p.row + 2})
	}
	return &job{to: spots[int(key%uint32(len(spots)))], face: "down", slot: -1, rest: 2 + float64(key%5)}
}

// tendJob is the work in sight of one who keeps a shop: a baker first sees
// to their firewood, a butcher now and then buys an animal from the dealer,
// and otherwise each stands to their shop.
func (m *mind) tendJob(f *farm) *job {
	v := m.v
	if v.town == nil {
		return nil
	}
	switch {
	case v.role == roleBaker:
		if j := m.bakerWood(f); j != nil {
			return j
		}
	case v.role == roleButcher && v.luck.Float() < 0.1:
		if j := m.butcherBuys(f); j != nil {
			return j
		}
	}
	return shopPost(f, v, uint32(v.luck.Intn(1<<16)))
}
