package main

import "github.com/gabrielluizsf/sge/inventory"

// kit is what one who works the land has, the player or a villager: their
// money, a bag of seeds, crops and wood, the hotbar slot in hand, the best
// axe they own, and how much they have sold.
type kit struct {
	money int
	bag   inventory.Bag
	slot  int
	axe   int // an index of axes, or noAxe
	sold  int
}
