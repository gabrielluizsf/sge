package random

// Between is the next number of the sequence from lo to hi, both in.
func (s *Source) Between(lo, hi int) int {
	return lo + s.Intn(hi-lo+1)
}
