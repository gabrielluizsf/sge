package farm

import "fmt"

// checkOwner is nil if p's owner is none, or a villager of the save with a
// crop of theirs on the farm's soil.
func (c *checker) checkOwner(p plotState) error {
	if err := inRange("plot owner", p.Owner, 0, len(c.s.People)); err != nil {
		return err
	}
	if p.Owner > 0 && (p.Town != nil || p.Crop < 0) {
		return fmt.Errorf("save: an owner of a garden, or of bare soil")
	}
	return nil
}
