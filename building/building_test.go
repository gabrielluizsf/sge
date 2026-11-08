package building

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// textured is a palette of every material textured, as Poly Haven's are,
// with small images standing in for theirs, each a span of its own.
func textured() Palette {
	p := Plain()
	img := func() image.Image { return image.NewNRGBA(image.Rect(0, 0, 4, 4)) }
	for m := range materials {
		p[m] = Surface{Albedo: img(), ARM: img(), Normal: img(), Span: 0.5 + float64(m)*0.25,
			Color: color.NRGBA{255, 255, 255, 255}, Roughness: 1, Metallic: 1}
	}
	return p
}

// Every building is a sound model, plain or textured: one node, its
// indices in range, its numbers finite, its normals of length 1 and
// facing the way its triangles wind, each texture repeat its span.
func TestBuildingsAreSound(t *testing.T) {
	for _, palette := range []Palette{Plain(), textured()} {
		for _, k := range All() {
			m := k.Model(palette)
			if len(m.Nodes) != 1 || len(m.Meshes) != 1 || len(m.Roots) != 1 {
				t.Fatalf("%s: %d nodes, %d meshes", k.Name, len(m.Nodes), len(m.Meshes))
			}
			for _, p := range m.Meshes[0].Primitives {
				sound(t, soundOf{k.Name, m, p})
			}
		}
	}
}

// soundOf is a primitive of a building's model, to check.
type soundOf struct {
	name string
	m    *gltf.Model
	p    gltf.Primitive
}

// sound checks a primitive as TestBuildingsAreSound says.
func sound(t *testing.T, s soundOf) {
	t.Helper()
	p := s.p
	if p.Material < 0 || p.Material >= len(s.m.Materials) {
		t.Fatalf("%s: material %d of %d", s.name, p.Material, len(s.m.Materials))
	}
	if len(p.Normals) != len(p.Positions) || len(p.UVs) != len(p.Positions) || len(p.Indices)%3 != 0 {
		t.Fatalf("%s: %d points, %d normals, %d UVs, %d indices", s.name, len(p.Positions), len(p.Normals), len(p.UVs), len(p.Indices))
	}
	for i, q := range p.Positions {
		n := p.Normals[i]
		for _, v := range append(q[:], append(n[:], p.UVs[i][:]...)...) {
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
				t.Fatalf("%s: point %v, normal %v, UV %v", s.name, q, n, p.UVs[i])
			}
		}
		if l := length(f64(n)); math.Abs(l-1) > 1e-3 {
			t.Fatalf("%s: a normal is %v long", s.name, l)
		}
		// Moved to start between 0 and 1, a face's UVs reach no further
		// than the largest face, 13 m, over the smallest span, 0.5 m.
		if uv := p.UVs[i]; uv[0] < 0 || uv[1] < 0 || uv[0] > 27 || uv[1] > 27 {
			t.Fatalf("%s: UV %v out of range", s.name, uv)
		}
	}
	for k := 0; k < len(p.Indices); k += 3 {
		i, j, l := p.Indices[k], p.Indices[k+1], p.Indices[k+2]
		if int(max(i, j, l)) >= len(p.Positions) {
			t.Fatalf("%s: index past %d points", s.name, len(p.Positions))
		}
		a, b, c := f64(p.Positions[i]), f64(p.Positions[j]), f64(p.Positions[l])
		w := cross(sub(b, a), sub(c, a))
		if length(w) > 1e-7 && dot(w, f64(p.Normals[i])) <= 0 {
			t.Fatalf("%s: a triangle winds against its normal", s.name)
		}
	}
}

func f64(v [3]float32) vec { return vec{float64(v[0]), float64(v[1]), float64(v[2])} }
