package human

import "testing"

func TestOptionsSuitTheirPerson(t *testing.T) {
	for _, sex := range []Sex{Woman, Man} {
		for _, p := range []Part{Hair, Top, Bottom, Dress, Suit, Shoes, Hat} {
			for _, item := range Options(p, sex, 30, 0.5) {
				if got, ok := PartOf(item); !ok || got != p {
					t.Errorf("%s is part %v, listed as %v", item, got, p)
				}
			}
		}
	}
	if len(Options(Dress, Man, 30, 0)) != 0 {
		t.Error("a man is offered dresses")
	}
	if len(Options(Hair, Woman, 30, 0)) < 3 || len(Options(Top, Man, 30, 0)) < 3 {
		t.Error("too few options")
	}
}

func TestEveryOptionIsInTheCatalogue(t *testing.T) {
	in := map[string]bool{}
	for _, c := range Catalogue() {
		in[c] = true
	}
	var all []string
	all = append(all, Eyebrows()...)
	all = append(all, Eyelashes(Woman)...)
	all = append(all, EyeMaterials()...)
	all = append(all, EyesProxy())
	for p := Hair; p <= Hat; p++ {
		all = append(all, Options(p, Woman, 30, 1)...)
		all = append(all, Options(p, Man, 30, 1)...)
	}
	for _, a := range all {
		if !in[a] {
			t.Errorf("%s offered, not in the catalogue", a)
		}
	}
	if _, ok := PartOf("eyes/high-poly/high-poly.mhclo"); ok {
		t.Error("eyes are a part")
	}
	if len(HairColours()) < 5 || len(ClothColours()) < 8 {
		t.Error("too few colours")
	}
}
