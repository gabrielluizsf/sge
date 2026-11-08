package human

import (
	"fmt"

	"github.com/gabrielluizsf/sge/gltf"
)

// Grip is how a tool is held, in one hand or two: in the
// animations called During, and out of sight in the others. Hands are the
// bones of the hand at the handle's end and of the hand up it, such as
// "wrist.L" and "wrist.R". The handle lies from the first palm through the
// second, the tool's X axis along it with the first palm at End on it, in
// metres, and its Z axis, the way its blade faces, as near as it can be
// to Facing, a direction in the person's space. A handle swung down is
// kept from pointing along Facing, by Flat of the way, 0 to 1, as the
// hands alone do not tell a wrist's turn: the blade then faces it.
//
// A tool held in one hand has the same bone twice in Hands, and lies
// where its Offset and Turn say, in the frame of that bone: its origin
// at Offset from the wrist, in metres, turned by Turn, a quaternion.
type Grip struct {
	Hands  [2]string
	During []string
	End    float64
	Facing [3]float64
	Flat   float64
	Offset [3]float64
	Turn   [4]float64
}

// tiny is the scale of a tool put away, too small to see.
const tiny = 1e-4

// Hold gives person m the tool to hold as g says. The tool is put, in
// each pose of the animations it is held in, where the hands put it, so
// the hand at its end never lets go. m's animations are made first. The
// tool's meshes, materials and images are added to m's.
func Hold(m, tool *gltf.Model, g Grip) error {
	for _, h := range g.Hands {
		if boneOf(m, h) < 0 {
			return fmt.Errorf("human: no bone %s to hold in", h)
		}
	}
	at := len(m.Nodes)
	grip := gltf.Node{Name: "grip", Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{tiny, tiny, tiny}}
	for _, r := range tool.Roots {
		grip.Children = append(grip.Children, at+1+r)
	}
	m.Nodes = append(m.Nodes, grip)
	adopt(m, tool)
	m.Nodes[m.Roots[0]].Children = append(m.Nodes[m.Roots[0]].Children, at)
	for _, name := range g.During {
		if err := place(m, at, name, g); err != nil {
			return err
		}
	}
	return nil
}

// boneOf is the index of the node called name in m, or -1.
func boneOf(m *gltf.Model, name string) int {
	for i, n := range m.Nodes {
		if n.Name == name {
			return i
		}
	}
	return -1
}

// adopt adds the tool's nodes, meshes, materials and images to m's, their
// indices moved to where they now are; the tool's node 0 is m's node
// len(m.Nodes)+1 as it was before.
func adopt(m, tool *gltf.Model) {
	nodes, meshes, mats, images := len(m.Nodes), len(m.Meshes), len(m.Materials), len(m.Images)
	for _, n := range tool.Nodes {
		n.Children = append([]int(nil), n.Children...)
		for i := range n.Children {
			n.Children[i] += nodes
		}
		if n.Mesh >= 0 {
			n.Mesh += meshes
		}
		n.Skin = -1
		m.Nodes = append(m.Nodes, n)
	}
	for _, me := range tool.Meshes {
		prims := make([]gltf.Primitive, len(me.Primitives))
		copy(prims, me.Primitives)
		for i := range prims {
			if prims[i].Material >= 0 {
				prims[i].Material += mats
			}
		}
		m.Meshes = append(m.Meshes, gltf.Mesh{Name: me.Name, Primitives: prims})
	}
	for _, ma := range tool.Materials {
		for _, t := range []*int{&ma.Texture, &ma.MetalRough, &ma.Normal, &ma.Occlusion} {
			if *t >= 0 {
				*t += images
			}
		}
		m.Materials = append(m.Materials, ma)
	}
	m.Images = append(m.Images, tool.Images...)
}
