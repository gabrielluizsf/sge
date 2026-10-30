package bake

import (
	"sort"

	"github.com/gabrielluizsf/sge/gltf"
)

// Shot is what to render of a model: an animation by name ("" for none)
// and the seconds into it, looped round its length; Yaw, the degrees it
// is turned round the vertical, counterclockwise seen from above, from
// facing straight at the viewer; and, when Only is set, only the scene's
// root Root, stood on its own foot: one of the variants a model of
// plants lays out side by side.
type Shot struct {
	Animation string
	Time      float64
	Yaw       float64
	Only      bool
	Root      int
}

// worlds are every node's matrix in the model's space, as shot: its parents'
// times its own, as the animation moves it at that moment.
func worlds(m *gltf.Model, p Shot) []mat {
	local := make([][3][4]float64, len(m.Nodes)) // translation, rotation, scale
	for i, n := range m.Nodes {
		local[i] = [3][4]float64{{n.Translation[0], n.Translation[1], n.Translation[2]}, n.Rotation, {n.Scale[0], n.Scale[1], n.Scale[2]}}
	}
	if a := animation(m, p.Animation); a != nil {
		at := p.Time
		if l := length(a); l > 0 {
			at -= l * float64(int(at/l))
		}
		for _, ch := range a.Channels {
			k := map[string]int{"translation": 0, "rotation": 1, "scale": 2}[ch.Path]
			local[ch.Node][k] = sample(ch, at)
		}
	}
	out := make([]mat, len(m.Nodes))
	var walk func(i int, parent mat)
	walk = func(i int, parent mat) {
		l := local[i]
		out[i] = parent.mul(trs([3]float64{l[0][0], l[0][1], l[0][2]}, l[1], [3]float64{l[2][0], l[2][1], l[2][2]}))
		for _, c := range m.Nodes[i].Children {
			walk(c, out[i])
		}
	}
	for i, r := range m.Roots {
		if p.Only && i != p.Root {
			continue
		}
		base := identity
		if p.Only {
			// Stood on its own foot: its place across the ground undone.
			t := m.Nodes[r].Translation
			base = trs([3]float64{-t[0], 0, -t[2]}, [4]float64{0, 0, 0, 1}, [3]float64{1, 1, 1})
		}
		walk(r, base)
	}
	return out
}

// animation is the model's animation called name, or nil.
func animation(m *gltf.Model, name string) *gltf.Animation {
	for i := range m.Animations {
		if name != "" && m.Animations[i].Name == name {
			return &m.Animations[i]
		}
	}
	return nil
}

// length is how long an animation runs: its last keyframe's time.
func length(a *gltf.Animation) float64 {
	var l float64
	for _, ch := range a.Channels {
		if n := len(ch.Times); n > 0 {
			l = max(l, float64(ch.Times[n-1]))
		}
	}
	return l
}

// sample is a channel's value at time t: held before its first key and
// after its last, stepped or eased between them.
func sample(ch gltf.Channel, t float64) [4]float64 {
	val := func(i int) [4]float64 {
		v := ch.Values[i]
		return [4]float64{float64(v[0]), float64(v[1]), float64(v[2]), float64(v[3])}
	}
	n := len(ch.Times)
	if n == 0 {
		return [4]float64{}
	}
	i := sort.Search(n, func(i int) bool { return float64(ch.Times[i]) > t })
	if i == 0 {
		return val(0)
	}
	if i == n || ch.Interpolation == "STEP" {
		return val(i - 1)
	}
	t0, t1 := float64(ch.Times[i-1]), float64(ch.Times[i])
	f := (t - t0) / max(t1-t0, 1e-9)
	a, b := val(i-1), val(i)
	if ch.Path == "rotation" {
		return slerp(a, b, f)
	}
	var out [4]float64
	for k := range out {
		out[k] = a[k] + (b[k]-a[k])*f
	}
	return out
}

// Length is how many seconds the animation name of m runs, 0 if m has no
// such animation.
func Length(m *gltf.Model, name string) float64 {
	if a := animation(m, name); a != nil {
		return length(a)
	}
	return 0
}
