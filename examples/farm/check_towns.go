package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/town"
)

// checkTowns makes the towns of the snapshot, the regions it names, each
// with a town in it, none twice.
func (c *checker) checkTowns() error {
	if len(c.s.Towns) > maxTowns {
		return fmt.Errorf("save: %d towns", len(c.s.Towns))
	}
	seen := map[town.Cell]bool{}
	for _, r := range c.s.Towns {
		if err := checkCell("town", r.Col, r.Row); err != nil {
			return err
		}
		t := c.f.towns.Of(r.Col, r.Row)
		if t == nil || seen[t.Region] {
			return fmt.Errorf("save: no town, or one twice, at region %d, %d", r.Col, r.Row)
		}
		seen[t.Region] = true
		c.out.towns = append(c.out.towns, t)
	}
	return nil
}

// townAt is the town of the snapshot in region r, or an error.
func (c *checker) townAt(r *cellState) (*town.Town, error) {
	for _, t := range c.out.towns {
		if t.Region.Col == r.Col && t.Region.Row == r.Row {
			return t, nil
		}
	}
	return nil, fmt.Errorf("save: a town at region %d, %d that was not woken", r.Col, r.Row)
}
