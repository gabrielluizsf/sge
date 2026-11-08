package human

import "github.com/gabrielluizsf/sge/random"

// faceOf is the eyes, brows and lashes everyone has: women's lashes are
// any of the four, men's the first.
func faceOf(sex Sex, r *random.Source) []string {
	lash := 1
	if sex == Woman {
		lash = r.Between(1, lashes)
	}
	return []string{eyesProxy, eyebrows(r.Between(1, brows)), eyelashes(lash)}
}

// hairOf is the hairstyle of a person: none to a toddler.
func hairOf(sex Sex, years, african float64, r *random.Source) []string {
	list := suited(hairStyles, sex, years, african)
	if len(list) == 0 {
		return nil
	}
	return []string{list[r.Intn(len(list))].item}
}

// eyesOf is the eyes' colour, by chance, as common as it is for a.
func eyesOf(a Ancestry, r *random.Source) string {
	w := make([]float64, len(eyeColours))
	for i, e := range eyeColours {
		w[i] = a.weigh(e.african, e.asian, e.european)
	}
	return eyeColours[choose(r, w)].item
}
