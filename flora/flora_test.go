package flora

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// Every crop at every stage is a sound model: indices in range, numbers
// finite, materials and texture there; it grows no smaller as it ripens,
// and two seeds differ.
func TestCropsGrow(t *testing.T) {
	for _, c := range Crops {
		last := 0.0
		for s := range c.Stages {
			m := c.Grow(s, 1)
			top := 0.0
			for _, p := range m.Meshes[0].Primitives {
				if p.Material < 0 || p.Material >= len(m.Materials) {
					t.Fatalf("%s stage %d: material %d of %d", c.Name, s, p.Material, len(m.Materials))
				}
				for _, i := range p.Indices {
					if int(i) >= len(p.Positions) {
						t.Fatalf("%s stage %d: index %d of %d", c.Name, s, i, len(p.Positions))
					}
				}
				for _, q := range p.Positions {
					for _, v := range q {
						if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
							t.Fatalf("%s stage %d: a point is %v", c.Name, s, q)
						}
					}
					top = max(top, float64(q[1]))
				}
			}
			if top < last*0.9 {
				t.Errorf("%s shrank from %.2f to %.2f at stage %d", c.Name, last, top, s)
			}
			last = top
		}
		if sum(c.Grow(2, 1)) == sum(c.Grow(2, 2)) {
			t.Errorf("%s: two seeds grew alike", c.Name)
		}
	}
}

// A leaf's texture is a leaf: opaque along its middle, clear at the
// corners, green.
func TestLeavesArePainted(t *testing.T) {
	img := Carrot.leaf.texture(64, 3)
	mid := img.NRGBAAt(32, 32)
	if mid.A != 255 || mid.G <= mid.R || mid.G <= mid.B {
		t.Errorf("the middle is %v", mid)
	}
	if img.NRGBAAt(0, 0).A != 0 || img.NRGBAAt(63, 0).A != 0 {
		t.Error("the corners are painted")
	}
}

// sum is every coordinate of a plant added up: two plants alike have the
// same.
func sum(m *gltf.Model) float64 {
	var s float64
	for _, p := range m.Meshes[0].Primitives {
		for _, q := range p.Positions {
			s += float64(q[0]) + float64(q[1])*3 + float64(q[2])*7
		}
	}
	return s
}
