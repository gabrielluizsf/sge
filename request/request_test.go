package request

import (
	"testing"

	"github.com/gabrielluizsf/sge/inventory"
	"github.com/gabrielluizsf/sge/random"
)

var book = Book{
	Offers: []Offer{{"wood", 5, 4}, {"carrot", 3, 12}, {"tomato", 2, 25}},
	Kinds:  2,
	Bonus:  1.5,
}

// Orders ask for one or two different items, some of each, no more than
// offered, for half again their price.
func TestOrdersAreDrawnFromTheBook(t *testing.T) {
	src := random.New(3)
	sizes := map[int]bool{}
	for range 200 {
		o := book.Draw(src)
		sizes[len(o.Want)] = true
		worth := 0
		seen := map[string]bool{}
		for _, s := range o.Want {
			var f Offer
			for _, f = range book.Offers {
				if f.Item == s.Item {
					break
				}
			}
			if s.N < 1 || s.N > f.Most || seen[s.Item] {
				t.Fatalf("order %+v", o)
			}
			seen[s.Item] = true
			worth += s.N * f.Price
		}
		if o.Reward != (worth*3+1)/2 {
			t.Fatalf("order %+v pays %d for %d worth", o, o.Reward, worth)
		}
	}
	if !sizes[1] || !sizes[2] || len(sizes) != 2 {
		t.Errorf("orders of %v kinds", sizes)
	}
	if o := (Book{}).Draw(src); len(o.Want) != 0 || o.Reward != 0 {
		t.Errorf("an empty book drew %+v", o)
	}
}

// An order is filled from a bag only when the bag holds all it wants.
func TestAnOrderIsFilledWhole(t *testing.T) {
	o := Order{Want: inventory.List{{Item: "wood", N: 3}, {Item: "carrot", N: 1}}, Reward: 30}
	var b inventory.Bag
	b.Add("wood", 3)
	if o.Fill(&b) || b.Count("wood") != 3 {
		t.Fatal("filled without the carrot")
	}
	b.Add("carrot", 2)
	if !o.Fill(&b) || b.Count("wood") != 0 || b.Count("carrot") != 1 {
		t.Fatalf("not filled: %v", b.Content())
	}
	if (Order{}).Fill(&b) {
		t.Error("filled an order for nothing")
	}
}
