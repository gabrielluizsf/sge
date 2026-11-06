package farm

import "github.com/gabrielluizsf/sge/backpack"

// chestSlots is how many stacks a chest keeps.
const chestSlots = 10

// chestSlot is a stack of an item kept in a chest: none when n is 0.
type chestSlot struct {
	item string
	n    int
}

// newPiece is a piece of furniture name on cell col, row turned heading,
// a chest empty.
func newPiece(name string, col, row, heading int) piece {
	p := piece{name: name, col: col, row: row, heading: heading}
	if name == "chest" {
		p.held = make([]chestSlot, chestSlots)
	}
	return p
}

// isEmpty reports whether nothing is kept in the piece: always so for
// what is no chest.
func (p piece) isEmpty() bool {
	for _, s := range p.held {
		if s.n > 0 {
			return false
		}
	}
	return true
}

// chestStore is the chest p seen by the inventory screen, beside the bag
// the items move to and from.
type chestStore struct {
	f *farm
	p *piece
}

var _ backpack.Store = chestStore{}

// Slots are the chest's slots.
func (c chestStore) Slots() []backpack.Item {
	out := make([]backpack.Item, len(c.p.held))
	for i, s := range c.p.held {
		if s.n > 0 {
			it := c.f.itemOf(s.item)
			it.Count = s.n
			out[i] = it
		}
	}
	return out
}

// Put moves all of item id carried into the chest: onto the slot that
// keeps it already, or the first that is free.
func (c chestStore) Put(id string) bool {
	n := c.f.bag.Count(id)
	if n == 0 {
		return false
	}
	at := -1
	for i, s := range c.p.held {
		if s.n > 0 && s.item == id {
			at = i
			break
		}
		if s.n == 0 && at < 0 {
			at = i
		}
	}
	if at < 0 || c.p.held[at].n+n > maxCount {
		return false
	}
	c.f.bag.Take(id, n)
	c.p.held[at] = chestSlot{item: id, n: c.p.held[at].n + n}
	return true
}

// Take moves all of item id kept in the chest into the bag.
func (c chestStore) Take(id string) bool {
	for i, s := range c.p.held {
		if s.n > 0 && s.item == id {
			c.f.bag.Add(id, s.n)
			c.p.held[i] = chestSlot{}
			return true
		}
	}
	return false
}
