package main

import (
	"fmt"

	"github.com/gabrielluizsf/sge/vector"
)

// deliver gives v what they ordered from the player's bag, if it holds it
// all, for the reward; v asks for something else a while later.
func (f *farm) deliver(v *villager) {
	if !v.order.Fill(&f.kit.bag) {
		f.float(v.pos, "FALTAM ITENS", 'r')
		f.sound("deny", nil)
		return
	}
	f.money += v.order.Reward
	f.float(v.pos, fmt.Sprint("OBRIGADO! +$", v.order.Reward), 'y')
	f.sound("sell", nil)
	v.order.Want, v.order.Reward = nil, 0
	v.wait = orderPause * (1 + v.luck.Float())
}

// tradeWood moves a piece of wood from seller to buyer for price, if the
// one has wood and the other the money.
func (f *farm) tradeWood(seller, buyer *kit, price int, at vector.Vector3) {
	switch {
	case !seller.bag.Take(wood, 1):
		f.float(at, "SEM MADEIRA", 'r')
	case buyer.money < price:
		seller.bag.Add(wood, 1)
		f.float(at, "SEM DINHEIRO", 'r')
	default:
		buyer.money -= price
		seller.money += price
		buyer.bag.Add(wood, 1)
		f.float(at, fmt.Sprint("MADEIRA $", price), 'y')
		f.sound("coin", nil)
		return
	}
	f.sound("deny", nil)
}

// doErrand gives the smith what they asked for, for the stone axe.
func (f *farm) doErrand(v *villager) {
	if !f.errand.Fill(&f.kit.bag) {
		f.float(v.pos, "FALTAM ITENS", 'r')
		f.sound("deny", nil)
		return
	}
	f.errandDone = true
	f.axe = max(f.axe, 0)
	f.float(v.pos, "MACHADO DE PEDRA!", 'y')
	f.sound("sell", nil)
	f.talk.choice = 0
}

// buyAxe buys the player axe i from the smith, if it is better than the
// one they have and they have the money.
func (f *farm) buyAxe(i int, at vector.Vector3) {
	switch {
	case f.axe >= i:
		f.float(at, "VOCE JA TEM", 'r')
	case f.money < axes[i].price:
		f.float(at, "SEM DINHEIRO", 'r')
	default:
		f.money -= axes[i].price
		f.axe = i
		f.float(at, "MACHADO DE "+axes[i].label+"!", 'y')
		f.sound("sell", nil)
		return
	}
	f.sound("deny", nil)
}
