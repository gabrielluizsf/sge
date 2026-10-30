package pose

import (
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// track is a channel ready to sample many times: what it moves, whether
// it steps, and, when its keys come at even steps from start, as most
// tools bake them, those steps, its key at a time then found by a
// division rather than a search.
type track struct {
	path        uint8
	stepped     bool
	even        bool
	start, step float64
}

// What a channel moves: none when its node is not the model's.
const (
	none uint8 = iota
	translation
	rotation
	scale
)

// tracksOf are the tracks of a's channels, in a model of nodes nodes.
func tracksOf(a *gltf.Animation, nodes int) []track {
	out := make([]track, len(a.Channels))
	for c := range a.Channels {
		ch := &a.Channels[c]
		if ch.Node < 0 || ch.Node >= nodes || len(ch.Values) < len(ch.Times) {
			continue
		}
		out[c] = trackOf(ch)
	}
	return out
}

// trackOf is ch's track, even when its keys are evenly spaced to a ten
// thousandth of a step.
func trackOf(ch *gltf.Channel) track {
	tr := track{path: map[string]uint8{"translation": translation, "rotation": rotation, "scale": scale}[ch.Path], stepped: ch.Interpolation == "STEP"}
	n := len(ch.Times)
	if n < 2 {
		return tr
	}
	start := float64(ch.Times[0])
	step := (float64(ch.Times[n-1]) - start) / float64(n-1)
	if !(step > 0) {
		return tr
	}
	for i, t := range ch.Times {
		if math.Abs(float64(t)-(start+float64(i)*step)) > step*1e-4 {
			return tr
		}
	}
	tr.even, tr.start, tr.step = true, start, step
	return tr
}

// sample is ch's value at t, as Sample has it.
func (tr *track) sample(ch *gltf.Channel, t float64) [4]float64 {
	if !tr.even {
		return Sample(ch, t)
	}
	x := (t - tr.start) / tr.step
	if !(x > 0) {
		return wide(ch.Values[0])
	}
	last := len(ch.Times) - 1
	i := int(x)
	if i >= last {
		return wide(ch.Values[last])
	}
	a, f := wide(ch.Values[i]), x-float64(i)
	if tr.stepped || f == 0 {
		return a
	}
	b := wide(ch.Values[i+1])
	if tr.path == rotation {
		return Slerp(a, b, f)
	}
	return [4]float64{a[0] + (b[0]-a[0])*f, a[1] + (b[1]-a[1])*f, a[2] + (b[2]-a[2])*f, a[3] + (b[3]-a[3])*f}
}
