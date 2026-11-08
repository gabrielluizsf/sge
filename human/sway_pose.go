package human

import (
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// headTrack is where the head is, and which way it turns, as an animation
// moves the body: found by walking the bones from the root to the head.
type headTrack struct {
	m     *gltf.Model
	an    *gltf.Animation
	path  []int
	chans map[int][2]*gltf.Channel // by node: its translation and rotation
}

// newHeadTrack is the head's track for animation an of m, nil if m has
// no head.
func newHeadTrack(m *gltf.Model, an *gltf.Animation) *headTrack {
	parent := map[int]int{}
	head := -1
	for i, n := range m.Nodes {
		for _, c := range n.Children {
			parent[c] = i
		}
		if n.Name == "head" {
			head = i
		}
	}
	if head < 0 {
		return nil
	}
	h := &headTrack{m: m, an: an, chans: map[int][2]*gltf.Channel{}}
	for n := head; ; {
		h.path = append([]int{n}, h.path...)
		p, ok := parent[n]
		if !ok {
			break
		}
		n = p
	}
	for i := range an.Channels {
		c := &an.Channels[i]
		pair := h.chans[c.Node]
		switch c.Path {
		case "translation":
			pair[0] = c
		case "rotation":
			pair[1] = c
		}
		h.chans[c.Node] = pair
	}
	return h
}

// at is the head's place and turn, in the model's space, at time t.
func (h *headTrack) at(t float64) ([3]float64, quat) {
	pos, rot := [3]float64{}, still
	for _, n := range h.path {
		node := h.m.Nodes[n]
		local, turn := node.Translation, quat(node.Rotation)
		pair := h.chans[n]
		if pair[0] != nil {
			v := keyAt(pair[0], t)
			local = [3]float64{v[0], v[1], v[2]}
		}
		if pair[1] != nil {
			turn = unit(keyAt(pair[1], t))
		}
		pos = add(pos, rot.turn(local))
		rot = turn.then(rot)
	}
	return pos, rot
}

// keyAt is channel c's value at time t, between its keys.
func keyAt(c *gltf.Channel, t float64) quat {
	n := len(c.Times)
	val := func(i int) quat {
		v := c.Values[i]
		return quat{float64(v[0]), float64(v[1]), float64(v[2]), float64(v[3])}
	}
	if n == 1 || t <= float64(c.Times[0]) {
		return val(0)
	}
	for i := 1; i < n; i++ {
		t1 := float64(c.Times[i])
		if t <= t1 {
			t0 := float64(c.Times[i-1])
			f := (t - t0) / max(t1-t0, 1e-9)
			a, b := val(i-1), val(i)
			var out quat
			for k := range out {
				out[k] = a[k] + (b[k]-a[k])*f
			}
			return out
		}
	}
	return val(n - 1)
}

func unit(q quat) quat {
	l := q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3]
	if l == 0 {
		return still
	}
	l = 1 / math.Sqrt(l)
	return quat{q[0] * l, q[1] * l, q[2] * l, q[3] * l}
}
