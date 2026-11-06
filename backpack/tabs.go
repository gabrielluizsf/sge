package backpack

// tabMouse is what the mouse did to the tabs: pressing on one shows it.
func (p *Pack) tabMouse(in Input) Event {
	if !in.Press {
		return Nothing
	}
	for i, r := range p.rects.tabs {
		if r.In(in.Mouse) {
			if i == p.Tab {
				return Nothing
			}
			p.SetTab(i)
			return Tabbed
		}
	}
	return Nothing
}
