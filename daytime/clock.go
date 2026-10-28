// Package daytime is a world's time of day: a Clock whose hours go round
// as the game runs, the light they give, from day to night through dusk
// and dawn, a Shade darkening the screen as night falls, and Schedules of
// what to do at what hour.
package daytime

import "github.com/gabrielluizsf/sge/game"

// Clock is the time of day: Hour from 0 up to 24, going round, Day
// counting the days gone round, and Speed, how many of its hours pass in
// a second of the game. Added to an object, it keeps time.
type Clock struct {
	Hour  float64
	Day   int
	Speed float64
}

// Update moves the clock on dt seconds of the game.
func (c *Clock) Update(_ *game.Object, dt float64) {
	c.Hour += dt * c.Speed
	for c.Hour >= 24 {
		c.Hour -= 24
		c.Day++
	}
}

// Between reports whether the clock is from hour from up to hour to, to
// left out, round midnight when to comes before from: 22 to 6 is the
// night.
func (c *Clock) Between(from, to float64) bool {
	if from <= to {
		return c.Hour >= from && c.Hour < to
	}
	return c.Hour >= from || c.Hour < to
}
