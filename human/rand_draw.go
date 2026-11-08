package human

import "github.com/gabrielluizsf/sge/random"

// newDraw is the random sequence of seed, warmed up: the first numbers
// of this generator of nearby seeds are nearly the same.
func newDraw(seed uint32) *random.Source {
	r := random.New(seed)
	for range 8 {
		r.Float()
	}
	return r
}

// spreadOf is centre, give either way, more likely near the centre.
func spreadOf(r *random.Source, centre, give float64) float64 {
	return min(max(centre+(r.Float()+r.Float()-1)*give, 0), 1)
}

// choose is the index of one of weights, by chance, as heavy as it is.
func choose(r *random.Source, weights []float64) int {
	var sum float64
	for _, w := range weights {
		sum += w
	}
	at := r.Float() * sum
	for i, w := range weights {
		if at < w {
			return i
		}
		at -= w
	}
	return len(weights) - 1
}
