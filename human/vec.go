package human

import (
	"math"
	"sort"
)

func add(a, b [3]float64) [3]float64 { return [3]float64{a[0] + b[0], a[1] + b[1], a[2] + b[2]} }
func sub(a, b [3]float64) [3]float64 { return [3]float64{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }
func scale(a [3]float64, s float64) [3]float64 {
	return [3]float64{a[0] * s, a[1] * s, a[2] * s}
}
func cross(a, b [3]float64) [3]float64 {
	return [3]float64{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}
func abs(v float64) float64 { return math.Abs(v) }

func normalize(a [3]float64) [3]float64 {
	l := math.Sqrt(a[0]*a[0] + a[1]*a[1] + a[2]*a[2])
	if l == 0 {
		return [3]float64{0, 1, 0}
	}
	return scale(a, 1/l)
}

// sortedKeys are m's keys in order, so what is made of a map is the same
// each time.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
