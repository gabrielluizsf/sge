package human

import (
	"fmt"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// BedPlace is where a bed stands in the space of the person who lies on
// it: At, in metres, from the spot on the floor under their feet when
// they stand, turned Yaw radians about the vertical, at Scale of its
// size, seen in the animation During and out of sight in the others.
type BedPlace struct {
	At     [3]float64
	Yaw    float64
	Scale  float64
	During string
}

// Lay gives person m the bed to lie on, in the person's space but not on
// their bones, so it stays flat while they lie down: a second root of the
// model. It is baked with them, one picture, so the bed hides the body
// where it should from any side. m's animations are made first. The bed's
// meshes, materials and images are added to m's.
func Lay(m, bed *gltf.Model, p BedPlace) error {
	at := len(m.Nodes)
	node := gltf.Node{Name: "bed", Mesh: -1, Skin: -1, Rotation: [4]float64{0, math.Sin(p.Yaw / 2), 0, math.Cos(p.Yaw / 2)},
		Scale: [3]float64{tiny, tiny, tiny}, Translation: p.At}
	for _, r := range bed.Roots {
		node.Children = append(node.Children, at+1+r)
	}
	m.Nodes = append(m.Nodes, node)
	adopt(m, bed)
	m.Roots = append(m.Roots, at)
	for i := range m.Animations {
		a := &m.Animations[i]
		if a.Name != p.During || len(a.Channels) == 0 {
			continue
		}
		size := gltf.Channel{Node: at, Path: "scale", Interpolation: "LINEAR", Times: a.Channels[0].Times}
		for range size.Times {
			size.Values = append(size.Values, [4]float32{float32(p.Scale), float32(p.Scale), float32(p.Scale)})
		}
		a.Channels = append(a.Channels, size)
		return nil
	}
	return fmt.Errorf("human: no animation %s to lie in", p.During)
}
