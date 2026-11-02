package building

import "testing"

// Every item is a sound model, plain or textured, small enough to hold,
// and lying on the ground.
func TestItemsAreSound(t *testing.T) {
	for _, palette := range []Palette{Plain(), textured()} {
		for _, k := range Items {
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
