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
