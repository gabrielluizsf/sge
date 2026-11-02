package human

import (
	"fmt"

	"github.com/gabrielluizsf/sge/gltf"
)

// Spill is a stream of water poured from a held can: the grip the can is
// in, where the tip of its spout is in the can's own frame, and which way
// the spout points there. The stream runs from the spout straight down to
// the ground, and is there only while the spout points down past Tilt, a
// least downward slope, 0 to 1, and the ground is Ground high.
type Spill struct {
	Can    Grip
	Spout  [3]float64
	Out    [3]float64
	Tilt   float64
	Ground float64
	During []string
}

// Pour gives person m the stream, a column of a metre standing from its
// origin, to show in the animations of s while the can is tipped. The
// can is held first, as Hold says, and the animations made.
func Pour(m, stream *gltf.Model, s Spill) error {
	at := len(m.Nodes)
	node := gltf.Node{Name: "stream", Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{tiny, tiny, tiny}}
	for _, r := range stream.Roots {
		node.Children = append(node.Children, at+1+r)
	}
	m.Nodes = append(m.Nodes, node)
	adopt(m, stream)
	m.Nodes[m.Roots[0]].Children = append(m.Nodes[m.Roots[0]].Children, at)
	for _, name := range s.During {
		if err := flow(m, at, name, s); err != nil {
			return err
		}
	}
	return nil
}

// flow gives node `at`, the stream, translation and scale channels in the
// animation called name: standing on the ground below the spout, as tall
// as the spout is high while it points down, and out of sight otherwise.
func flow(m *gltf.Model, node int, name string, s Spill) error {
	for i := range m.Animations {
		a := &m.Animations[i]
		if a.Name != name || len(a.Channels) == 0 {
			continue
		}
		times := a.Channels[0].Times
		move := gltf.Channel{Node: node, Path: "translation", Interpolation: "LINEAR", Times: times}
		size := gltf.Channel{Node: node, Path: "scale", Interpolation: "LINEAR", Times: times}
		for k := range times {
			tip, down := spoutAt(m, a, s, k)
			move.Values = append(move.Values, [4]float32{float32(tip[0]), float32(s.Ground), float32(tip[2])})
			if down {
				size.Values = append(size.Values, [4]float32{1, float32(tip[1] - s.Ground), 1})
			} else {
				size.Values = append(size.Values, [4]float32{tiny, tiny, tiny})
			}
		}
		a.Channels = append(a.Channels, move, size)
		return nil
	}
	return fmt.Errorf("human: no animation %s to pour in", name)
}

// spoutAt is where the spout's tip is in key k of animation a, and
// whether it points down past the tilt, above the ground.
func spoutAt(m *gltf.Model, a *gltf.Animation, s Spill, k int) (tip [3]float64, down bool) {
	at, q := handle(m, a, s.Can, k)
	out := q.turn(normalize(s.Out))
	tip = add(at, q.turn(s.Spout))
	return tip, -out[1] > s.Tilt && tip[1]-s.Ground > 0.02
}
