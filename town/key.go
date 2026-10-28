package town

// reach is how far from the origin, in regions, a region can be keyed.
const reach = 1 << 20

// key is region c as one number, which JavaScript maps compare fast; ok
// is false beyond reach.
func key(c Cell) (k int, ok bool) {
	if c.Col <= -reach || c.Col >= reach || c.Row <= -reach || c.Row >= reach {
		return 0, false
	}
	return (c.Col+reach)*(2*reach) + c.Row + reach, true
}

// floorDiv is a divided by b, rounded down, negative a and all.
func floorDiv(a, b int) int {
	q := a / b
	if a%b != 0 && (a < 0) != (b < 0) {
		q--
	}
	return q
}
