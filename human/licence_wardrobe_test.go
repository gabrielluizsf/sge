package human

import "testing"

// Everything the wardrobe offers is under CC0, in its proxy and in its
// mesh, which may say something else than its proxy: the game refuses a
// person made of any other.
func TestEveryGarmentOfTheWardrobeIsCC0(t *testing.T) {
	k := realKit(t)
	b := baselines[0]
	for _, g := range wardrobe() {
		lics, err := k.Licences(b.body, Look{Eyes: b.look.Eyes, Wear: []string{g.item}})
		if err != nil {
			t.Errorf("%s: %v", g.item, err)
			continue
		}
		for name, lic := range lics {
			if lic != "CC0" && !isSkinFile(name) {
				t.Errorf("%s: %s is under %q", g.item, name, lic)
			}
		}
	}
}

func isSkinFile(name string) bool { return hasFolder(name, "skins") || hasFolder(name, "eyes") }
