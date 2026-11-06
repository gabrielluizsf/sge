package backpack

// Store is a second place to keep items, such as a chest, shown as a grid
// beside the pack's.
type Store interface {
	// Slots are the store's slots, always as many; a slot with no item
	// has an empty ID, and Count below 0 is not shown.
	Slots() []Item
	// Put moves the whole stack of item id from what is carried into the
	// store, and reports whether it did: not when there is no slot free.
	Put(id string) bool
	// Take moves the whole stack of item id from the store to what is
	// carried, and reports whether it did.
	Take(id string) bool
}

// SideCols is how many slots a row of the store has.
const sideCols = 5

// storeAt is the item on store slot i, or an empty one.
func (p *Pack) storeAt(i int) Item {
	if p.Side == nil {
		return Item{}
	}
	if s := p.Side.Slots(); i >= 0 && i < len(s) {
		return s[i]
	}
	return Item{}
}

// send moves the stack held, or under the cursor, across: from the grid
// into the store, from the store to the grid.
func (p *Pack) send(s Spot) Event {
	id := p.itemAt(s)
	if p.Side == nil || id == "" {
		return Nothing
	}
	switch s.Area {
	case Grid:
		if p.Side.Put(id) {
			return Stored
		}
	case Chest:
		if p.Side.Take(id) {
			return Stored
		}
	}
	return Refused
}

// stepSide moves the cursor on the store, dx slots across and dy down; off
// its right edge it goes to the grid, and off its bottom to the hotbar.
func (p *Pack) stepSide(dx, dy int) {
	c := p.cursor
	col, row := c.Index%sideCols, c.Index/sideCols
	rows := (len(p.Side.Slots()) + sideCols - 1) / sideCols
	switch {
	case dx > 0 && col == sideCols-1:
		p.cursor = Spot{Grid, min(row, p.cells()/p.Cols-1) * p.Cols}
		return
	case dy > 0 && row == rows-1:
		p.cursor = Spot{Bar, min(c.Index%sideCols, len(p.Hotbar)-1)}
		return
	}
	col = max(0, min(sideCols-1, col+dx))
	row = max(0, row+dy)
	p.cursor.Index = min(len(p.Side.Slots())-1, row*sideCols+col)
}
