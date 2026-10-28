package inventory

// Count is how many of item are in the bag.
func (b *Bag) Count(item string) int { return b.counts[item] }
