package daytime

// The day's hours: the sun rises from Dawn to Morning, and sets from Dusk
// to Night.
const (
	Dawn    = 5.0
	Morning = 7.0
	Dusk    = 18.0
	Night   = 20.0
)

// Light is how light the clock's hour is: 1 by day, 0 by night, eased in
// between at dawn and dusk.
func (c *Clock) Light() float64 {
	h := c.Hour
	switch {
	case h >= Morning && h < Dusk:
		return 1
	case h >= Dawn && h < Morning:
		return ease((h - Dawn) / (Morning - Dawn))
	case h >= Dusk && h < Night:
		return 1 - ease((h-Dusk)/(Night-Dusk))
	}
	return 0
}

// ease is a smooth step from 0 to 1.
func ease(t float64) float64 { return t * t * (3 - 2*t) }
