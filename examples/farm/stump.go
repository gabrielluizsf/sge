package farm

import "github.com/gabrielluizsf/sge/tilemap"

// stump is a tree felled, growing back: its cell, the tile of the tree it
// was, and the seconds left.
type stump struct {
	col, row int
	tree     int
	left     float64
}

// regrowIn and rockRegrowIn are how long a felled tree and a mined rock
// take to grow back.
const (
	regrowIn     = 90
	rockRegrowIn = 360
)

// regrowTime is how long what stood as tile id takes to grow back.
func regrowTime(id int) float64 {
	if isRock(id) {
		return rockRegrowIn
	}
	return regrowIn
}

// leftBy is the tile that stays where tile id was felled or mined: a
// tree's stump, nothing where a rock was.
func leftBy(id int) int {
	if isRock(id) {
		return tilemap.Empty
	}
	return tileID["stump"]
}

// regrow grows the stumps back into trees when their time is up, unless
// someone stands on the cell, who would be caught in the trunk: then a
// little later.
func (f *farm) regrow(dt float64) {
	kept := f.stumps[:0]
	for _, s := range f.stumps {
		if s.left -= dt; s.left > 0 || f.someoneOn(s.col, s.row) {
			kept = append(kept, s)
			continue
		}
		if f.thingAt(s.col, s.row) == leftBy(s.tree) {
			f.setThing(s.col, s.row, s.tree)
		}
	}
	f.stumps = kept
}

// someoneOn reports whether anyone about stands on cell col, row.
func (f *farm) someoneOn(col, row int) bool {
	for _, o := range f.bodies() {
		if c, r := f.cellOf(feet(o)); c == col && r == row {
			return true
		}
	}
	return false
}
