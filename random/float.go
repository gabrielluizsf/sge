package random

// Float is the next number of the sequence, from 0 up to 1.
func (s *Source) Float() float64 {
	*s = *s*1664525 + 1013904223
	return float64(*s>>8) / float64(1<<24)
}
