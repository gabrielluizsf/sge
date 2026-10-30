package human

import (
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// metre is how many of MakeHuman's units, decimetres, make a metre.
const metre = 10

// done is the person: their parts as one skinned mesh, in metres, the
// lowest of them on the ground, and the skeleton's bones as nodes under
// the root, whose rest is the shape of the body.
func (d *dresser) done() *gltf.Model {
	floor := math.Inf(1)
	for _, p := range d.parts {
		for _, v := range p.pos {
			floor = min(floor, v[1])
		}
	}
	place := func(v [3]float64) [3]float64 {
		v[1] -= floor
		return scale(v, 1.0/metre)
	}
	m := d.model
	bones := d.kit.skel.bones
	first := 2
	root := gltf.Node{Name: "human", Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{1, 1, 1}, Children: []int{1}}
	body := gltf.Node{Name: "body", Mesh: 0, Skin: 0, Rotation: [4]float64{3: 1}, Scale: [3]float64{1, 1, 1}}
	m.Nodes = append(m.Nodes, root, body)
	skin := gltf.Skin{}
	heads := make([][3]float64, len(bones))
	for i, b := range bones {
		heads[i] = place(d.kit.skel.joint(b.head, d.pos))
		n := gltf.Node{Name: b.name, Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{1, 1, 1}, Translation: heads[i]}
		if b.parent < 0 {
			m.Nodes[0].Children = append(m.Nodes[0].Children, first+i)
		} else {
			n.Translation = sub(heads[i], heads[b.parent])
			m.Nodes[first+b.parent].Children = append(m.Nodes[first+b.parent].Children, first+i)
		}
		m.Nodes = append(m.Nodes, n)
		skin.Joints = append(skin.Joints, first+i)
		skin.InverseBind = append(skin.InverseBind, [16]float64{0: 1, 5: 1, 10: 1, 15: 1, 12: -heads[i][0], 13: -heads[i][1], 14: -heads[i][2]})
	}
	m.Skins = []gltf.Skin{skin}
	mesh := gltf.Mesh{Name: "human"}
	for _, p := range d.parts {
		mesh.Primitives = append(mesh.Primitives, p.primitive(place))
	}
	m.Meshes = []gltf.Mesh{mesh}
	m.Roots = []int{0}
	return m
}

// primitive is the part as a primitive, its vertices placed by place.
func (p part) primitive(place func([3]float64) [3]float64) gltf.Primitive {
	g := gltf.Primitive{Material: p.mat}
	normals := p.normals()
	for i, v := range p.pos {
		at, n := place(v), normals[i]
		g.Positions = append(g.Positions, [3]float32{float32(at[0]), float32(at[1]), float32(at[2])})
		g.Normals = append(g.Normals, [3]float32{float32(n[0]), float32(n[1]), float32(n[2])})
		g.UVs = append(g.UVs, [2]float32{float32(p.uv[i][0]), float32(p.uv[i][1])})
		j, w := strongest(p.infl[i])
		g.Joints = append(g.Joints, j)
		g.Weights = append(g.Weights, w)
	}
	for _, t := range p.tris {
		g.Indices = append(g.Indices, uint32(t))
	}
	return g
}
