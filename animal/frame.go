package animal

import (
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// Frame is the turn about the vertical and the move that bring a model to
// face +Z, stand on Y=0 and centre on X and Z, as package bake wants.
type Frame struct {
	Turn Quat
	Move Vec
}

// Bounds are the lowest and highest corners of the model's meshes as
// they lie in the scene, rigid meshes where their nodes put them, skinned
// ones as their vertices are, as a bind pose.
func Bounds(m *gltf.Model) (lo, hi Vec) {
	lo, hi = Vec{1e18, 1e18, 1e18}, Vec{-1e18, -1e18, -1e18}
	w, _ := worlds(m)
	each(m, func(n int, p *gltf.Primitive) {
		for _, v := range p.Positions {
			q := Vec{float64(v[0]), float64(v[1]), float64(v[2])}
			if m.Nodes[n].Skin < 0 {
				q = add(w[n].pos, scale(turn(w[n].rot, q), w[n].size))
			}
			for k := range q {
				lo[k], hi[k] = min(lo[k], q[k]), max(hi[k], q[k])
			}
		}
	})
	return lo, hi
}

// each calls f for every primitive of every node showing a mesh.
func each(m *gltf.Model, f func(node int, p *gltf.Primitive)) {
	for n, node := range m.Nodes {
		if node.Mesh < 0 {
			continue
		}
		for i := range m.Meshes[node.Mesh].Primitives {
			f(n, &m.Meshes[node.Mesh].Primitives[i])
		}
	}
}

// FrameFor is the frame that turns the model so that forward, a
// direction along the ground in the scene, points to +Z, and stands it.
func FrameFor(m *gltf.Model, forward Vec) Frame {
	f := Frame{Turn: About(Vec{0, 1, 0}, math.Atan2(-forward[0], forward[2]))}
	lo, hi := Bounds(m)
	low, high := Vec{1e18, 1e18, 1e18}, Vec{-1e18, -1e18, -1e18}
	for _, x := range []float64{lo[0], hi[0]} {
		for _, y := range []float64{lo[1], hi[1]} {
			for _, z := range []float64{lo[2], hi[2]} {
				q := turn(f.Turn, Vec{x, y, z})
				for k := range q {
					low[k], high[k] = min(low[k], q[k]), max(high[k], q[k])
				}
			}
		}
	}
	f.Move = Vec{-(low[0] + high[0]) / 2, -low[1], -(low[2] + high[2]) / 2}
	return f
}

// Wrap puts the frame over the whole scene, as a node above its roots,
// so a model with a skeleton already is brought to the frame.
func Wrap(m *gltf.Model, f Frame) {
	m.Nodes = append(m.Nodes, gltf.Node{
		Name: "frame", Children: m.Roots, Mesh: -1, Skin: -1,
		Translation: f.Move, Rotation: f.Turn, Scale: [3]float64{1, 1, 1},
	})
	m.Roots = []int{len(m.Nodes) - 1}
}

// Bake puts rigid meshes where their nodes put them and then in the frame,
// in their vertices, and leaves a node for each at the scene's roots; for a
// model with no skeleton, to be given one in the frame's axes.
func Bake(m *gltf.Model, f Frame) {
	w, _ := worlds(m)
	var roots []int
	var nodes []gltf.Node
	for n, node := range m.Nodes {
		if node.Mesh < 0 {
			continue
		}
		for i := range m.Meshes[node.Mesh].Primitives {
			p := &m.Meshes[node.Mesh].Primitives[i]
			place := func(v Vec, dir bool) Vec {
				v = turn(w[n].rot, v)
				if !dir {
					v = add(w[n].pos, scale(v, w[n].size))
				}
				v = turn(f.Turn, v)
				if !dir {
					v = add(v, f.Move)
				}
				return v
			}
			for k, v := range p.Positions {
				p.Positions[k] = f3(place(Vec{float64(v[0]), float64(v[1]), float64(v[2])}, false))
			}
			for k, v := range p.Normals {
				p.Normals[k] = f3(place(Vec{float64(v[0]), float64(v[1]), float64(v[2])}, true))
			}
			for k, v := range p.Tangents {
				d := place(Vec{float64(v[0]), float64(v[1]), float64(v[2])}, true)
				p.Tangents[k] = [4]float32{float32(d[0]), float32(d[1]), float32(d[2]), v[3]}
			}
		}
		roots = append(roots, len(nodes))
		nodes = append(nodes, gltf.Node{Name: node.Name, Mesh: node.Mesh, Skin: -1, Rotation: identity, Scale: [3]float64{1, 1, 1}})
	}
	m.Nodes, m.Roots = nodes, roots
}

func f3(v Vec) [3]float32 { return [3]float32{float32(v[0]), float32(v[1]), float32(v[2])} }
