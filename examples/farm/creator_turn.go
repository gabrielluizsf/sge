package main

// keyTurn is how many degrees a second a key turns the person.
const keyTurn = 120

// onHit reports whether x, y is over something the mouse acts on.
func (c *creator) onHit(x, y int) bool {
	for _, h := range c.hits {
		if h.at.has(x, y) {
			return true
		}
	}
	return false
}

// turn turns the person by what the player did this frame, and by that
// alone: the mouse dragged on the stage sideways, or the keys held.
func (c *creator) turn(in cin) {
	if c.turning && in.Held {
		c.view.Turn(float64(in.X-c.turnX) * turnPerPixel)
		c.turnX = in.X
		c.stage.Dirty()
	}
	if !in.Held {
		c.turning = false
	}
	if in.Turn != 0 {
		c.view.Turn(float64(in.Turn) * keyTurn * c.dt)
		c.stage.Dirty()
	}
}
