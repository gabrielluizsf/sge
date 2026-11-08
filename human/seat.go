package human

import (
	"fmt"

	"github.com/gabrielluizsf/sge/gltf"
)

// SeatPlace is where a piece of furniture stands in a person's space while
// they use it: At, in metres, from the spot under their hips, at Scale of
// its size, seen in the animations During and out of sight in the others.
type SeatPlace struct {
	At     [3]float64
	Scale  float64
	During []string
}

// Seat gives person m the chair to sit on, in their own space, so it turns
// with them and is baked with them: what is behind them hides them where it
// should from any side, as no tile drawn apart could. m's animations are
// made first. The chair's meshes, materials and images are added to m's.
func Seat(m, chair *gltf.Model, s SeatPlace) error {
	at := len(m.Nodes)
	node := gltf.Node{Name: "seat", Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{tiny, tiny, tiny},
		Translation: s.At}
	for _, r := range chair.Roots {
		node.Children = append(node.Children, at+1+r)
	}
	m.Nodes = append(m.Nodes, node)
	adopt(m, chair)
	m.Nodes[m.Roots[0]].Children = append(m.Nodes[m.Roots[0]].Children, at)
	for _, name := range s.During {
		if err := showSeat(m, at, name, s.Scale); err != nil {
			return err
		}
	}
	return nil
}

// showSeat shows node `at`, the chair, at full size all through the
// animation called name.
func showSeat(m *gltf.Model, node int, name string, scale float64) error {
	for i := range m.Animations {
		a := &m.Animations[i]
		if a.Name != name || len(a.Channels) == 0 {
			continue
		}
		size := gltf.Channel{Node: node, Path: "scale", Interpolation: "LINEAR", Times: a.Channels[0].Times}
		s := float32(scale)
		for range size.Times {
			size.Values = append(size.Values, [4]float32{s, s, s})
		}
		a.Channels = append(a.Channels, size)
		return nil
	}
	return fmt.Errorf("human: no animation %s to sit in", name)
}
