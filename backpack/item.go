package backpack

// Item is something carried, as the screen shows it: a picture, a name,
// how many and a few lines of what it is for.
type Item struct {
	// ID tells it from other items.
	ID string
	// Name is what the tooltip calls it, Icon the Art's picture of it.
	Name, Icon string
	// Count is how many are carried, shown on its slot; below 0 it is not
	// shown.
	Count int
	// Info are the lines under the name in the tooltip: its price, what it
	// does.
	Info []string
	// Hand says it may be put on the hotbar, to be held.
	Hand bool
	// Tool says it is equipment: a slot of the hotbar that has it is not
	// given up to what comes out of a store.
	Tool bool
}

// Source is what the game carries.
type Source interface {
	// Items are the things carried, in the order the grid shows them.
	Items() []Item
	// Find is the item with that ID, carried or not: what is on the hotbar
	// may have run out.
	Find(id string) (Item, bool)
}
