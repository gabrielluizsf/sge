package clock

// Tick is the seconds since the clock was made or last ticked. A clock source
// that steps backwards counts as no time passing.
func (c *Clock) Tick() float64 {
	t := c.now()
	dt := t - c.last
	c.last = t
	return max(dt, 0)
}
