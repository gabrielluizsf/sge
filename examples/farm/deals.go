package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/vector"
)

// deliver gives v what they ordered from the player's bag, if it holds it
// all, for the reward; v asks for something else a while later.
func (f *farm) deliver(v *villager) {
	if !v.order.Fill(&f.kit.bag) {
		f.float(v.pos, "Faltam itens", 'r')
		f.sound("deny", nil)
		return
	}
	f.money += v.order.Reward
	f.float(v.pos, fmt.Sprint("Obrigado! +$", v.order.Reward), 'y')
	f.sound("sell", nil)
	v.order.Want, v.order.Reward = nil, 0
	v.wait = orderPause * (1 + v.luck.Float())
}

// deal is a sale of a piece of wood: who sells it, who buys it, and for
// how much.
type deal struct {
	seller, buyer *kit
	price         int
}

// tradeWood moves a piece of wood from seller to buyer for the deal's
// price, if the one has wood and the other the money, said at at.
func (f *farm) tradeWood(d deal, at vector.Vector3) {
	seller, buyer, price := d.seller, d.buyer, d.price
	switch {
	case !seller.bag.Take(wood, 1):
		f.float(at, "Sem madeira", 'r')
	case buyer.money < price:
		seller.bag.Add(wood, 1)
		f.float(at, "Sem dinheiro", 'r')
	default:
		buyer.money -= price
		seller.money += price
		buyer.bag.Add(wood, 1)
		f.float(at, fmt.Sprint("Madeira $", price), 'y')
		f.sound("coin", nil)
		return
	}
	f.sound("deny", nil)
}

// doErrand gives the smith what they asked for, for the stone axe.
func (f *farm) doErrand(v *villager) {
	if !f.errand.Fill(&f.kit.bag) {
		f.float(v.pos, "Faltam itens", 'r')
		f.sound("deny", nil)
		return
	}
	f.errandDone = true
	f.axe = max(f.axe, 0)
	f.float(v.pos, "Machado de pedra!", 'y')
	f.sound("sell", nil)
	f.talk.choice = 0
}

// buyAxe buys the player axe i from the smith, if it is better than the
// one they have and they have the money.
func (f *farm) buyAxe(i int, at vector.Vector3) {
	switch {
	case f.axe >= i:
		f.float(at, "Você já tem", 'r')
	case f.money < axes[i].price:
		f.float(at, "Sem dinheiro", 'r')
	default:
		f.money -= axes[i].price
		f.axe = i
		f.float(at, "Machado de "+axes[i].label+"!", 'y')
		f.sound("sell", nil)
		return
	}
	f.sound("deny", nil)
}
