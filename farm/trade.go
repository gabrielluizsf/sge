package farm

import (
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/vector"
)

// buy buys k a packet of crop c's seeds, if the money is there, and tells
// so above at.
func (f *farm) buy(k *kit, c int, at vector.Vector3) bool {
	price := crops[c].seedPrice
	if k.money < price {
		f.float(at, i18n.T("Not enough money"), 'r')
		f.sound("deny", nil)
		return false
	}
	k.money -= price
	k.bag.Add(seedsOf(c), 1)
	f.float(at, i18n.T("-$%d", price), 'r')
	f.sound("coin", &at)
	return true
}

// sell sells k's whole harvest and wood, and tells what it made above at.
func (f *farm) sell(k *kit, at vector.Vector3) bool {
	earned := worth(k)
	if earned == 0 {
		f.float(at, i18n.T("Nothing to sell"), 'r')
		f.sound("deny", nil)
		return false
	}
	k.sellOut()
	f.earnBy(k, xpSale(earned))
	f.float(at, i18n.T("+$%d", earned), 'y')
	f.sound("sell", &at)
	return true
}

// sellOut empties k's bag of the harvest and wood, for what it is worth,
// which goes to k's money; it gives that.
func (k *kit) sellOut() int {
	earned := worth(k)
	for c := range crops {
		n := k.bag.Count(cropOf(c))
		k.bag.Take(cropOf(c), n)
		k.sold += n
	}
	sellGoods(k)
	k.sold += k.bag.Count(wood)
	k.bag.Take(wood, k.bag.Count(wood))
	k.money += earned
	return earned
}
