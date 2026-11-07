package human

// retired are the proxies the wardrobe no longer offers, each with the one
// that takes its place: a look saved before still has the old name.
var retired = map[string]string{
	cloth("culturalibre_male_boots"):      cloth("toigo_ankle_boots_male"),
	cloth("toigo_skirt_with_lace_ruffle"): cloth("toigo_tiered_skirt"),
	hair("braid01"):                       hair("ponytail01"),
	hair("bob01"):                         hair("toigo_blunt_bob"),
	hair("bob02"):                         hair("toigo_curled_under_bob"),
}

// Mended is l with every retired proxy it wears replaced by its successor,
// in the same place, so what an older game saved still loads. A look that
// wears none is returned as it is.
func Mended(l Look) Look {
	var wear []string
	for i, w := range l.Wear {
		to, ok := retired[w]
		if !ok {
			continue
		}
		if wear == nil {
			wear = append([]string(nil), l.Wear...)
		}
		wear[i] = to
	}
	if wear == nil {
		return l
	}
	l.Wear = wear
	return l
}
