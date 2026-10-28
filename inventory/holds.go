package inventory

// Holds reports whether the bag has everything l lists, a stack of an item
// listed twice counting twice.
func (b *Bag) Holds(l List) bool {
	need := map[string]int{}
	for _, s := range l {
		need[s.Item] += s.N
		if b.Count(s.Item) < need[s.Item] {
			return false
		}
	}
	return true
}
