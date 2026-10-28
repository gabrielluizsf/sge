package clock

// New is a clock reading the time, in seconds, from now. It starts counting
// from the moment it is made.
func New(now func() float64) *Clock {
	return &Clock{now: now, last: now()}
}
