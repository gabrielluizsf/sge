package human

import (
	"strings"
	"testing"
)

// womanly and manly are what a name says of who wears the garment.
var (
	womanly = []string{"female", "flats", "ballet", "stiletto", "mj_cloth", "heroine", "skirt", "dress", "camisole", "halter", "bodice", "keyhole", "tube_top", "cloche", "kimono", "harem", "tank_top", "bob", "ponytail", "braid", "long01"}
	manly   = []string{"male_", "_male", "polo", "fedora", "hero_boots"}
)

// isFeminine and isMasculine are whether the name of item says it is a
// woman's or a man's.
func isFeminine(item string) bool {
	return !isMasculine(item) && containsAny(item, womanly)
}

func isMasculine(item string) bool {
	return !strings.Contains(item, "female") && !strings.Contains(item, "heroine") && containsAny(item, manly)
}

func containsAny(s string, subs []string) bool {
	for _, x := range subs {
		if strings.Contains(s, x) {
			return true
		}
	}
	return false
}

// No garment made for women is offered to men or given to them by
// chance, nor the reverse: the wardrobe's flags say what its names say.
func TestTheWardrobeSaysWhoWearsWhat(t *testing.T) {
	for _, g := range wardrobe() {
		if isFeminine(g.item) && g.men {
			t.Errorf("%s is a woman's and offered to men", g.item)
		}
		if isMasculine(g.item) && g.women {
			t.Errorf("%s is a man's and offered to women", g.item)
		}
		if !g.women && !g.men {
			t.Errorf("%s is nobody's", g.item)
		}
	}
	for _, s := range []Sex{Man, Woman} {
		for _, p := range []Part{Hair, Top, Bottom, Dress, Suit, Shoes, Hat} {
			for _, item := range Options(p, s, 30, 1) {
				if s == Man && isFeminine(item) || s == Woman && isMasculine(item) {
					t.Errorf("%v is offered %s", s, item)
				}
			}
		}
	}
}

// A man given clothes by chance wears nothing a woman's, whatever the seed.
func TestRandomPeopleWearWhatSuitsTheirSex(t *testing.T) {
	for seed := uint32(1); seed <= 400; seed++ {
		for _, s := range []Sex{Man, Woman} {
			_, l := Random(Spec{Sex: s, Seed: seed})
			for _, w := range l.Wear {
				g, ok := worn(w)
				if !ok {
					continue
				}
				if s == Man && !g.men || s == Woman && !g.women {
					t.Fatalf("seed %d: a %v wears %s", seed, s, w)
				}
			}
		}
	}
}
