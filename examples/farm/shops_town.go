package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/town"
)

// The ground behind a town's northern houses, for its shops: wide enough
// for all of them side by side, deep enough for the deepest.
const (
	backCols = 24
	backRows = penRows
)

// townPen is the pen of town t's stable, at the west of its back ground.
func townPen(t *town.Town) rect {
	if t.Back == (town.Rect{}) {
		return rect{}
	}
	return rect{t.Back.C0, t.Back.R1 - penRows + 1, penCols, penRows}
}

// stableRect is the cells of town t's stable, a's footprint for it, beside
// its pen, its front to the row free to walk along.
func stableRect(a Art, t *town.Town) rect {
	p := townPen(t)
	if p.cols == 0 {
		return rect{}
	}
	cols, rows := a.Footprint("stable")
	return rect{p.col + p.cols + 1, t.Back.R1 - rows + 1, cols, rows}
}

// stableBy is the town whose stable o stands in front of, or nil.
func (f *farm) stableBy(o *game.Object) *town.Town {
	col, row := f.cellOf(feet(o))
	for _, t := range f.towns.Near(col, row, f.towns.Reach()+backRows+8) {
		r := stableRect(f.art, t)
		if row == r.row+r.rows && col >= r.col-1 && col <= r.col+r.cols {
			return t
		}
	}
	return nil
}

// openStable opens the stable's menu of animals for sale.
func (f *farm) openStable(t *town.Town) {
	f.shop.open, f.shop.animals, f.shop.choice, f.shop.at = true, true, 0, town.Cell{Col: t.Back.C0, Row: t.Back.R1}
	f.sound("ui_open", nil)
}

// dealerPost is where a livestock dealer stands at a moment of their day:
// at the stable's front, or at the gate of the pen, by key.
func dealerPost(f *farm, t *town.Town, key uint32) *job {
	r, p := stableRect(f.art, t), townPen(t)
	if p.cols == 0 {
		return nil
	}
	spots := []pathfind.Cell{
		{Col: r.col + r.cols/2, Row: r.row + r.rows},
		{Col: p.col + p.cols/2, Row: p.row + p.rows},
		{Col: p.col + p.cols/2, Row: p.row + p.rows - 2},
		{Col: p.col + 1, Row: p.row + 1},
		{Col: p.col + p.cols - 2, Row: p.row + 2},
	}
	return &job{to: spots[int(key%uint32(len(spots)))], face: "up", slot: -1, rest: 2 + float64(key%5)}
}

// tendJob is a livestock dealer's work in sight: to one of their posts.
func (m *mind) tendJob(f *farm) *job {
	if m.v.town == nil {
		return nil
	}
	return dealerPost(f, m.v.town, uint32(m.v.luck.Intn(1<<16)))
}
