package request

import (
	"math"

	"github.com/gabrielluizsf/sge/inventory"
	"github.com/gabrielluizsf/sge/random"
)

// Draw is an order drawn from the book by src: one to Kinds different
// items, at least one of each and at most its offer's Most, paying their
// price times the Bonus, rounded up. A book with no offers draws an empty
// order.
func (b Book) Draw(src *random.Source) Order {
	var o Order
	if len(b.Offers) == 0 {
		return o
	}
	left := append([]Offer(nil), b.Offers...)
	kinds := src.Between(1, max(1, min(b.Kinds, len(left))))
	worth := 0
	for range kinds {
		i := src.Intn(len(left))
		f := left[i]
		left = append(left[:i], left[i+1:]...)
		n := src.Between(1, max(1, f.Most))
		o.Want = append(o.Want, inventory.Stack{Item: f.Item, N: n})
		worth += n * f.Price
	}
	o.Reward = int(math.Ceil(float64(worth) * max(b.Bonus, 1)))
	return o
}
