package random

// Intn is the next number of the sequence from 0 up to n, n left out; 0
// when n is not above 0.
func (s *Source) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return min(int(s.Float()*float64(n)), n-1)
}
