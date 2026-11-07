package human

import (
	"strings"
	"testing"
)

// skinAndBrows are the skin's and the eyebrows' meshes of m.
func skinAndBrows(t *testing.T, k *Kit, b Body, brow int) (skin, brows [][3]float32, normals [][3]float32) {
	t.Helper()
	l := Look{Eyes: "eyes/materials/brown.mhmat", Wear: []string{eyesProxy, eyebrows(brow)}}
	model, err := k.Make(b, l)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range model.Meshes[0].Primitives {
		name := model.Materials[p.Material].Name
		if strings.HasPrefix(name, "eyebrows/") {
			brows = p.Positions
		} else if name == "skin" || strings.HasPrefix(name, "skins/") {
			skin, normals = p.Positions, p.Normals
		}
	}
	return
}

// Eyebrows are laid on the skin and stand off it, for every brow, sex and
// age: none of their vertices sinks into the flesh.
func TestEyebrowsDoNotSinkIntoTheSkin(t *testing.T) {
	k := realKit(t)
	type c struct {
		sex, age float64
		brow     int
	}
	var cases []c
	for brow := 1; brow <= brows; brow++ {
		cases = append(cases, c{0, 0.5, brow}, c{1, 0.5, brow})
	}
	for _, age := range []float64{0.1, 0.25, 1} {
		cases = append(cases, c{0, age, 6}, c{1, age, 6})
	}
	for _, cs := range cases {
		b := Body{Gender: cs.sex, Age: cs.age, European: 1, Muscle: .5, Weight: .5, Height: .5, Proportions: .5}
		skin, brow, normals := skinAndBrows(t, k, b, cs.brow)
		if len(brow) == 0 || len(skin) == 0 {
			t.Fatalf("%+v: no skin or brows", cs)
		}
		for i, v := range brow {
			best, bd := 0, float32(1e18)
			for j, w := range skin {
				d := (v[0]-w[0])*(v[0]-w[0]) + (v[1]-w[1])*(v[1]-w[1]) + (v[2]-w[2])*(v[2]-w[2])
				if d < bd {
					bd, best = d, j
				}
			}
			n, w := normals[best], skin[best]
			if h := (v[0]-w[0])*n[0] + (v[1]-w[1])*n[1] + (v[2]-w[2])*n[2]; h < 0.0001 {
				t.Fatalf("%+v: brow vertex %d is %.5f m above the skin, in the flesh", cs, i, h)
			}
		}
	}
}
