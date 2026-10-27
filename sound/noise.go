package sound

// Noise is n samples of white noise, the same for the same seed.
func Noise(seed uint32, n int) []float32 {
	s := make([]float32, n)
	x := seed | 1
	for i := range s {
		x ^= x << 13
		x ^= x >> 17
		x ^= x << 5
		s[i] = float32(x)/float32(1<<31) - 1
	}
	return s
}
