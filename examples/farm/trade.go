package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/vector"
)

// buy buys k a packet of crop c's seeds, if the money is there, and tells
// so above at.
func (f *farm) buy(k *kit, c int, at vector.Vector3) bool {
	price := crops[c].seedPrice
	if k.money < price {
		f.float(at, "Sem dinheiro", 'r')
		f.sound("deny", nil)
		return false
	}
	k.money -= price
	k.bag.Add(seedsOf(c), 1)
	f.float(at, fmt.Sprint("-$", price), 'r')
	f.sound("coin", &at)
	return true
}

// sell sells k's whole harvest and wood, and tells what it made above at.
func (f *farm) sell(k *kit, at vector.Vector3) bool {
	earned := worth(k)
	if earned == 0 {
		f.float(at, "Nada a vender", 'r')
		f.sound("deny", nil)
		return false
	}
	for c := range crops {
		n := k.bag.Count(cropOf(c))
		k.bag.Take(cropOf(c), n)
		k.sold += n
	}
	k.sold += k.bag.Count(wood)
	k.bag.Take(wood, k.bag.Count(wood))
	k.money += earned
	f.float(at, fmt.Sprint("+$", earned), 'y')
	f.sound("sell", &at)
	return true
}
