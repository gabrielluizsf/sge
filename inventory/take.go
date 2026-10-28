package inventory

// Take takes n of item out of the bag, if there are as many, and reports
// whether it did; it takes none otherwise.
func (b *Bag) Take(item string, n int) bool {
	if n < 0 || b.Count(item) < n {
		return false
	}
	if n > 0 {
		b.counts[item] -= n
	}
	return true
}
