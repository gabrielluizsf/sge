package human

import (
	"fmt"

	"github.com/gabrielluizsf/sge/gltf"
)

// Fall is a stream of water running into a held can: the grip the can is
// in, where its mouth is in the can's own frame, and how tall the stream
// stands over it. The stream is there only from From to To of the
// animation, 0 to 1, and falls straight down onto the mouth.
type Fall struct {
	Can      Grip
	Mouth    [3]float64
	Height   float64
	From, To float64
	During   []string
}

// Fill gives person m the stream to show in the animations of f while the
// can is filled, a column of a metre standing from its origin. The can is
// held first, as Hold says, and the animations made.
func Fill(m, stream *gltf.Model, f Fall) error {
	at := attach(m, stream)
	for _, name := range f.During {
		if err := fall(m, at, name, f); err != nil {
			return err
		}
	}
	return nil
}

// fall gives node `at`, the stream, translation and scale channels in the
// animation called name: standing on the can's mouth, as tall as f says
// while the water runs, and out of sight otherwise.
func fall(m *gltf.Model, node int, name string, f Fall) error {
	for i := range m.Animations {
		a := &m.Animations[i]
		if a.Name != name || len(a.Channels) == 0 {
			continue
		}
		times := a.Channels[0].Times
		move := gltf.Channel{Node: node, Path: "translation", Interpolation: "LINEAR", Times: times}
		size := gltf.Channel{Node: node, Path: "scale", Interpolation: "LINEAR", Times: times}
		last := float64(times[len(times)-1])
		for k := range times {
			at, q := handle(m, a, f.Can, k)
			mouth := add(at, q.turn(f.Mouth))
			move.Values = append(move.Values, [4]float32{float32(mouth[0]), float32(mouth[1]), float32(mouth[2])})
			if u := float64(times[k]) / last; u >= f.From && u <= f.To {
				size.Values = append(size.Values, [4]float32{1, float32(f.Height), 1})
			} else {
				size.Values = append(size.Values, [4]float32{tiny, tiny, tiny})
			}
		}
		a.Channels = append(a.Channels, move, size)
		return nil
	}
	return fmt.Errorf("human: no animation %s to fill in", name)
}
