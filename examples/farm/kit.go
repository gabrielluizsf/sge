package farm

import "github.com/gabrielluizsf/sge/inventory"

// kit is what one who works the land has, the player or a villager: their
// money, a bag of seeds, crops and wood, the hotbar slot in hand, the
// crop whose seeds are sown, the best axe they own, how much they have
// sold, and how much of the can's water they have poured.
type kit struct {
	money int
	bag   inventory.Bag
	slot  int
	seed  int // an index of crops: the seeds the seeds slot holds
	axe   int // an index of axes, or noAxe
	sold  int
	used  int // the can's water poured; none is a full can
}
