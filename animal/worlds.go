package animal

import "github.com/gabrielluizsf/sge/gltf"

// placed is where a node is in the scene: its place, turn and size (the
// scales models are made with are the same all ways).
type placed struct {
	pos  Vec
	rot  Quat
	size float64
}

// worlds are every node's place in the scene, and parents each node's
// parent, -1 for a root.
func worlds(m *gltf.Model) (out []placed, parents []int) {
	out = make([]placed, len(m.Nodes))
	parents = make([]int, len(m.Nodes))
	for i := range parents {
		parents[i] = -1
	}
	var walk func(i int, p placed)
	walk = func(i int, p placed) {
		n := m.Nodes[i]
		local := Vec(n.Translation)
		out[i] = placed{
			pos:  add(p.pos, scale(turn(p.rot, local), p.size)),
			rot:  mul(p.rot, Quat(n.Rotation)),
			size: p.size * n.Scale[0],
		}
		for _, c := range n.Children {
			parents[c] = i
			walk(c, out[i])
		}
	}
	for _, r := range m.Roots {
		walk(r, placed{rot: identity, size: 1})
	}
	return out, parents
}
