package farm

import "fmt"

// checkPlots makes the plots: the farm's, each as the snapshot says or
// bare, and a garden for each house of the towns woken. A plot of the
// snapshot that is on no such cell is an error.
func (c *checker) checkPlots() error {
	if len(c.s.Plots) > maxPlots {
		return fmt.Errorf("save: %d plots", len(c.s.Plots))
	}
	saved := map[string]plotState{} // by place
	for _, p := range c.s.Plots {
		if err := c.checkPlot(p); err != nil {
			return err
		}
		k := plotSpot(p.Town, p.Col, p.Row)
		if _, twice := saved[k]; twice {
			return fmt.Errorf("save: plot %d, %d twice", p.Col, p.Row)
		}
		saved[k] = p
	}
	for _, p := range c.f.plots {
		if p.town == nil {
			c.take(saved, &plot{col: p.col, row: p.row, crop: -1})
		}
	}
	for _, t := range c.out.towns {
		for _, h := range t.Houses {
			c.take(saved, &plot{col: h.Garden.Col, row: h.Garden.Row, crop: -1, town: t})
		}
	}
	if len(saved) > 0 {
		return fmt.Errorf("save: %d plots on cells that have none", len(saved))
	}
	return nil
}

// plotSpot is where a plot is, as a key: its town's region or none, then
// its cell.
func plotSpot(t *cellState, col, row int) string {
	if t == nil {
		return fmt.Sprint("farm ", col, " ", row)
	}
	return fmt.Sprint("town ", t.Col, " ", t.Row, " ", col, " ", row)
}

// take puts the plot the snapshot holds for p's place in p, adds p to the
// plots made and leaves it out of saved.
func (c *checker) take(saved map[string]plotState, p *plot) {
	k := plotSpot(takeRegion(p.town), p.col, p.row)
	if s, ok := saved[k]; ok {
		p.crop, p.stage, p.days, p.watered, p.dead = s.Crop, s.Stage, s.Days, s.Watered, s.Dead
		delete(saved, k)
	}
	c.out.plots = append(c.out.plots, p)
}

// checkPlot is nil if p is a plot with a crop the game could have grown.
func (c *checker) checkPlot(p plotState) error {
	if err := firstError(
		checkCell("plot", p.Col, p.Row), c.checkRegion(p.Town),
		inRange("crop", p.Crop, -1, len(crops)-1),
		inRange("crop stage", p.Stage, 0, ripe),
		inRange("crop days", p.Days, 0, maxDays),
	); err != nil {
		return err
	}
	if p.Crop < 0 && (p.Stage != 0 || p.Days != 0 || p.Dead) {
		return fmt.Errorf("save: a plot with no crop has one's state")
	}
	return nil
}

// checkRegion is nil if r is nil or a town region woken.
func (c *checker) checkRegion(r *cellState) error {
	if r == nil {
		return nil
	}
	_, err := c.townAt(r)
	return err
}
