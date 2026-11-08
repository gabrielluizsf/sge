package backpack

import "testing"

// chest is a Store of three slots.
type chest struct{ slots [3]Item }

func (c *chest) Slots() []Item { return c.slots[:] }
func (c *chest) Put(id string) bool {
	for i := range c.slots {
		if c.slots[i].ID == "" {
			c.slots[i] = Item{ID: id, Name: id, Icon: id, Count: 1}
			return true
		}
	}
	return false
}

func (c *chest) Take(id string) bool {
	for i := range c.slots {
		if c.slots[i].ID == id {
			c.slots[i] = Item{}
			return true
		}
	}
	return false
}

func withChest() (*Pack, *chest) {
	c := &chest{}
	p := opened()
	p.Side, p.SideTitle = c, "Chest"
	p.Place(View{W: 1280, H: 720, Scale: 2})
	return p, c
}

func TestAStoreSitsBesideTheGridAndTakesWhatIsSentAcross(t *testing.T) {
	p, c := withChest()
	if p.rects.side.X+p.rects.side.W > p.rects.panel.X {
		t.Fatal("the chest overlaps the pack")
	}
	p.Handle(Input{Dx: 1}) // the second item
	if ev := p.Handle(Input{Select: true, Quick: true}); ev != Stored || c.slots[0].ID != "b_seeds" {
		t.Fatalf("quick send: %v %v", ev, c.slots)
	}
	for range 3 {
		p.Handle(Input{Dx: -1}) // off the grid's left edge into the chest
	}
	if p.Cursor().Area != Chest {
		t.Fatalf("cursor %+v", p.Cursor())
	}
}

func TestHeldItemsGoBothWaysByKeyAndByMouse(t *testing.T) {
	p, c := withChest()
	p.Handle(Input{Select: true}) // take up the first item
	p.Handle(Input{Dx: -1})       // the chest's last slot
	if ev := p.Handle(Input{Select: true}); ev != Stored || c.slots[0].ID != "a_seeds" {
		t.Fatalf("put: %v %v", ev, c.slots)
	}
	// Back by the mouse: drag from the chest's slot onto the grid.
	from, to := centre(p, Spot{Chest, 0}), centre(p, Spot{Grid, 3})
	p.Handle(Input{Press: true, Mouse: from})
	if ev := p.Handle(Input{Release: true, Mouse: to}); ev != Stored || c.slots[0].ID != "" {
		t.Errorf("take: %v %v", ev, c.slots)
	}
}

func TestAFullStoreRefuses(t *testing.T) {
	p, c := withChest()
	for i := range c.slots {
		c.slots[i] = Item{ID: "x", Count: 1}
	}
	if ev := p.Handle(Input{Select: true, Quick: true}); ev != Refused {
		t.Errorf("event %v", ev)
	}
}

// tooled is a pack with a chest keeping a thing to hold and a thing not to,
// and a tool on the hotbar's first slot.
func tooled() (*Pack, *chest, func()) {
	p, c := withChest()
	was := things
	things = append(append([]Item{}, things...), Item{ID: "spade", Name: "spade", Icon: "spade", Hand: true, Tool: true}, Item{ID: "loaf", Name: "loaf", Icon: "loaf", Hand: true}, Item{ID: "log", Name: "log", Icon: "log"})
	p.Hotbar[0] = "spade"
	c.slots[0] = Item{ID: "loaf", Name: "loaf", Icon: "loaf", Count: 1, Hand: true}
	c.slots[1] = Item{ID: "log", Name: "log", Icon: "log", Count: 1}
	return p, c, func() { things = was }
}

// Dragged from the chest onto a slot of the hotbar with no tool on it, an
// item leaves the chest and goes on the slot.
func TestAnItemGoesFromTheChestToTheHotbar(t *testing.T) {
	p, c, done := tooled()
	defer done()
	p.Handle(Input{Press: true, Mouse: centre(p, Spot{Chest, 0})})
	if ev := p.Handle(Input{Release: true, Mouse: centre(p, Spot{Bar, 2})}); ev != Bound || p.Hotbar[2] != "loaf" || c.slots[0].ID != "" {
		t.Errorf("%v, hotbar %v, chest %v", ev, p.Hotbar, c.slots)
	}
}

// A slot that has equipment keeps it, and the item stays in the chest.
func TestAToolsSlotIsNotGivenUp(t *testing.T) {
	p, c, done := tooled()
	defer done()
	p.Handle(Input{Press: true, Mouse: centre(p, Spot{Chest, 0})})
	if ev := p.Handle(Input{Release: true, Mouse: centre(p, Spot{Bar, 0})}); ev != Refused || p.Hotbar[0] != "spade" || c.slots[0].ID != "loaf" {
		t.Errorf("%v, hotbar %v, chest %v", ev, p.Hotbar, c.slots)
	}
}

// A number key does the same with the item under the cursor in the chest.
func TestANumberKeyTakesFromTheChest(t *testing.T) {
	p, c, done := tooled()
	defer done()
	p.cursor = Spot{Chest, 0}
	if ev := p.Handle(Input{Number: 1}); ev != Refused || c.slots[0].ID != "loaf" {
		t.Errorf("onto the tool: %v %v", ev, c.slots)
	}
	if ev := p.Handle(Input{Number: 4}); ev != Bound || p.Hotbar[3] != "loaf" || c.slots[0].ID != "" {
		t.Errorf("onto a free slot: %v, hotbar %v, chest %v", ev, p.Hotbar, c.slots)
	}
}

// What can not be held comes to the bag instead: the hotbar has no place
// for it, and it is not lost.
func TestWhatCanNotBeHeldComesToTheBag(t *testing.T) {
	p, c, done := tooled()
	defer done()
	p.Handle(Input{Press: true, Mouse: centre(p, Spot{Chest, 1})})
	if ev := p.Handle(Input{Release: true, Mouse: centre(p, Spot{Bar, 2})}); ev != Stored || p.Hotbar[2] != "" || c.slots[1].ID != "" {
		t.Errorf("%v, hotbar %v, chest %v", ev, p.Hotbar, c.slots)
	}
}
