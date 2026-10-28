package main

// stump is a tree felled, growing back: its cell, and the seconds left.
type stump struct {
	col, row int
	left     float64
}

// regrowIn is how long a felled tree takes to grow back.
const regrowIn = 90

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
		if f.thingAt(s.col, s.row) == tileID["stump"] {
			f.setThing(s.col, s.row, tileID["tree0"])
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
