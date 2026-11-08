package human

import "github.com/gabrielluizsf/sge/gltf"

// spot is a node's place in the person's space: where it is and how it
// is turned.
type spot struct {
	at   [3]float64
	turn quat
}

// worldAt is where every node of m is in key k of animation a, each as
// its parent's and its own turns and translations, the animation's where
// it has them.
func worldAt(m *gltf.Model, a *gltf.Animation, k int) []spot {
	at := make([][3]float64, len(m.Nodes))
	turn := make([]quat, len(m.Nodes))
	for i, n := range m.Nodes {
		at[i], turn[i] = n.Translation, quat(n.Rotation)
	}
	for _, ch := range a.Channels {
		v := ch.Values[k]
		switch ch.Path {
		case "translation":
			at[ch.Node] = [3]float64{float64(v[0]), float64(v[1]), float64(v[2])}
		case "rotation":
			turn[ch.Node] = quat{float64(v[0]), float64(v[1]), float64(v[2]), float64(v[3])}
		}
	}
	out := make([]spot, len(m.Nodes))
	var walk func(i int, parent spot)
	walk = func(i int, parent spot) {
		out[i] = spot{add(parent.at, parent.turn.turn(at[i])), turn[i].then(parent.turn)}
		for _, c := range m.Nodes[i].Children {
			walk(c, out[i])
		}
	}
	for _, r := range m.Roots {
		walk(r, spot{turn: still})
	}
	return out
}

// palm is the middle of the palm of the hand whose wrist is the node
// wrist, in the wrist's frame at rest: from the wrist to the middle
// finger's metacarpal and half way on to its first knuckle, or the wrist
// itself where the hand has no fingers.
func palm(m *gltf.Model, wrist int) [3]float64 {
	side := m.Nodes[wrist].Name[len("wrist"):]
	meta := childNamed(m, wrist, "metacarpal3"+side)
	if meta < 0 {
		return [3]float64{}
	}
	p := m.Nodes[meta].Translation
	if k := childNamed(m, meta, "finger3-1"+side); k >= 0 {
		p = add(p, scale(m.Nodes[k].Translation, 0.5))
	}
	return p
}

// childNamed is the child of node parent called name, or -1.
func childNamed(m *gltf.Model, parent int, name string) int {
	for _, c := range m.Nodes[parent].Children {
		if m.Nodes[c].Name == name {
			return c
		}
	}
	return -1
}
