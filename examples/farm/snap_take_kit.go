package farm

import (
	"github.com/gabrielluizsf/sge/inventory"
	"github.com/gabrielluizsf/sge/request"
)

// takeKit is k as plain data.
func takeKit(k *kit) kitState {
	return kitState{Money: k.money, Bag: takeStacks(k.bag.Content()), Slot: k.slot, Seed: k.seed, Axe: k.axe, Sold: k.sold, Used: k.used, Exact: k.exact}
}

// takeStacks are l as plain data, never nil.
func takeStacks(l inventory.List) []stackState {
	out := make([]stackState, 0, len(l))
	for _, s := range l {
		out = append(out, stackState{Item: s.Item, N: s.N})
	}
	return out
}

// takeOrder is o as plain data.
func takeOrder(o request.Order) orderState {
	return orderState{Want: takeStacks(o.Want), Reward: o.Reward}
}
