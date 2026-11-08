package human

import (
	"math"
	"strings"

	"github.com/gabrielluizsf/sge/gltf"
)

// swayBone is the name of the bone the hair's free ends hang from: a
// child of the head, at the scalp's pivot, which the animations swing.
const swayBone = "hair_sway"

// hairLength is how far, in metres, below the top of the hair a vertex is
// freest: it follows the bone entirely from there down, and not at all at
// the top.
const hairLength = 0.30

// swayNode is the index of m's hair bone, adding it, and the share of the
// hair's vertices that hang from it, if m has hair. It is false for a
// model that has none, or no head.
func swayNode(m *gltf.Model) (int, bool) {
	head, hair := -1, false
	for i, n := range m.Nodes {
		if n.Name == swayBone {
			return i, true
		}
		if n.Name == "head" {
			head = i
		}
	}
	for _, p := range hairPrimitives(m) {
		hair = hair || len(p.Positions) > 0
	}
	if head < 0 || !hair || len(m.Skins) == 0 {
		return 0, false
	}
	node := gltf.Node{Name: swayBone, Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{1, 1, 1}}
	m.Nodes = append(m.Nodes, node)
	id := len(m.Nodes) - 1
	m.Nodes[head].Children = append(m.Nodes[head].Children, id)
	sk := &m.Skins[0]
	joint := len(sk.Joints)
	var inverse [16]float64
	for k, j := range sk.Joints {
		if j == head {
			inverse = sk.InverseBind[k]
		}
	}
	sk.Joints = append(sk.Joints, id)
	sk.InverseBind = append(sk.InverseBind, inverse)
	hang(hairPrimitives(m), joint)
	return id, true
}

// hairPrimitives are the primitives of m's hair, which the pack names by
// the folder of its material.
func hairPrimitives(m *gltf.Model) []*gltf.Primitive {
	var out []*gltf.Primitive
	for mi := range m.Meshes {
		for pi := range m.Meshes[mi].Primitives {
			p := &m.Meshes[mi].Primitives[pi]
			if p.Material >= 0 && p.Material < len(m.Materials) && strings.HasPrefix(m.Materials[p.Material].Name, "hair/") {
				out = append(out, p)
			}
		}
	}
	return out
}

// hang gives the hair's vertices to the sway bone, joint of the skin, by
// how far below the top of the hair each is: the roots stay with the head,
// the ends are free.
func hang(hair []*gltf.Primitive, joint int) {
	top := math.Inf(-1)
	for _, p := range hair {
		for _, v := range p.Positions {
			top = max(top, float64(v[1]))
		}
	}
	for _, p := range hair {
		for i, v := range p.Positions {
			s := freedom((top - float64(v[1])) / hairLength)
			if s > 0 {
				give(p, i, joint, s)
			}
		}
	}
}

// freedom is how free a hair is, 0 to 1, at depth, its distance from the
// top in hair lengths: none at the roots, growing smoothly.
func freedom(depth float64) float64 {
	d := min(max(depth, 0), 1)
	return d * d * (3 - 2*d)
}

// give moves share s of vertex i's weight to joint, taking it from the
// others in proportion, the weakest dropped if all four places are full.
func give(p *gltf.Primitive, i, joint int, s float64) {
	j, w := p.Joints[i], p.Weights[i]
	slot := 0
	for k := range w {
		if w[k] < w[slot] {
			slot = k
		}
	}
	var sum float64
	for k := range w {
		if k != slot {
			w[k] *= float32(1 - s)
			sum += float64(w[k])
		}
	}
	j[slot], w[slot] = uint16(joint), float32(s)
	if sum+s > 0 {
		for k := range w {
			w[k] /= float32(sum + s)
		}
	}
	p.Joints[i], p.Weights[i] = j, w
}
