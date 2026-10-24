package terrain

// corners are the diagonal bits, and the two sides beside each.
var corners = [4][3]int{
	{NorthEast, North, East},
	{SouthEast, South, East},
	{SouthWest, South, West},
	{NorthWest, North, West},
}

// Canonical is mask with the corners that make no difference cleared: a
// corner counts only when both sides beside it are of the area, as a
// corner of land between two of water is what bends a shore.
func Canonical(mask int) int {
	for _, c := range corners {
		if mask&c[1] == 0 || mask&c[2] == 0 {
			mask &^= c[0]
		}
	}
	return mask
}
