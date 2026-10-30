// Package life is a person's years in a game: born on a day of the game's
// calendar, older as its years pass, through childhood to old age, and
// dying of it at the end of a span drawn at birth.
package life

// Calendar is how the game's days make years.
type Calendar struct {
	DaysPerYear int
}

// Life is a person's: the day they were born, which is before the game's
// first for one already grown when it starts, and how many years they
// live, dying of old age then.
type Life struct {
	Born int
	Span float64
}

// Age is how old someone of l is on day, in years.
func (c Calendar) Age(l Life, day int) float64 {
	return float64(day-l.Born) / float64(max(c.DaysPerYear, 1))
}

// Gone is whether someone of l has died of old age by day.
func (c Calendar) Gone(l Life, day int) bool { return c.Age(l, day) >= l.Span }
