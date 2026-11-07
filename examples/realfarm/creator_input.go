package main

// cin is what the player did to the creation screen in a frame: arrows,
// Enter, and the mouse (pressed this frame, and held).
type cin struct {
	Dx, Dy      int
	Turn        int    // -1 or 1 while a key turns the person
	Typed       string // the text typed this frame
	Back        bool   // Backspace was pressed
	Cancel      bool   // Escape was pressed
	Pick        bool   // the eyedropper's key was pressed
	Select      bool
	X, Y        int
	Press, Held bool
}

// handle does what in says, and sets the preview to show what it made.
func (c *creator) handle(in cin) {
	if c.finished {
		return
	}
	if c.drop != nil {
		c.pickStep(in)
		c.refresh()
		return
	}
	in = c.typeName(in)
	if in.Pick && c.at().field != nil {
		c.startPick()
	}
	switch {
	case in.Dy != 0:
		c.move(in.Dy)
	case in.Dx != 0 && !c.typing:
		c.change(in.Dx)
	case in.Select:
		c.activate()
	case in.Press:
		c.typing = false
		c.click(in.X, in.Y)
		c.turning, c.turnX = c.drag == nil && c.stageAt.has(in.X, in.Y) && !c.onHit(in.X, in.Y), in.X
	case in.Held && c.drag != nil:
		c.slide(*c.drag, in.X)
	}
	c.turn(in)
	if !in.Held {
		c.drag = nil
	}
	if c.finished {
		return
	}
	c.keep()
	if c.at().kind != textRow {
		c.typing = false
	}
	c.refresh()
}

// click does what a press at x, y is for: it chooses what it is over, and
// steps, presses or slides it.
func (c *creator) click(x, y int) {
	for i := len(c.hits) - 1; i >= 0; i-- {
		h := c.hits[i]
		if !h.at.has(x, y) {
			continue
		}
		switch h.kind {
		case hitCategory:
			c.pick(h.cat)
			if c.at().pressable() {
				c.at().press(c)
			}
		case hitRow:
			c.cat, c.row = h.cat, h.row
			if c.at().pressable() {
				c.at().press(c)
			}
		case hitLess, hitMore:
			c.cat, c.row = h.cat, h.row
			c.change(map[hitKind]int{hitLess: -1, hitMore: 1}[h.kind])
		case hitDropper:
			c.cat, c.row = h.cat, h.row
			c.startPick()
		case hitTrack:
			c.cat, c.row = h.cat, h.row
			d := h
			c.drag = &d
			c.slide(h, x)
		}
		return
	}
}

// slide puts the slider of hit h where x is along its track.
func (c *creator) slide(h hit, x int) {
	f := float64(x-h.at.X) / float64(max(h.at.W, 1))
	c.cats[h.cat].rows[h.row].slider.set(c.st, f)
}
