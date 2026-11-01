package tree

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// extent is how tall m stands, and how wide: the larger of its spans
// along X and Z.
func extent(m *gltf.Model) (height, width float64) {
	lo, hi := [3]float64{math.Inf(1), math.Inf(1), math.Inf(1)}, [3]float64{math.Inf(-1), math.Inf(-1), math.Inf(-1)}
	for _, p := range m.Meshes[0].Primitives {
		for _, q := range p.Positions {
			for k, v := range q {
				if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
					return math.NaN(), math.NaN()
				}
				lo[k], hi[k] = math.Min(lo[k], float64(v)), math.Max(hi[k], float64(v))
			}
		}
	}
	return hi[1], math.Max(hi[0]-lo[0], hi[2]-lo[2])
}

// sum is every coordinate of a tree added up: two trees alike have the
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

// Every species at every size is a sound model, as tall (or wide) as the
// range its species says, whatever the seed.
func TestSizesAreRealistic(t *testing.T) {
	for _, s := range All() {
		for size := Sapling; size < Sizes; size++ {
			r := s.Metres[size]
			for seed := uint32(1); seed <= 3; seed++ {
				m := s.Grow(Spec{Size: size, Seed: seed})
				for _, p := range m.Meshes[0].Primitives {
					if p.Material < 0 || p.Material >= len(m.Materials) {
						t.Fatalf("%s: material %d of %d", s.Name, p.Material, len(m.Materials))
					}
					for _, i := range p.Indices {
						if int(i) >= len(p.Positions) {
							t.Fatalf("%s: index %d of %d", s.Name, i, len(p.Positions))
						}
					}
				}
				h, w := extent(m)
				got := h
				if s.Wide {
					got = w
				}
				if math.IsNaN(got) || got < r[0]-1e-3 || got > r[1]+1e-3 {
					t.Errorf("%s size %d seed %d: %.2f m, want %.1f to %.1f", s.Name, size, seed, got, r[0], r[1])
				}
			}
		}
	}
}

// A mature tree is taller than a young one, and that than a sapling.
func TestTreesGrow(t *testing.T) {
	for _, s := range All() {
		last := 0.0
		for size := Sapling; size < Sizes; size++ {
			h, w := extent(s.Grow(Spec{Size: size, Seed: 7}))
			if s.Wide {
				h = w
			}
			if h <= last {
				t.Errorf("%s size %d: %.2f m after %.2f", s.Name, size, h, last)
			}
			last = h
		}
	}
}

// Two seeds grow two trees, and one seed the same tree.
func TestSeedsDiffer(t *testing.T) {
	for _, s := range All() {
		a, b := s.Grow(Spec{Size: Young, Seed: 1}), s.Grow(Spec{Size: Young, Seed: 2})
		if sum(a) == sum(b) {
			t.Errorf("%s: two seeds grew alike", s.Name)
		}
		if sum(a) != sum(s.Grow(Spec{Size: Young, Seed: 1})) {
			t.Errorf("%s: a seed grew two trees", s.Name)
		}
	}
}

// Where a tree has cards, they are cut out where clear and textured:
// the cactus has none, every other species some.
func TestCardsAreCutOut(t *testing.T) {
	for _, s := range All() {
		cards := 0
		m := s.Grow(Spec{Size: Mature, Seed: 1})
		for _, mat := range m.Materials {
			switch mat.Name {
			case "needles", "leaf", "blossom", "frond":
				cards++
				if mat.AlphaMode != "MASK" || mat.Texture < 0 || mat.Texture >= len(m.Images) {
					t.Errorf("%s: %s is %s, texture %d", s.Name, mat.Name, mat.AlphaMode, mat.Texture)
				}
			default:
				if mat.AlphaMode != "OPAQUE" {
					t.Errorf("%s: %s is %s, want opaque", s.Name, mat.Name, mat.AlphaMode)
				}
			}
		}
		if (cards == 0) != (s.Name == "mandacaru") {
			t.Errorf("%s has %d card materials", s.Name, cards)
		}
	}
}

// A bark given is the trunk's material, with its textures; with none the
// trunk is a plain colour.
func TestBarkIsWorn(t *testing.T) {
	m := Pine.Grow(Spec{Size: Young, Seed: 1})
	if m.Materials[0].Name != "bark" || m.Materials[0].Texture != -1 {
		t.Errorf("plain bark is %+v", m.Materials[0])
	}
	img := newPen(4, 4).img
	m = Pine.Grow(Spec{Size: Young, Seed: 1, Barks: Barks{Pine: Bark{Albedo: img, Normal: img, ARM: img, Span: 2}}})
	b := m.Materials[0]
	if b.Texture < 0 || b.Normal < 0 || b.MetalRough < 0 {
		t.Errorf("worn bark is %+v", b)
	}
}
