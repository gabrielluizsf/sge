package gltf

// trianglesMode is the only primitive mode read: separate triangles.
const trianglesMode = 4

// meshes are the model's meshes, each primitive's attributes read and
// its indices checked against its vertices.
func meshes(doc *document, v views, nMaterials int) ([]Mesh, error) {
	var out []Mesh
	for mi, dm := range doc.Meshes {
		mesh := Mesh{Name: dm.Name}
		for pi, dp := range dm.Primitives {
			if dp.Mode != nil && *dp.Mode != trianglesMode {
				continue // points and lines are not drawn
			}
			p, err := primitive(v, dp, nMaterials)
			if err != nil {
				return nil, malformed("mesh %d primitive %d: %v", mi, pi, err)
			}
			mesh.Primitives = append(mesh.Primitives, p)
		}
		out = append(out, mesh)
	}
	return out, nil
}

// primitive reads one primitive: positions it must have, the rest when
// there, all of as many vertices, and indices into them, three a
// triangle; with none, its vertices taken three at a time.
func primitive(v views, dp docPrimitive, nMaterials int) (Primitive, error) {
	p := Primitive{Material: -1}
	if dp.Material != nil {
		if *dp.Material < 0 || *dp.Material >= nMaterials {
			return p, malformed("material %d", *dp.Material)
		}
		p.Material = *dp.Material
	}
	pos, ok := dp.Attributes["POSITION"]
	if !ok {
		return p, malformed("no positions")
	}
	a, err := v.accessor(pos, "VEC3")
	if err != nil {
		return p, err
	}
	n := a.count
	p.Positions = vec3s(a)
	if i, ok := dp.Attributes["NORMAL"]; ok {
		if p.Normals, err = attribute(v, i, want[[3]float32]{n, vec3s, "VEC3"}); err != nil {
			return p, err
		}
	}
	if i, ok := dp.Attributes["TANGENT"]; ok {
		if p.Tangents, err = attribute(v, i, want[[4]float32]{n, vec4s, "VEC4"}); err != nil {
			return p, err
		}
	}
	if i, ok := dp.Attributes["TEXCOORD_0"]; ok {
		if p.UVs, err = attribute(v, i, want[[2]float32]{n, vec2s, "VEC2"}); err != nil {
			return p, err
		}
	}
	if i, ok := dp.Attributes["WEIGHTS_0"]; ok {
		if p.Weights, err = attribute(v, i, want[[4]float32]{n, vec4s, "VEC4"}); err != nil {
			return p, err
		}
	}
	if i, ok := dp.Attributes["JOINTS_0"]; ok {
		ja, err := v.accessor(i, "VEC4")
		if err != nil || ja.count != n || !ja.integer() {
			return p, malformed("joints")
		}
		p.Joints = make([][4]uint16, n)
		for e := range n {
			for c := range 4 {
				p.Joints[e][c] = uint16(min(ja.uint(e, c), 0xffff))
			}
		}
	}
	p.Indices, err = indices(v, dp, n)
	return p, err
}

// indices reads the primitive's indices, each below n, or makes them.
func indices(v views, dp docPrimitive, n int) ([]uint32, error) {
	if dp.Indices == nil {
		if n%3 != 0 {
			return nil, malformed("%d vertices are no triangles", n)
		}
		out := make([]uint32, n)
		for i := range out {
			out[i] = uint32(i)
		}
		return out, nil
	}
	a, err := v.accessor(*dp.Indices, "SCALAR")
	if err != nil || !a.integer() || a.count%3 != 0 {
		return nil, malformed("indices")
	}
	out := make([]uint32, a.count)
	for e := range a.count {
		x := a.uint(e, 0)
		if int64(x) >= int64(n) {
			return nil, malformed("index %d of %d vertices", x, n)
		}
		out[e] = x
	}
	return out, nil
}

// want is what an attribute must be: n elements of kind, read by conv.
type want[T any] struct {
	n    int
	conv func(accessor) []T
	kind string
}

// attribute is accessor i read as w wants.
func attribute[T any](v views, i int, w want[T]) ([]T, error) {
	n, conv := w.n, w.conv
	a, err := v.accessor(i, w.kind)
	if err != nil {
		return nil, err
	}
	if a.count != n {
		return nil, malformed("%d values for %d vertices", a.count, n)
	}
	return conv(a), nil
}

func vec2s(a accessor) [][2]float32 {
	out := make([][2]float32, a.count)
	for e := range out {
		out[e] = [2]float32{float32(a.float(e, 0)), float32(a.float(e, 1))}
	}
	return out
}

func vec3s(a accessor) [][3]float32 {
	out := make([][3]float32, a.count)
	for e := range out {
		out[e] = [3]float32{float32(a.float(e, 0)), float32(a.float(e, 1)), float32(a.float(e, 2))}
	}
	return out
}

func vec4s(a accessor) [][4]float32 {
	out := make([][4]float32, a.count)
	for e := range out {
		for c := range 4 {
			out[e][c] = float32(a.float(e, c))
		}
	}
	return out
}
