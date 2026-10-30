package flora

import (
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// vec is a point or direction in the model's space: metres, Y up.
type vec = [3]float64

// LeafShape is where a leaf grows and how it hangs: from Base, heading
// Yaw radians round the stem and rising Pitch radians above level, Length
// long and Width wide, its tip drooping Droop metres under its weight.
type LeafShape struct {
	Base       vec
	Yaw, Pitch float64
	Length     float64
	Width      float64
	Droop      float64
}

// leaf is a leaf's strip of quads along its length, bent as it droops,
// mapped onto a leaf texture: its stalk at the texture's bottom, its tip
// at the top; material mat.
func leaf(s LeafShape, mat int) gltf.Primitive {
	const steps = 6
	dir := vec{math.Cos(s.Pitch) * math.Sin(s.Yaw), math.Sin(s.Pitch), math.Cos(s.Pitch) * math.Cos(s.Yaw)}
	side := vec{math.Cos(s.Yaw), 0, -math.Sin(s.Yaw)}
	var p gltf.Primitive
	p.Material = mat
	for i := 0; i <= steps; i++ {
		t := float64(i) / steps
		c := add(s.Base, scale(dir, s.Length*t))
		c[1] -= s.Droop * t * t
		// The surface faces up across its axis.
		ahead := normalize(vec{dir[0], dir[1] - 2*s.Droop*t/math.Max(s.Length, 1e-6), dir[2]})
		up := normalize(cross(side, ahead))
		for k, u := range []float64{-0.5, 0.5} {
			q := add(c, scale(side, u*s.Width))
			p.Positions = append(p.Positions, f32(q))
			p.Normals = append(p.Normals, f32(up))
			p.UVs = append(p.UVs, [2]float32{float32(k), float32(1 - t)})
		}
		if i > 0 {
			b := uint32(2 * (i - 1))
			p.Indices = append(p.Indices, b, b+2, b+3, b, b+3, b+1)
		}
	}
	return p
}

// rod is a round stem from a to b, radius r, of material mat.
type rod struct {
	a, b vec
	r    float64
	mat  int
}

// stem is a rod's mesh.
func stem(s rod) gltf.Primitive {
	a, b, r, mat := s.a, s.b, s.r, s.mat
	const sides = 6
	axis := normalize(sub(b, a))
	u := normalize(cross(axis, vec{1, 0, 0}))
	if length(u) < 0.5 {
		u = normalize(cross(axis, vec{0, 0, 1}))
	}
	w := cross(axis, u)
	var p gltf.Primitive
	p.Material = mat
	for i := 0; i <= sides; i++ {
		ang := 2 * math.Pi * float64(i) / sides
		n := add(scale(u, math.Cos(ang)), scale(w, math.Sin(ang)))
		for k, end := range []vec{a, b} {
			p.Positions = append(p.Positions, f32(add(end, scale(n, r))))
			p.Normals = append(p.Normals, f32(n))
			p.UVs = append(p.UVs, [2]float32{float32(i) / sides, float32(k)})
		}
		if i > 0 {
			q := uint32(2 * (i - 1))
			p.Indices = append(p.Indices, q, q+2, q+3, q, q+3, q+1)
		}
	}
	return p
}

// Fruit is a round fruit: its middle, its radius across and up, and how
// deeply it is ribbed, as a pumpkin is, with ribs ribs round it.
type Fruit struct {
	Center vec
	Radius float64
	Squash float64 // its height as a share of its width
	Rib    float64
	Ribs   int
}

// fruit is a fruit's sphere, of material mat.
func fruit(f Fruit, mat int) gltf.Primitive {
	const rings, segs = 8, 12
	var p gltf.Primitive
	p.Material = mat
	for r := 0; r <= rings; r++ {
		th := math.Pi * float64(r) / rings
		for s := 0; s <= segs; s++ {
			ph := 2 * math.Pi * float64(s) / segs
			bulge := 1.0
			if f.Ribs > 0 {
				bulge = 1 - f.Rib*0.5*(1-math.Cos(float64(f.Ribs)*ph))
			}
			n := vec{math.Sin(th) * math.Cos(ph), math.Cos(th), math.Sin(th) * math.Sin(ph)}
			q := vec{n[0] * f.Radius * bulge, n[1] * f.Radius * f.Squash, n[2] * f.Radius * bulge}
			p.Positions = append(p.Positions, f32(add(f.Center, q)))
			p.Normals = append(p.Normals, f32(normalize(vec{n[0], n[1] / math.Max(f.Squash, 0.1), n[2]})))
			p.UVs = append(p.UVs, [2]float32{float32(s) / segs, float32(r) / rings})
			if r > 0 && s > 0 {
				a := uint32((r-1)*(segs+1) + s - 1)
				b := uint32(r*(segs+1) + s - 1)
				p.Indices = append(p.Indices, a, b, b+1, a, b+1, a+1)
			}
		}
	}
	return p
}

func add(a, b vec) vec           { return vec{a[0] + b[0], a[1] + b[1], a[2] + b[2]} }
func sub(a, b vec) vec           { return vec{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }
func scale(a vec, k float64) vec { return vec{a[0] * k, a[1] * k, a[2] * k} }
func length(a vec) float64       { return math.Sqrt(a[0]*a[0] + a[1]*a[1] + a[2]*a[2]) }
func f32(a vec) [3]float32       { return [3]float32{float32(a[0]), float32(a[1]), float32(a[2])} }

func normalize(a vec) vec {
	if l := length(a); l > 0 {
		return scale(a, 1/l)
	}
	return a
}

func cross(a, b vec) vec {
	return vec{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}
