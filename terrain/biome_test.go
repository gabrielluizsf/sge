package terrain

import "testing"

func climate() Climate {
	return Climate{Seed: 5, Scale: 40, Ragged: 0.06, Biomes: []Biome{
		{Name: "forest", Warmth: 0.3, Wetness: 0.7},
		{Name: "meadow", Warmth: 0.5, Wetness: 0.5},
		{Name: "dryland", Warmth: 0.75, Wetness: 0.2},
		{Name: "coast", Warmth: 0.6, Wetness: 0.9},
	}}
}

// Every biome is found over a wide land, each in large patches, most
// cells beside another of their own; the same cell is always the same.
func TestBiomesLieInPatches(t *testing.T) {
	c := climate()
	count := map[int]int{}
	same, pairs := 0, 0
	for r := -150; r < 150; r++ {
		for col := -150; col < 150; col++ {
			b := c.At(col, r)
			count[b]++
			if c.At(col, r) != b {
				t.Fatal("a cell changed its biome")
			}
			if c.At(col+1, r) == b {
				same++
			}
			pairs++
		}
	}
	for i, b := range c.Biomes {
		if count[i] < 300*300/50 {
			t.Errorf("%s covers only %d cells", b.Name, count[i])
		}
	}
	if share := float64(same) / float64(pairs); share < 0.85 {
		t.Errorf("only %.2f of cells lie beside one of their own biome", share)
	}
	if !c.In("forest")(firstOf(c, 0)) || c.In("nowhere")(0, 0) {
		t.Error("In does not agree with At")
	}
}

// firstOf is a cell of biome b.
func firstOf(c Climate, b int) (int, int) {
	for r := 0; ; r++ {
		for col := range 200 {
			if c.At(col, r) == b {
				return col, r
			}
		}
	}
}
