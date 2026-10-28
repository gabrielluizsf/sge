package terrain

// Fractal is octaves of Noise, each twice as fine and half as strong as
// the last, from 0 to 1: shapes with large bends and small wiggles, as
// coasts and woods have.
func Fractal(seed uint32, x, y float64, octaves int) float64 {
	sum, amp, total := 0.0, 1.0, 0.0
	for i := range octaves {
		sum += amp * Noise(seed+uint32(i)*0x9e3779b9, x, y)
		total += amp
		x, y, amp = x*2, y*2, amp/2
	}
	return sum / total
}
