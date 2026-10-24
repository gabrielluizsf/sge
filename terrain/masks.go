package terrain

// Masks are the canonical masks, the 47 tiles a set that joins up in
// every way needs.
func Masks() []int {
	var out []int
	for mask := range AllSides + 1 {
		if Canonical(mask) == mask {
			out = append(out, mask)
		}
	}
	return out
}
