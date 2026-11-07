package farm

import (
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/examples/i18n"
)

// What the animals give.
const (
	milk = "milk"
	wool = "wool"
)

// good is something got from the animals, to sell or to make things of:
// its item, its name and what it is for on screen, as English keys, and
// what it sells for.
type good struct {
	item, label, use string
	price            int
}

// goods are all that the animals give.
var goods = []good{
	{item: milk, label: "Milk bottle", use: "From your cow, once a day", price: 18},
	{item: wool, label: "Wool", use: "Beds are made of it", price: 12},
}

// goodOf is the good item is, and whether it is one.
func goodOf(item string) (good, bool) {
	for _, g := range goods {
		if g.item == item {
			return g, true
		}
	}
	return good{}, false
}

// goodsWorth is what k's goods sell for.
func goodsWorth(k *kit) int {
	w := 0
	for _, g := range goods {
		w += k.bag.Count(g.item) * g.price
	}
	return w
}

// sellGoods empties k's bag of the goods, counted as sold.
func sellGoods(k *kit) {
	for _, g := range goods {
		n := k.bag.Count(g.item)
		k.bag.Take(g.item, n)
		k.sold += n
	}
}

// good is the good g as the inventory shows it.
func (c carried) good(g good) backpack.Item {
	return backpack.Item{ID: g.item, Name: i18n.T(g.label), Icon: "icon_" + g.item, Count: c.f.bag.Count(g.item),
		Info: []string{i18n.T(g.use), i18n.T("Sells for $%d", g.price)}}
}
