package human

import "github.com/gabrielluizsf/sge/random"

// pickOne is the item of one of list, by chance, "" if it is empty.
func pickOne(list []garment, r *random.Source) string {
	if len(list) == 0 {
		return ""
	}
	return list[r.Intn(len(list))].item
}

// outfitOf is what a person of sex and years wears: a dress, a suit or a
// top and bottom, shoes, and now and then a hat.
func outfitOf(sex Sex, years float64, r *random.Source) []string {
	kind := choose(r, []float64{0.35, 0.25, 0.4})
	if sex == Man {
		kind = []int{1, 2, 2, 1, 2}[r.Intn(5)]
	}
	at := func(list []garment) []garment { return suited(list, sex, years, 0) }
	var out []string
	if d := pickOne(at(dresses), r); kind == 0 && d != "" {
		out = []string{d}
	} else if s := pickOne(at(suits), r); kind == 1 && s != "" {
		out = []string{s}
	} else {
		out = []string{pickOne(at(tops), r), pickOne(at(bottoms), r)}
	}
	out = append(out, pickOne(at(shoes), r))
	if r.Float() < 0.15 {
		if h := pickOne(at(hats), r); h != "" {
			out = append(out, h)
		}
	}
	return out
}
