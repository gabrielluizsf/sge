package farm

import (
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/i18n"
)

// What the animals give, and the ovens and the butchers make.
const (
	milk   = "milk"
	wool   = "wool"
	bread  = "bread"
	beef   = "beef"
	mutton = "mutton"
)

// good is something got from the animals or made of what was grown, to sell
// or to make things of: its item, its name and what it is for on screen, as
// English keys, what it sells for and buys for, and how much it feeds, none
// when it is not eaten.
type good struct {
	item, label, use string
	price, cost      int
	nutrition        float64
	// drink says it is drunk, not eaten.
	drink bool
}

// goods are all of them. Bread feeds and sells for more than the wheat it
// is baked from, and meat more still.
var goods = []good{
	{item: milk, label: "Milk bottle", use: "From your cow, once a day", price: 18, nutrition: 18, drink: true},
	{item: wool, label: "Wool", use: "Beds are made of it", price: 12},
	{item: bread, label: "Bread", use: "Baked from wheat at a forge", price: 28, cost: 36, nutrition: 40},
	{item: beef, label: "Beef", use: "From the butcher, or your own cow", price: 34, cost: 45, nutrition: 55},
	{item: mutton, label: "Mutton", use: "From the butcher, or your own sheep", price: 26, cost: 35, nutrition: 45},
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

// good is the good g as the inventory shows it: what is eaten can be taken
// in hand.
func (c carried) good(g good) backpack.Item {
	info := []string{i18n.T(g.use), i18n.T("Sells for $%d", g.price)}
	if g.nutrition > 0 {
		verb := "Eat: +%d food"
		if g.drink {
			verb = "Drink: +%d food"
		}
		info = append(info, i18n.T(verb, int(g.nutrition)))
	}
	return backpack.Item{ID: g.item, Name: i18n.T(g.label), Icon: "icon_" + g.item, Count: c.f.bag.Count(g.item), Hand: g.nutrition > 0, Info: info}
}
