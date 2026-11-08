package human

// garment is a proxy of the packs, and for whom it is: women, men, and the
// years of age it suits; and, for hair, the least African ancestry it
// suits, 0 for any.
type garment struct {
	item       string
	women, men bool
	from, to   float64
	african    float64
}

// suits is whether the garment suits a person of sex and age years.
func (g garment) suits(sex Sex, years float64) bool {
	ok := g.women
	if sex == Man {
		ok = g.men
	}
	return ok && years >= g.from && years <= g.to
}

// forBoth is a garment of any sex, from to years of age.
func forBoth(item string, from, to float64) garment {
	return garment{item: item, women: true, men: true, from: from, to: to}
}

// forWomen is a garment of women, from to years of age.
func forWomen(item string, from, to float64) garment {
	return garment{item: item, women: true, from: from, to: to}
}

// forMen is a garment of men, from to years of age.
func forMen(item string, from, to float64) garment {
	return garment{item: item, men: true, from: from, to: to}
}

// cloth is the proxy of the clothes folder named name.
func cloth(name string) string { return "clothes/" + name + "/" + name + ".mhclo" }

// suited are the garments of list that suit a person of sex and years,
// and ancestry of african.
func suited(list []garment, sex Sex, years, african float64) []garment {
	var out []garment
	for _, g := range list {
		if g.suits(sex, years) && african >= g.african {
			out = append(out, g)
		}
	}
	return out
}

// wardrobe is every garment there is: hair, clothes, shoes and hats.
func wardrobe() []garment {
	var all []garment
	for _, l := range [][]garment{hairStyles, tops, bottoms, dresses, suits, shoes, hats} {
		all = append(all, l...)
	}
	return all
}

// worn is the garment of the proxy item, false for any that is not one:
// eyes, brows and lashes are everyone's.
func worn(item string) (garment, bool) {
	for _, g := range wardrobe() {
		if g.item == item {
			return g, true
		}
	}
	return garment{}, false
}
