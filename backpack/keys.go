package backpack

// keys is what the keyboard did: a number key puts the item held, or
// under the cursor, on that slot of the hotbar; Enter takes up the item
// under the cursor and, with one held, puts it down there; Escape lets go
// of what is held, or closes the screen; the arrows move the cursor.
func (p *Pack) keys(in Input) Event {
	if p.Crafting() {
		return p.craftKeys(in)
	}
	switch {
	case in.Select && in.Quick && p.held == "":
		return p.send(p.cursor)
	case in.Number > 0:
		return p.number(in.Number - 1)
	case in.Select:
		return p.enter()
	case in.Cancel && p.held != "":
		p.let()
		return Dropped
	case in.Cancel:
		return p.Hide()
	case in.Dx != 0 || in.Dy != 0:
		p.step(in.Dx, in.Dy)
		return Moved
	}
	return Nothing
}

// number binds the item held, or under the cursor, to slot.
func (p *Pack) number(slot int) Event {
	id := p.held
	if id == "" {
		id = p.itemAt(p.cursor)
	}
	p.let()
	if id == "" {
		return Nothing
	}
	return p.Bind(id, slot)
}

// enter takes up the item under the cursor, or puts down the one held.
func (p *Pack) enter() Event {
	if p.held != "" {
		return p.place(p.cursor)
	}
	id := p.itemAt(p.cursor)
	if id == "" {
		return Nothing
	}
	p.held, p.from = id, p.cursor
	return Picked
}

// step moves the cursor dx slots across and dy down: the grid's bottom row
// leads down to the hotbar, and the hotbar up to the grid.
func (p *Pack) step(dx, dy int) {
	c := p.cursor
	if c.Area == Chest && p.Side != nil {
		p.stepSide(dx, dy)
		return
	}
	if c.Area == Grid && p.Side != nil && dx < 0 && c.Index%p.Cols == 0 {
		rows := (len(p.Side.Slots()) + sideCols - 1) / sideCols
		p.cursor = Spot{Chest, min(c.Index/p.Cols, rows-1)*sideCols + sideCols - 1}
		return
	}
	if c.Area == Bar {
		switch {
		case dy < 0 && p.Side != nil && c.Index >= p.Cols:
			rows := (len(p.Side.Slots()) + sideCols - 1) / sideCols
			p.cursor = Spot{Chest, (rows-1)*sideCols + min(c.Index-p.Cols, sideCols-1)}
		case dy < 0:
			p.cursor = Spot{Grid, p.cells() - p.Cols + min(c.Index, p.Cols-1)}
		default:
			p.cursor.Index = max(0, min(len(p.Hotbar)-1, c.Index+dx))
		}
		return
	}
	col, row := c.Index%p.Cols, c.Index/p.Cols
	col = max(0, min(p.Cols-1, col+dx))
	row += dy
	if row >= p.cells()/p.Cols {
		p.cursor = Spot{Bar, min(col, len(p.Hotbar)-1)}
		return
	}
	p.cursor.Index = max(0, row)*p.Cols + col
}
