package farm

import "github.com/gabrielluizsf/sge/vector"

// bed is a bed of a room: the cells it covers, whether it is the couple's
// double one, and who lies in it: used once anyone does, and each side of a
// double one taken or not.
type bed struct {
	span
	double bool
	used   bool
	taken  [2]bool
}

// center is the middle of the bed in the world.
func (b bed) center(f *farm) vector.Vector3 {
	c := f.ground.CellCenter(b.c0, b.r0)
	return c.Add(vector.Vector3{X: float64(b.cols-1) * tile / 2, Y: float64(b.rows-1) * tile / 2})
}

// anim is the way someone lies on side of the bed: asleep in a bed of
// their own, or on one side of a double bed.
func (b bed) anim(side int) string {
	if !b.double {
		return AnimSleep
	}
	return [2]string{AnimSleepLeft, AnimSleepRight}[side]
}

// bedsOf are the beds ordered with the double one first, whichever way
// they are laid.
func bedsOf(beds []bed) []bed {
	out := make([]bed, 0, len(beds))
	for _, double := range []bool{true, false} {
		for _, b := range beds {
			if b.double == double {
				out = append(out, b)
			}
		}
	}
	return out
}

// berth is a place to lie: a bed, and the side of it, 0 or 1, which is
// always 0 in a single bed.
type berth struct{ bed, side int }

// bedFor is where someone sleeps. A couple's half of the double bed, the
// side of their order first; a single sleeper's own bed, or a free side of
// the double one when theirs is taken. It reports false when none is free.
func (r *room) bedFor(couple bool, order int) (berth, bool) {
	sides := [2]int{order & 1, 1 - order&1}
	if couple {
		if b, ok := r.freeSide(true, sides); ok {
			return b, true
		}
		return r.freeSide(false, sides)
	}
	if b, ok := r.freeSide(false, sides); ok {
		return b, true
	}
	return r.freeSide(true, sides)
}

// freeSide is a free place in a double bed or in a single one, the sides
// tried in the order given.
func (r *room) freeSide(double bool, sides [2]int) (berth, bool) {
	for i, b := range r.beds {
		if b.double != double {
			continue
		}
		for _, s := range sidesOf(double, sides) {
			if !b.taken[s] {
				return berth{i, s}, true
			}
		}
	}
	return berth{}, false
}

// sidesOf are the sides of a bed to try: those given for a double bed,
// the one of a single bed.
func sidesOf(double bool, sides [2]int) []int {
	if double {
		return sides[:]
	}
	return []int{0}
}
