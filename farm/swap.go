package farm

// stonePrice is what a villager pays for a stone.
const stonePrice = 3

// ware is something a villager sells: the item, what it is called on
// screen as an English key, and its price.
type ware struct {
	item, label string
	price       int
}

// tradeWares are what villagers trade among themselves: wood, stone, and
// each crop.
var tradeWares = wares()

func wares() []ware {
	out := []ware{{wood, "Wood", woodPrice}, {stone, "Stone", stonePrice}}
	for c := range crops {
		out = append(out, ware{cropOf(c), crops[c].label, crops[c].sellPrice})
	}
	return out
}

// sale is a ware that changes hands: from one, to the other.
type sale struct {
	ware     ware
	from, to *villager
}

// saleBetween is what a would sell to b, or else what b would sell to a:
// the first ware the seller has and the buyer can pay for.
func saleBetween(a, b *villager) (sale, bool) {
	for _, pair := range [2][2]*villager{{a, b}, {b, a}} {
		for _, w := range tradeWares {
			if pair[0].kit.bag.Count(w.item) > 0 && pair[1].kit.money >= w.price {
				return sale{w, pair[0], pair[1]}, true
			}
		}
	}
	return sale{}, false
}

// settle moves the ware of s from seller to buyer for its price: nothing
// else changes, so the money of the two is the same before and after.
func (s sale) settle() {
	s.from.kit.bag.Take(s.ware.item, 1)
	s.to.kit.bag.Add(s.ware.item, 1)
	s.to.kit.money -= s.ware.price
	s.from.kit.money += s.ware.price
}
