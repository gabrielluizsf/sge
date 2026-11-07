package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/town"
)

// checkAway makes the times the towns were last seen: each of a town woken,
// none twice, from the first day's start to now.
func (c *checker) checkAway() error {
	if len(c.s.Away) > maxTowns {
		return fmt.Errorf("save: %d times away", len(c.s.Away))
	}
	now := float64(c.s.Day)*24 + c.s.Hour
	for _, a := range c.s.Away {
		t, err := c.townAt(&cellState{Col: a.Col, Row: a.Row})
		if err != nil {
			return err
		}
		if _, twice := c.out.away[t]; twice {
			return fmt.Errorf("save: a town twice away, at region %d, %d", a.Col, a.Row)
		}
		if err := finite("time away", a.At, 0, now); err != nil {
			return err
		}
		if c.out.away == nil {
			c.out.away = map[*town.Town]float64{}
		}
		c.out.away[t] = a.At
	}
	return nil
}
