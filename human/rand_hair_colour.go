package human

import "github.com/gabrielluizsf/sge/random"

// hairColourOf is a hair colour as common for a, each strand's shade a
// little off the shade's own.
func hairColourOf(a Ancestry, r *random.Source) Colour {
	w := make([]float64, len(hairShades))
	for i, s := range hairShades {
		w[i] = a.weigh(s.african, s.asian, s.european)
	}
	c := hairShades[choose(r, w)].colour
	k := 0.9 + 0.2*r.Float()
	return RGB(min(c.R*k, 1), min(c.G*k, 1), min(c.B*k, 1))
}
