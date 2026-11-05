package farm

import "fmt"

// deliver gives v what they ordered from the player's bag, if it holds it
// all, for the reward; v asks for something else a while later.
func (f *farm) deliver(v *villager) {
	if !v.order.Fill(&f.kit.bag) {
		f.speak(v, "Ainda faltam itens para o pedido.")
		f.sound("deny", nil)
		return
	}
	f.money += v.order.Reward
	f.speak(v, fmt.Sprint("Obrigado! Aqui estão $", v.order.Reward, "."))
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
		f.speak(v, "Você não tem madeira para isso.")
	case buyer.money < price:
		seller.bag.Add(wood, 1)
		f.speak(v, "Falta dinheiro para essa troca.")
	default:
		buyer.money -= price
		seller.money += price
		buyer.bag.Add(wood, 1)
		f.speak(v, fmt.Sprint("Negócio fechado: uma madeira por $", price, "."))
		f.sound("coin", nil)
		return
	}
	f.sound("deny", nil)
}

// doErrand gives the smith what they asked for, for the stone axe.
func (f *farm) doErrand(v *villager) {
	if !f.errand.Fill(&f.kit.bag) {
		f.speak(v, "Ainda falta algo da encomenda.")
		f.sound("deny", nil)
		return
	}
	f.errandDone = true
	f.axe = max(f.axe, 0)
	f.stock(itemAxe)
	f.speak(v, "Aqui está o machado de pedra. Boa colheita de lenha!")
	f.sound("sell", nil)
	f.talk.choice = 0
}

// buyAxe buys the player axe i from the smith, if it is better than the
// one they have and they have the money.
func (f *farm) buyAxe(i int, v *villager) {
	switch {
	case f.axe >= i:
		f.speak(v, "Você já tem um machado tão bom quanto esse.")
	case f.money < axes[i].price:
		f.speak(v, "Falta dinheiro para essa troca.")
	default:
		f.money -= axes[i].price
		f.axe = i
		f.stock(itemAxe)
		f.speak(v, "Machado de "+axes[i].label+" é seu. Ele corta bem!")
		f.sound("sell", nil)
		return
	}
	f.sound("deny", nil)
}
