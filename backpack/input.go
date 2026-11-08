package backpack

// Input is what the player did this frame: arrows (or WASD) as steps, Enter
// or space (Quick, with Shift held: send the item across to the chest or
// back, or make the recipe), Escape, Tab (-1 or 1: the tab before or after),
// a number key (1 for the first hotbar slot, 0 for none),
// and the mouse: where it is, and whether its button went down or up.
type Input struct {
	Dx, Dy         int
	Select, Cancel bool
	Quick          bool
	Tab            int
	Number         int
	Mouse          Point
	Press, Release bool
}

// Handle does what in says to the open pack, and is what it did. A closed
// pack does nothing.
func (p *Pack) Handle(in Input) Event {
	if !p.Open {
		return Nothing
	}
	if in.Tab != 0 {
		return p.tab(in.Tab)
	}
	if ev := p.mouse(in); ev != Nothing {
		return ev
	}
	return p.keys(in)
}
