package daytime

// Ambient is the light of an hour: how much of each of red, green and
// blue reaches the world, from 0 to 1.
type Ambient [3]float64

// lightKey is the ambient light at an hour.
type lightKey struct {
	hour  float64
	light Ambient
}

// The hours' light, each between its neighbours eased from one to the
// other: moonlit blue, the mauve before dawn, a warm sunrise, full day,
// gold, rose, and blue again.
var (
	moonlit = Ambient{0.36, 0.40, 0.62}
	dayKeys = []lightKey{
		{0, moonlit},
		{4.5, moonlit},
		{5.3, Ambient{0.62, 0.54, 0.66}},
		{6, Ambient{0.96, 0.82, 0.70}},
		{7.5, Ambient{1, 1, 1}},
		{16.5, Ambient{1, 1, 1}},
		{18, Ambient{1, 0.88, 0.70}},
		{19, Ambient{0.82, 0.64, 0.66}},
		{20, Ambient{0.50, 0.50, 0.70}},
		{21.5, moonlit},
		{24, moonlit},
	}
)

// Ambient is the light at the clock's hour.
func (c *Clock) Ambient() Ambient {
	for i := 1; i < len(dayKeys); i++ {
		a, b := dayKeys[i-1], dayKeys[i]
		if c.Hour > b.hour {
			continue
		}
		t := ease((c.Hour - a.hour) / (b.hour - a.hour))
		var out Ambient
		for k := range out {
			out[k] = a.light[k] + (b.light[k]-a.light[k])*t
		}
		return out
	}
	return moonlit
}
