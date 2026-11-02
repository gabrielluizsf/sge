package backpack

// Bind puts the item id on hotbar slot slot. An item is on one slot at
// most: if it is on another, the two slots trade. It is Refused for an
// item that can not be held, or a slot there is not.
func (p *Pack) Bind(id string, slot int) Event {
	it, ok := p.Source.Find(id)
	if !ok || !it.Hand || slot < 0 || slot >= len(p.Hotbar) {
		return Refused
	}
	prev := p.Hotbar[slot]
	for j, h := range p.Hotbar {
		if h == id {
			p.Hotbar[j] = prev
		}
	}
	p.Hotbar[slot] = id
	return Bound
}

// Unbind empties hotbar slot slot.
func (p *Pack) Unbind(slot int) Event {
	if slot < 0 || slot >= len(p.Hotbar) || p.Hotbar[slot] == "" {
		return Nothing
	}
	p.Hotbar[slot] = ""
	return Unbound
}

// Slot is the hotbar slot that holds item id, or -1.
func (p *Pack) Slot(id string) int {
	for i, h := range p.Hotbar {
		if h == id {
			return i
		}
	}
	return -1
}

// place does with the item held what letting it go on spot s means: onto
// a hotbar slot it is bound there; onto the grid, a held item that came
// from the hotbar leaves it, and one from the grid goes back.
func (p *Pack) place(s Spot) Event {
	id, from := p.held, p.from
	p.let()
	switch {
	case s.Area == Bar:
		return p.Bind(id, s.Index)
	case from.Area == Bar:
		return p.Unbind(from.Index)
	}
	return Dropped
}
