// Package clock measures the time between frames.
package clock

// Clock hands out the seconds elapsed since it was last asked.
type Clock struct {
	now  func() float64
	last float64
}
