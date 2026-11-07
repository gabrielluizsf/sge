package farm

// trading is a villager at the market from hour from to hour to.
type trading struct {
	v        *villager
	from, to float64
}

// swapEvery is how long the two of them must be at the market together for
// each ware that changes hands between them, and swapMost the most that do.
const (
	swapEvery = 0.5
	swapMost  = 3
)

// settleSwaps settles the sales between villagers who traded at the same
// time: each pair, in the order they are in, for each stretch they have
// together, what one has that the other will pay for, until they have
// nothing of the kind or it is the most. What the buyer pays is the seller's:
// no money is made or lost between them.
func (f *farm) settleSwaps(traders []trading) {
	for i, a := range traders {
		for _, b := range traders[i+1:] {
			together := min(a.to, b.to) - max(a.from, b.from)
			for n := 0; n < min(int(together/swapEvery), swapMost); n++ {
				s, ok := saleBetween(a.v, b.v)
				if !ok {
					break
				}
				s.settle()
			}
		}
	}
}
