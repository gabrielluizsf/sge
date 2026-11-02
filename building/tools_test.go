package building

import "testing"

// Every item, and every smith's tool, is a sound model, plain or textured, small enough to hold,
// and lying on the ground.
func TestItemsAreSound(t *testing.T) {
	for _, palette := range []Palette{Plain(), textured()} {
		for _, k := range append(append([]Building{}, Items...), Smithing...) {
			m := k.Model(palette)
			if len(m.Nodes) != 1 || len(m.Meshes) != 1 {
				t.Fatalf("%s: %d nodes, %d meshes", k.Name, len(m.Nodes), len(m.Meshes))
			}
			lo, hi := 1e9, -1e9
			for _, p := range m.Meshes[0].Primitives {
				sound(t, soundOf{k.Name, m, p})
				for _, q := range p.Positions {
					lo, hi = min(lo, float64(q[1])), max(hi, float64(q[1]))
					for _, c := range q {
						if c < -1.5 || c > 1.5 {
							t.Errorf("%s is %v away from its origin", k.Name, c)
						}
					}
				}
			}
			if lo < -0.01 || hi > 0.5 {
				t.Errorf("%s spans heights %v to %v", k.Name, lo, hi)
			}
		}
	}
}

// extent is a model's size along each axis, in metres.
func extent(k Building) (lo, hi [3]float64) {
	lo, hi = [3]float64{1e9, 1e9, 1e9}, [3]float64{-1e9, -1e9, -1e9}
	for _, p := range k.Model(Plain()).Meshes[0].Primitives {
		for _, q := range p.Positions {
			for c := range 3 {
				lo[c], hi[c] = min(lo[c], float64(q[c])), max(hi[c], float64(q[c]))
			}
		}
	}
	return lo, hi
}

// The smith's tools are life size: a hammer with a handle of about 35 cm
// and a head of about 12 cm across, a bar of iron of about 30 cm, tongs
// of about 40 cm that hold it by its cold end, its far end glowing.
func TestSmithingToolsAreLifeSize(t *testing.T) {
	for _, c := range []struct {
		k        Building
		long     [2]float64 // its length along X, least and most
		wide     [2]float64 // its width along Z
		thick    float64    // the most it rises
		glow     bool
		wantName string
	}{
		{Hammer, [2]float64{0.30, 0.42}, [2]float64{0.10, 0.16}, 0.07, false, "hammer"},
		{IronBar, [2]float64{0.28, 0.32}, [2]float64{0.015, 0.04}, 0.04, true, "iron_bar"},
		{Tongs, [2]float64{0.35, 0.45}, [2]float64{0.03, 0.12}, 0.06, false, "tongs"},
		{TongsIron, [2]float64{0.45, 0.75}, [2]float64{0.03, 0.12}, 0.06, true, "tongs_iron"},
	} {
		if c.k.Name != c.wantName {
			t.Errorf("%s is called %s", c.wantName, c.k.Name)
		}
		lo, hi := extent(c.k)
		if l := hi[0] - lo[0]; l < c.long[0] || l > c.long[1] {
			t.Errorf("%s is %.2f m long", c.k.Name, l)
		}
		if w := hi[2] - lo[2]; w < c.wide[0] || w > c.wide[1] {
			t.Errorf("%s is %.2f m wide", c.k.Name, w)
		}
		if hi[1] > c.thick {
			t.Errorf("%s rises %.2f m", c.k.Name, hi[1])
		}
		if c.glow != glows(c.k) {
			t.Errorf("%s glows: %v", c.k.Name, glows(c.k))
		}
	}
}

// glows is whether k has an orange, hot part: a surface far redder than
// it is blue.
func glows(k Building) bool {
	m := k.Model(Plain())
	for _, mt := range m.Materials {
		if mt.BaseColor[0] > 0.8 && mt.BaseColor[2] < 0.25 && mt.BaseColor[1] > 0.15 && mt.BaseColor[1] < 0.6 {
			return true
		}
	}
	return false
}
