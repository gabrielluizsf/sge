package inventory

// Add puts n of item in the bag; n of 0 or less puts nothing.
func (b *Bag) Add(item string, n int) {
	if n <= 0 {
		return
	}
	if b.counts == nil {
		b.counts = map[string]int{}
	}
	if _, ok := b.counts[item]; !ok {
		b.names = append(b.names, item)
	}
	b.counts[item] += n
}
