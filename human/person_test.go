package human

import (
	"reflect"
	"strings"
	"testing"
)

// worn is the kinds of what a look wears, by folder and name.
func kindOf(name string) string {
	switch {
	case strings.HasPrefix(name, "hair/"):
		return "hair"
	case strings.HasPrefix(name, "clothes/"):
		return "clothes"
	}
	return "other"
}

func TestRandomIsTheSamePersonForTheSameSeed(t *testing.T) {
	for seed := uint32(0); seed < 20; seed++ {
		b1, l1 := Random(Spec{Seed: seed})
		b2, l2 := Random(Spec{Seed: seed})
		if b1 != b2 || !reflect.DeepEqual(l1, l2) {
			t.Fatalf("seed %d made two people", seed)
		}
	}
}

func TestRandomSeedsMakeDifferentPeople(t *testing.T) {
	hairs, outfits, eyes, hairColours, skins := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[Colour]bool{}, map[float64]bool{}
	for seed := uint32(1); seed <= 50; seed++ {
		_, l := Random(Spec{Seed: seed})
		eyes[l.Eyes] = true
		hairColours[l.HairColour] = true
		skins[l.SkinTone] = true
		var outfit []string
		for _, w := range l.Wear {
			switch kindOf(w) {
			case "hair":
				hairs[w] = true
			case "clothes":
				outfit = append(outfit, w)
			}
		}
		outfits[strings.Join(outfit, "+")] = true
	}
	if len(hairs) < 6 || len(outfits) < 20 || len(eyes) < 4 || len(hairColours) < 10 || len(skins) < 20 {
		t.Errorf("50 seeds: %d hair styles, %d outfits, %d eye colours, %d hair colours, %d skin tones", len(hairs), len(outfits), len(eyes), len(hairColours), len(skins))
	}
}

func TestRandomRespectsSexAndAge(t *testing.T) {
	for seed := uint32(1); seed <= 50; seed++ {
		for _, c := range []struct {
			sex   Sex
			years float64
		}{{Woman, 30}, {Man, 30}, {Woman, 8}, {Man, 70}, {Woman, 85}, {Man, 14}} {
			b, l := Random(Spec{Sex: c.sex, Age: c.years, Seed: seed})
			if (c.sex == Woman) != (b.Gender < 0.5) {
				t.Fatalf("seed %d: %v has gender %v", seed, c.sex, b.Gender)
			}
			if y := b.Years(); y < c.years-0.5 || y > c.years+0.5 {
				t.Fatalf("seed %d: asked %v years, got %v", seed, c.years, y)
			}
			for _, w := range l.Wear {
				if g, ok := worn(w); ok && !g.suits(c.sex, c.years) {
					t.Errorf("seed %d: %v of %v years wears %s", seed, c.sex, c.years, w)
				}
			}
		}
	}
}

func TestRandomAdultsAreAdultsAndBodiesAreNatural(t *testing.T) {
	for seed := uint32(1); seed <= 100; seed++ {
		b, _ := Random(Spec{Seed: seed})
		if y := b.Years(); y < 18 || y > 80 {
			t.Fatalf("seed %d: a random adult of %v years", seed, y)
		}
		for _, v := range []float64{b.Muscle, b.Weight, b.Height, b.Proportions, b.Gender} {
			if v < 0 || v > 1 {
				t.Fatalf("seed %d: body %+v out of range", seed, b)
			}
		}
		if sum := b.African + b.Asian + b.European; sum < 0.99 || sum > 1.01 {
			t.Fatalf("seed %d: ancestry sums to %v", seed, sum)
		}
	}
}

func TestRandomFollowsAncestry(t *testing.T) {
	mean := func(a Ancestry) (tone, dark float64) {
		for seed := uint32(1); seed <= 100; seed++ {
			b, l := Random(Spec{Ancestry: a, Seed: seed})
			if a.African > 0 && b.African < 0.99 || a.European > 0 && b.European < 0.99 || a.Asian > 0 && b.Asian < 0.99 {
				t.Fatalf("asked %+v, got %+v", a, b)
			}
			tone += l.SkinTone / 100
			if strings.Contains(l.Eyes, "brown") {
				dark += 0.01
			}
		}
		return
	}
	af, afDark := mean(Ancestry{African: 1})
	eu, euDark := mean(Ancestry{European: 1})
	as, asDark := mean(Ancestry{Asian: 1})
	if !(af > as && as > eu) {
		t.Errorf("skin tones african %.2f, asian %.2f, european %.2f", af, as, eu)
	}
	if afDark < 0.8 || asDark < 0.8 || euDark > 0.7 || euDark < 0.15 {
		t.Errorf("brown eyes: african %.2f asian %.2f european %.2f", afDark, asDark, euDark)
	}
	lum := func(a Ancestry) (l float64) {
		for seed := uint32(1); seed <= 100; seed++ {
			_, look := Random(Spec{Ancestry: a, Age: 25, Seed: seed})
			l += (look.HairColour.R + look.HairColour.G + look.HairColour.B) / 300
		}
		return
	}
	if lum(Ancestry{African: 1}) > 0.2 || lum(Ancestry{Asian: 1}) > 0.2 || lum(Ancestry{European: 1}) < 0.2 {
		t.Errorf("hair: african %.2f asian %.2f european %.2f", lum(Ancestry{African: 1}), lum(Ancestry{Asian: 1}), lum(Ancestry{European: 1}))
	}
}

func TestRandomMixesAncestryWhenNoneIsAsked(t *testing.T) {
	var mixed int
	for seed := uint32(1); seed <= 50; seed++ {
		b, _ := Random(Spec{Seed: seed})
		if max(b.African, b.Asian, b.European) < 0.95 {
			mixed++
		}
	}
	if mixed < 10 || mixed > 45 {
		t.Errorf("%d of 50 are of mixed ancestry", mixed)
	}
}

func TestRandomDressesWomenAndMenAlike(t *testing.T) {
	dresses, suits := 0, 0
	for seed := uint32(1); seed <= 80; seed++ {
		_, w := Random(Spec{Sex: Woman, Age: 30, Seed: seed})
		_, m := Random(Spec{Sex: Man, Age: 30, Seed: seed})
		for _, it := range m.Wear {
			if strings.Contains(it, "dress") || strings.Contains(it, "skirt") {
				t.Fatalf("a man in %s", it)
			}
		}
		for _, it := range w.Wear {
			if strings.Contains(it, "dress") {
				dresses++
			}
			if strings.Contains(it, "suit") {
				suits++
			}
		}
	}
	if dresses < 5 || suits < 3 {
		t.Errorf("women in %d dresses and %d suits of 80", dresses, suits)
	}
}
