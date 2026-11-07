package human

import "github.com/gabrielluizsf/sge/gltf"

// swayHair adds to animation an the swing of the hair of m: a rotation
// channel for its sway bone, keyed as the pose channels are, keys of them,
// evenly, as the hold code reads every channel by pose. The baked frames
// and a live model move the hair alike. A model with no hair is left as
// it is.
func swayHair(m *gltf.Model, an *gltf.Animation, keys int, once bool) {
	node, ok := swayNode(m)
	if !ok {
		return
	}
	h := newHeadTrack(m, an)
	if h == nil {
		return
	}
	length := float64(animationLength(an))
	lean := swayed(h, length, keys, once)
	ch := gltf.Channel{Node: node, Path: "rotation", Interpolation: "LINEAR"}
	var last quat
	for k, u := range lean {
		q := between([3]float64{0, -1, 0}, [3]float64{u[0], -1, u[1]})
		if k > 0 && dot4([4]float32{float32(last[0]), float32(last[1]), float32(last[2]), float32(last[3])}, q) < 0 {
			q = quat{-q[0], -q[1], -q[2], -q[3]}
		}
		last = q
		ch.Times = append(ch.Times, float32(length*float64(k)/float64(keys-1)))
		ch.Values = append(ch.Values, [4]float32{float32(q[0]), float32(q[1]), float32(q[2]), float32(q[3])})
	}
	if !once { // the loop closes on where it began
		ch.Values[len(ch.Values)-1] = ch.Values[0]
	}
	an.Channels = append(an.Channels, ch)
}

// animationLength is the time of the last key of an.
func animationLength(an *gltf.Animation) float32 {
	var l float32
	for _, c := range an.Channels {
		if n := len(c.Times); n > 0 {
			l = max(l, c.Times[n-1])
		}
	}
	return l
}
