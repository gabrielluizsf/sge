package animal

import "github.com/gabrielluizsf/sge/gltf"

// DropRigid hides every mesh not bound to a skin: the pieces that are
// only placed by their nodes, which a model may leave far from its body.
func DropRigid(m *gltf.Model) {
	for i := range m.Nodes {
		if m.Nodes[i].Skin < 0 {
			m.Nodes[i].Mesh = -1
		}
	}
}

// Filter keeps only the triangles keep allows, by the name of the
// triangle's material and the centre of its corners, and the vertices they
// use, so what is left is no larger than what shows. A mesh a node shares
// with another is cut for both.
func Filter(m *gltf.Model, keep func(material string, centre Vec) bool) {
	for mi := range m.Meshes {
		for pi := range m.Meshes[mi].Primitives {
			p := &m.Meshes[mi].Primitives[pi]
			name := ""
			if p.Material >= 0 && p.Material < len(m.Materials) {
				name = m.Materials[p.Material].Name
			}
			cut(p, func(a, b, c uint32) bool {
				var s Vec
				for _, i := range []uint32{a, b, c} {
					v := p.Positions[i]
					s = add(s, Vec{float64(v[0]), float64(v[1]), float64(v[2])})
				}
				return keep(name, scale(s, 1.0/3))
			})
		}
	}
}

// cut keeps of p's triangles those keep allows, and the vertices they use.
func cut(p *gltf.Primitive, keep func(a, b, c uint32) bool) {
	remap := make([]int64, len(p.Positions))
	for i := range remap {
		remap[i] = -1
	}
	var idx []uint32
	n := int64(0)
	var used []uint32
	for t := 0; t+2 < len(p.Indices); t += 3 {
		a, b, c := p.Indices[t], p.Indices[t+1], p.Indices[t+2]
		if !keep(a, b, c) {
			continue
		}
		for _, v := range []uint32{a, b, c} {
			if remap[v] < 0 {
				remap[v] = n
				n++
				used = append(used, v)
			}
			idx = append(idx, uint32(remap[v]))
		}
	}
	p.Positions = pick(p.Positions, used)
	p.Normals = pick(p.Normals, used)
	p.Tangents = pick(p.Tangents, used)
	p.UVs = pick(p.UVs, used)
	p.Joints = pick(p.Joints, used)
	p.Weights = pick(p.Weights, used)
	p.Indices = idx
}

// pick is the elements of s at used, or none if s has none.
func pick[T any](s []T, used []uint32) []T {
	if len(s) == 0 {
		return nil
	}
	out := make([]T, len(used))
	for i, u := range used {
		out[i] = s[u]
	}
	return out
}
