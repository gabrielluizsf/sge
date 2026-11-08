package backpack

// Need is one material of a recipe: its picture and name, how many the
// recipe takes and how many are carried.
type Need struct {
	Name, Icon string
	Want, Have int
}

// Recipe is something that can be made, as the Craft tab lists it: a
// picture, a name, the materials it takes and a few lines of what it is.
type Recipe struct {
	ID, Name, Icon string
	Needs          []Need
	Info           []string
}

// Can reports whether every material is carried in the amount wanted.
func (r Recipe) Can() bool {
	for _, n := range r.Needs {
		if n.Have < n.Want {
			return false
		}
	}
	return true
}

// Crafter is what the Craft tab makes things with.
type Crafter interface {
	// Recipes are the things that can be made, in the order the tab lists
	// them, with the materials as carried now.
	Recipes() []Recipe
	// Make makes recipe id, taking its materials, and reports whether it
	// did: it does nothing when a material is lacking.
	Make(id string) bool
}

// craftState is the Craft tab's cursor and the first row it shows.
type craftState struct{ cursor, top int }

// recipes are the recipes now.
func (p *Pack) recipes() []Recipe {
	if p.Crafts == nil {
		return nil
	}
	return p.Crafts.Recipes()
}

// Crafting reports whether the Craft tab is the one shown.
func (p *Pack) Crafting() bool {
	return p.Tab == 1 && p.Crafts != nil && len(p.Tabs) > 1 && p.Side == nil
}

// CraftCursor is the row of the Craft tab the cursor is on.
func (p *Pack) CraftCursor() int { return p.craft.cursor }

// tab shows the tab d steps on from the one shown, round the tabs.
func (p *Pack) tab(d int) Event {
	n := len(p.Tabs)
	if n < 2 || p.Side != nil {
		return Nothing
	}
	p.Tab = ((p.Tab+d)%n + n) % n
	p.let()
	return Tabbed
}

// SetTab shows tab i.
func (p *Pack) SetTab(i int) {
	if i >= 0 && i < len(p.Tabs) {
		p.Tab = i
		p.let()
	}
}

// craftStep moves the cursor dy rows, keeping it in view.
func (p *Pack) craftStep(dy int) {
	n := len(p.recipes())
	c := &p.craft
	c.cursor = max(0, min(n-1, c.cursor+dy))
	if c.cursor < c.top {
		c.top = c.cursor
	}
	if c.cursor >= c.top+p.Rows {
		c.top = c.cursor - p.Rows + 1
	}
}

// craftKeys is what the keyboard did on the Craft tab: up and down move
// the cursor, Enter makes the recipe under it, Escape closes the screen.
func (p *Pack) craftKeys(in Input) Event {
	switch {
	case in.Select:
		return p.make(p.craft.cursor)
	case in.Cancel:
		return p.Hide()
	case in.Dy != 0:
		p.craftStep(in.Dy)
		return Moved
	}
	return Nothing
}

// make makes the recipe on row i: Crafted, or Refused when a material is
// lacking.
func (p *Pack) make(i int) Event {
	rs := p.recipes()
	if i < 0 || i >= len(rs) || !p.Crafts.Make(rs[i].ID) {
		return Refused
	}
	p.made = rs[i].ID
	return Crafted
}

// Made is the ID of the recipe made last.
func (p *Pack) Made() string { return p.made }

// craftMouse is what the mouse did on the Craft tab: moving over a row
// puts the cursor there, pressing on it makes the recipe.
func (p *Pack) craftMouse(in Input) Event {
	for k, r := range p.rects.rows {
		if !r.In(in.Mouse) {
			continue
		}
		i := p.craft.top + k
		switch {
		case in.Press:
			p.craft.cursor = i
			return p.make(i)
		case i != p.craft.cursor && i < len(p.recipes()):
			p.craft.cursor = i
			return Moved
		}
	}
	return Nothing
}
