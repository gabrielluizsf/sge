package inventory

// TakeAll takes everything l lists out of the bag, if it holds it all,
// and reports whether it did; it takes nothing otherwise.
func (b *Bag) TakeAll(l List) bool {
	if !b.Holds(l) {
		return false
	}
	for _, s := range l {
		b.Take(s.Item, s.N)
	}
	return true
}
