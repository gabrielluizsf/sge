package terrain

// reach is how far from the origin, in cells or chunks, cellKey can key.
const reach = 1 << 20

// cellKey is col, row as one number, a key JavaScript maps compare as
// fast as Go's, where an array key is made into a string every time; ok
// is false beyond reach, where the number would lose precision.
func cellKey(col, row int) (key int, ok bool) {
	if col <= -reach || col >= reach || row <= -reach || row >= reach {
		return 0, false
	}
	return (col+reach)*(2*reach) + row + reach, true
}
