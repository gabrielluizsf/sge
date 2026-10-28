package daytime

// Slot is something to do from hour From up to hour To, round midnight
// when To comes before From.
type Slot struct {
	From, To float64
	Do       string
}

// Schedule is a day's slots, the first that holds deciding.
type Schedule []Slot

// At is what the schedule says to do at the clock's hour, or "" for
// nothing.
func (s Schedule) At(c *Clock) string {
	for _, slot := range s {
		if c.Between(slot.From, slot.To) {
			return slot.Do
		}
	}
	return ""
}
