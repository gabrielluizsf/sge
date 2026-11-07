package farm

import "github.com/gabrielluizsf/sge/i18n"

// deliver gives v what they ordered from the player's bag, if it holds it
// all, for the reward; v asks for something else a while later.
func (f *farm) deliver(v *villager) {
	if !v.order.Fill(&f.kit.bag) {
		f.speak(v, i18n.T("The order still needs items."))
		f.sound("deny", nil)
		return
	}
	f.money += v.order.Reward
	f.earn(xpOrder)
	f.speak(v, i18n.T("Thank you! Here is $%d.", v.order.Reward))
	f.sound("sell", nil)
	v.order.Want, v.order.Reward = nil, 0
	v.wait = orderPause * (1 + v.luck.Float())
}

// deal is a sale of a piece of wood: who sells it, who buys it, and for
// how much.
type deal struct {
	seller, buyer *kit
	price         int
	with          *villager
}

// tradeWood moves a piece of wood from seller to buyer for the deal's
// price, if the one has wood and the other the money, as the villager the
// deal is with says.
func (f *farm) tradeWood(d deal) {
	seller, buyer, price, v := d.seller, d.buyer, d.price, d.with
	switch {
	case !seller.bag.Take(wood, 1):
		f.speak(v, i18n.T("I don't have wood for that."))
	case buyer.money < price:
		seller.bag.Add(wood, 1)
		f.speak(v, i18n.T("Not enough money for this trade."))
	default:
		buyer.money -= price
		seller.money += price
		buyer.bag.Add(wood, 1)
		f.speak(v, i18n.T("Deal: one wood for $%d.", price))
		f.sound("coin", nil)
		return
	}
	f.sound("deny", nil)
}

// doErrand gives the smith what they asked for, for the stone axe.
func (f *farm) doErrand(v *villager) {
	if !f.errand.Fill(&f.kit.bag) {
		f.speak(v, i18n.T("The errand still needs something."))
		f.sound("deny", nil)
		return
	}
	f.errandDone = true
	f.earn(xpErrand)
	f.axe = max(f.axe, 0)
	f.stock(itemAxe)
	f.speak(v, i18n.T("Here's the stone axe. Happy woodcutting!"))
	f.sound("sell", nil)
	f.talk.choice = 0
}

// buyAxe buys the player axe i from the smith, if it is better than the
// one they have and they have the money.
func (f *farm) buyAxe(i int, v *villager) {
	switch {
	case f.axe >= i:
		f.speak(v, i18n.T("You already have an axe as good as this one."))
	case f.money < axes[i].price:
		f.speak(v, i18n.T("Not enough money for this trade."))
	default:
		f.money -= axes[i].price
		f.axe = i
		f.stock(itemAxe)
		f.speak(v, i18n.T("The %s axe is yours. It cuts well!", i18n.T(axes[i].label)))
		f.sound("sell", nil)
		return
	}
	f.sound("deny", nil)
}
