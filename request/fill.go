package request

import "github.com/gabrielluizsf/sge/inventory"

// Fill takes what the order wants out of bag, if it holds it all, and
// reports whether it did; it takes nothing otherwise.
func (o Order) Fill(bag *inventory.Bag) bool {
	return len(o.Want) > 0 && bag.TakeAll(o.Want)
}
