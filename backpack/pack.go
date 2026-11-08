package backpack

// Pack is the screen of what is carried and the hotbar. The game sets
// Source, Hotbar and Active, and each frame hands it the Input and, to
// draw, the Art.
type Pack struct {
	Source Source
	// Hotbar is the ID of the item on each hotbar slot, "" for none.
	Hotbar []string
	// Active is the hotbar slot in hand.
	Active int
	// Open is whether the screen is shown.
	Open bool
	// Cols is how many slots the grid has across, Rows the least it has
	// down.
	Cols, Rows int
	// Title and Hint are the words on the screen.
	Title, Hint string
	// Tabs are the names of the screen's tabs, when it has more than one:
	// the grid of what is carried, then the list of what Crafts can make.
	Tabs []string
	// Tab is the tab shown, an index of Tabs.
	Tab int
	// Crafts makes what the Craft tab lists.
	Crafts Crafter
	// SideTitle is the store's title.
	SideTitle string
	// Side, when set, is a store shown beside the pack's grid, a chest.
	Side Store
	// Scale is how many pixels a pixel of the game's art takes.
	Scale int

	cursor   Spot
	craft    craftState
	made     string // the recipe made last
	held     string // the ID of the item taken up
	from     Spot   // where it was taken from
	dragging bool   // whether the mouse holds it
	at       Point
	rects    rects
}

// New is a closed pack of source, a hotbar of slots slots, and a grid cols
// across.
func New(source Source, slots, cols int) *Pack {
	return &Pack{Source: source, Hotbar: make([]string, slots), Cols: cols, Rows: 3, cursor: Spot{Area: Grid}}
}

// Cursor is the slot the cursor is on.
func (p *Pack) Cursor() Spot { return p.cursor }

// Held is the ID of the item taken up, or "".
func (p *Pack) Held() string { return p.held }

// itemAt is the ID of the item on spot, or "" when the slot is empty.
func (p *Pack) itemAt(s Spot) string {
	switch s.Area {
	case Grid:
		if items := p.Source.Items(); s.Index >= 0 && s.Index < len(items) {
			return items[s.Index].ID
		}
	case Bar:
		if s.Index >= 0 && s.Index < len(p.Hotbar) {
			return p.Hotbar[s.Index]
		}
	case Chest:
		return p.storeAt(s.Index).ID
	}
	return ""
}

// cells is how many slots the grid has: whole rows, enough for what is
// carried.
func (p *Pack) cells() int {
	rows := max(p.Rows, (len(p.Source.Items())+p.Cols-1)/p.Cols)
	return rows * p.Cols
}

// Show opens the screen with the cursor at the first slot.
func (p *Pack) Show() Event {
	p.Open, p.cursor = true, Spot{Area: Grid}
	p.craft = craftState{}
	p.let()
	return Opened
}

// Hide closes the screen, letting go of what is held.
func (p *Pack) Hide() Event {
	p.Open = false
	p.let()
	return Closed
}

// let lets go of the item held.
func (p *Pack) let() { p.held, p.dragging = "", false }
