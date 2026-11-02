package human

import "github.com/gabrielluizsf/sge/random"

// skinToneOf is a skin tone, 0 lightest to 1 darkest, for ancestry a,
// each person a little off the middle of their ancestry's.
func skinToneOf(a Ancestry, r *random.Source) float64 {
	t := a.African*0.82 + a.Asian*0.3 + a.European*0.14
	return min(max(t+(r.Float()+r.Float()-1)*0.2, 0.02), 1)
}
