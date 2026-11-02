package human

import "github.com/gabrielluizsf/sge/random"

// ancestryOf is a in shares of 1, or, if it is none, a random one: most
// are of one ancestry, the rest a mix.
func ancestryOf(a Ancestry, r *random.Source) Ancestry {
	if sum := a.African + a.Asian + a.European; sum > 0 {
		return Ancestry{a.African / sum, a.Asian / sum, a.European / sum}
	}
	if r.Float() < 0.6 {
		var out [3]float64
		out[r.Intn(3)] = 1
		return Ancestry{out[0], out[1], out[2]}
	}
	x, y, z := r.Float()+0.05, r.Float()+0.05, r.Float()+0.05
	sum := x + y + z
	return Ancestry{x / sum, y / sum, z / sum}
}

// weigh is how common something is for ancestry a, given how common it
// is among each.
func (a Ancestry) weigh(african, asian, european float64) float64 {
	return a.African*african + a.Asian*asian + a.European*european
}
