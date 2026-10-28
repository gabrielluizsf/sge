package terrain

// Chance is Hash as a number from 0 up to 1.
func Chance(seed uint32, col, row int) float64 {
	return float64(Hash(seed, col, row)) / (1 << 32)
}
