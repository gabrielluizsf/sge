package building

import "math"

// face is a flat, convex polygon of a solid: its corners, counterclockwise
// seen from outside; its paint; and its grain, the way a texture's grain
// runs where the face allows, as a plank's runs along it, down the image
// or across it as its surface says, or zero for the image's down down
// the face, or, on a level face, toward the front.
type face struct {
	corners []vec
	paint   paint
	grain   vec
}

// face adds f, fanned into triangles, flat shaded, its texture laid on
// at its surface's scale: each repeat its span in metres, wherever the
// face is, so neighbouring faces in a plane run on without a seam. Its
// UVs are moved by whole repeats to start between 0 and 1.
func (b *builder) face(f face) {
	n := newell(f.corners)
	if n == (vec{}) {
		return
	}
	i := b.primitive(f.paint)
	span := b.laid[i].Span
	down := downOf(n, f.grain)
	right := cross(n, down)
	if b.laid[i].Across && f.grain != (vec{}) {
		// The grain runs across the image: lay the image's across along it.
		right = down
		down = cross(right, n)
	}
	uvs := make([][2]float64, len(f.corners))
	lo := [2]float64{math.Inf(1), math.Inf(1)}
	for k, c := range f.corners {
		uvs[k] = [2]float64{dot(c, right) / span, dot(c, down) / span}
		lo = [2]float64{min(lo[0], uvs[k][0]), min(lo[1], uvs[k][1])}
	}
	shift := [2]float64{math.Floor(lo[0]), math.Floor(lo[1])}
	p := &b.prims[i]
	base := uint32(len(p.Positions))
	for k, c := range f.corners {
		p.Positions = append(p.Positions, f32(c))
		p.Normals = append(p.Normals, f32(n))
		p.UVs = append(p.UVs, [2]float32{float32(uvs[k][0] - shift[0]), float32(uvs[k][1] - shift[1])})
	}
	for k := 1; k+1 < len(f.corners); k++ {
		p.Indices = append(p.Indices, base, base+uint32(k), base+uint32(k+1))
	}
}

// downOf is the way a texture's down runs on a face facing n: along the
// grain, where it lies in the face; else down the face; else, on a level
// face, toward the front.
func downOf(n, grain vec) vec {
	flat := func(v vec) vec { return sub(v, scale(n, dot(v, n))) }
	if g := flat(normalize(grain)); length(g) > 0.2 {
		return normalize(g)
	}
	if g := flat(vec{0, -1, 0}); length(g) > 0.05 {
		return normalize(g)
	}
	return normalize(flat(vec{0, 0, 1}))
}
