package animal

import (
	"math"
	"sort"

	"github.com/gabrielluizsf/sge/gltf"
)

// Track is where bone b is in the scene, at each of times, in the
// animation called name of the rig's model, turns eased by their
// shorter way and not along a curve.
func (r *Rig) Track(name string, b int, times []float64) []Vec {
	var anim *gltf.Animation
	for i := range r.M.Animations {
		if r.M.Animations[i].Name == name {
			anim = &r.M.Animations[i]
		}
	}
	out := make([]Vec, len(times))
	for k, t := range times {
		local := make([][3][4]float64, len(r.M.Nodes))
		for i, n := range r.M.Nodes {
			local[i] = [3][4]float64{{n.Translation[0], n.Translation[1], n.Translation[2]}, n.Rotation, {n.Scale[0], n.Scale[1], n.Scale[2]}}
		}
		if anim != nil {
			for _, ch := range anim.Channels {
				i := map[string]int{"translation": 0, "rotation": 1, "scale": 2}[ch.Path]
				local[ch.Node][i] = value(ch, t)
			}
		}
		out[k] = r.posed(local, b)
	}
	return out
}

// posed is where node b lies when the nodes' parts are local.
func (r *Rig) posed(local [][3][4]float64, b int) Vec {
	var chain []int
	for n := b; n >= 0; n = r.parent[n] {
		chain = append([]int{n}, chain...)
	}
	p := placed{rot: identity, size: 1}
	for _, n := range chain {
		l := local[n]
		p = placed{
			pos:  add(p.pos, scale(turn(p.rot, Vec{l[0][0], l[0][1], l[0][2]}), p.size)),
			rot:  mul(p.rot, Quat(l[1])),
			size: p.size * l[2][0],
		}
	}
	return p.pos
}

// value is a channel's value at t, held beyond its keys.
func value(ch gltf.Channel, t float64) [4]float64 {
	n := len(ch.Times)
	get := func(i int) [4]float64 { return floats(ch.Values[i]) }
	i := sort.Search(n, func(i int) bool { return float64(ch.Times[i]) > t })
	if i == 0 {
		return get(0)
	}
	if i == n {
		return get(n - 1)
	}
	t0, t1 := float64(ch.Times[i-1]), float64(ch.Times[i])
	f := (t - t0) / math.Max(t1-t0, 1e-9)
	a, b := get(i-1), get(i)
	if ch.Path == "rotation" {
		b = Quat(sameSide(Quat(b), Quat(a)))
	}
	var v [4]float64
	for k := range v {
		v[k] = a[k] + (b[k]-a[k])*f
	}
	if ch.Path == "rotation" {
		s := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2] + v[3]*v[3])
		for k := range v {
			v[k] /= s
		}
	}
	return v
}
