package inventory

// Content is what is in the bag, a stack an item, in the order the items
// were first put in; items there are none of left out.
func (b *Bag) Content() List {
	var l List
	for _, name := range b.names {
		if n := b.counts[name]; n > 0 {
			l = append(l, Stack{Item: name, N: n})
		}
	}
	return l
}
