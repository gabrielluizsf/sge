package human

import "math"

// irises is where the irises are, in the pack's eye textures: two discs,
// the rest being the whites of the eyes and what is around them.
func irises(x, y float64) float64 {
	const radius, edge = 0.118, 0.012
	m := 0.0
	for _, c := range [][2]float64{{0.721, 0.302}, {0.298, 0.720}} {
		d := math.Hypot(x-c[0], y-c[1])
		m = max(m, min(max((radius+edge-d)/(2*edge), 0), 1))
	}
	return m
}
