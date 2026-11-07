package kit

import (
	"testing"

	"github.com/gabrielluizsf/sge/save"
)

// A character saved before a garment was retired still loads, wearing what
// took its place where it was.
func TestASavedCharacterWithARetiredGarmentStillLoads(t *testing.T) {
	old := map[string]bool{
		"clothes/culturalibre_male_boots/culturalibre_male_boots.mhclo":           true,
		"clothes/toigo_skirt_with_lace_ruffle/toigo_skirt_with_lace_ruffle.mhclo": true,
		"hair/braid01/braid01.mhclo":                                              true,
		"hair/bob01/bob01.mhclo":                                                  true,
		"hair/bob02/bob02.mhclo":                                                  true,
	}
	for item := range old {
		c := someone(11)
		c.Body.Gender = 1
		c.Look.Wear = append(append([]string(nil), c.Look.Wear...), item)
		if c.Check() == nil {
			t.Fatalf("%s is still in the catalogue", item)
		}
		s := save.In(t.TempDir(), characterVersion)
		if err := SaveCharacter(s, c); err != nil {
			t.Fatal(err)
		}
		got, ok := LoadCharacter(s)
		if !ok {
			t.Fatalf("a character in %s no longer loads", item)
		}
		for _, w := range got.Look.Wear {
			if old[w] {
				t.Errorf("%s was kept", w)
			}
		}
		if len(got.Look.Wear) != len(c.Look.Wear) || got.Look.Wear[len(c.Look.Wear)-1] == item {
			t.Errorf("the garment was not replaced in its place: %v", got.Look.Wear)
		}
	}
}
