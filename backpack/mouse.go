package backpack

// mouse is what the mouse did. Pressing on an item takes it up and drags
// it; letting go over a hotbar slot binds it there, and over the grid or
// outside it, an item taken from the hotbar leaves it. Pressing and letting
// go on the same slot is a click: the item stays up, to be put down by the
// next click. Moving over a slot puts the cursor there.
func (p *Pack) mouse(in Input) Event {
	p.at = in.Mouse
	spot, over := p.spotAt(in.Mouse)
	switch {
	case in.Press && p.held != "" && !p.dragging:
		return p.put(spot, over)
	case in.Press && over && p.itemAt(spot) != "":
		p.cursor, p.held, p.from, p.dragging = spot, p.itemAt(spot), spot, true
		return Picked
	case in.Release && p.dragging:
		p.dragging = false
		if over && spot == p.from {
			return Nothing // a click: still up
		}
		return p.put(spot, over)
	case over && !p.dragging && p.held == "" && spot != p.cursor:
		p.cursor = spot
		return Moved
	}
	return Nothing
}

// put lets the item held go on spot, if the mouse is over one, or
// outside the screen's slots.
func (p *Pack) put(s Spot, over bool) Event {
	if over {
		p.cursor = s
	} else {
		s = Spot{}
	}
	return p.place(s)
}
