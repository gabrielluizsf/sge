package main

// at is the row chosen.
func (c *creator) at() *row { return &c.cats[c.cat].rows[c.row] }

// shown are the rows of category cat that are shown now.
func (c *creator) shown(cat int) []int {
	var out []int
	for i := range c.cats[cat].rows {
		if c.cats[cat].rows[i].visible(c.st) {
			out = append(out, i)
		}
	}
	return out
}

// move chooses the row d (+1 or -1) places on, across categories, round
// the ends.
func (c *creator) move(d int) {
	type place struct{ cat, row int }
	var all []place
	cur := 0
	for cat := range c.cats {
		for _, r := range c.shown(cat) {
			if cat == c.cat && r == c.row {
				cur = len(all)
			}
			all = append(all, place{cat, r})
		}
	}
	n := all[(cur+d+len(all))%len(all)]
	c.cat, c.row = n.cat, n.row
}

// pick chooses category cat, at its first row shown.
func (c *creator) pick(cat int) {
	c.cat, c.row = cat, c.shown(cat)[0]
}

// keep makes sure the row chosen is one that is shown: a row can vanish
// when a choice (a dress for a top) takes its need away.
func (c *creator) keep() {
	if !c.at().visible(c.st) {
		c.pick(c.cat)
	}
}

// change does what left (-1) or right (+1) is for on the row chosen.
func (c *creator) change(d int) {
	r := c.at()
	switch r.kind {
	case choiceRow:
		r.step(c.st, d)
	case sliderRow:
		r.slider.nudge(c.st, d)
	}
}

// activate does what Enter is for on the row chosen: presses a button,
// steps a choice on.
func (c *creator) activate() {
	switch r := c.at(); r.kind {
	case buttonRow, textRow:
		r.press(c)
	case choiceRow:
		r.step(c.st, 1)
	}
}
