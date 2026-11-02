package human

import (
	"fmt"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// handleY is how far the handle's axis lies above a tool's origin, which
// is on the ground, in metres.
const handleY = 0.03

// handle is where the hands put the handle of a tool held as g says, in
// key k of animation a: where the tool's origin is, and how it is turned.
func handle(m *gltf.Model, a *gltf.Animation, g Grip, k int) (at [3]float64, turn quat) {
	w := worldAt(m, a, k)
	if g.Hands[0] == g.Hands[1] {
		n := boneOf(m, g.Hands[0])
		return add(w[n].at, w[n].turn.turn(g.Offset)), quat(g.Turn).then(w[n].turn)
	}
	var hand [2][3]float64
	for i, h := range g.Hands {
		n := boneOf(m, h)
		hand[i] = add(w[n].at, w[n].turn.turn(palm(m, n)))
	}
	along := flatten(normalize(sub(hand[1], hand[0])), g)
	turn = aimed(along, g.Facing)
	at = add(hand[0], scale(along, -g.End))
	at = sub(at, scale(turn.turn([3]float64{0, 1, 0}), handleY))
	return at, turn
}

// flatten is the handle's direction along, with as much of its turn toward
// g.Facing taken out as g.Flat says, more the further down it points.
func flatten(along [3]float64, g Grip) [3]float64 {
	down := min(max(-along[1]/0.5, 0), 1)
	w := g.Flat * down * down * (3 - 2*down)
	f := normalize(g.Facing)
	return normalize(sub(along, scale(f, w*dot3(along, f))))
}

// aimed is the turn that lays a tool's X axis along `along`, and its Z
// axis as near as it can come to facing.
func aimed(along, facing [3]float64) quat {
	q := between([3]float64{1, 0, 0}, along)
	z := q.turn([3]float64{0, 0, 1})
	perp := func(v [3]float64) [3]float64 { return normalize(sub(v, scale(along, dot3(v, along)))) }
	want := perp(facing)
	if dot3(facing, along)*dot3(facing, along) > 0.98 {
		want = perp([3]float64{0, 1, 0})
	}
	s := cross(perp(z), want)
	return q.then(around(along, math.Atan2(dot3(s, along), dot3(perp(z), want))))
}

func dot3(a, b [3]float64) float64 { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }

// place gives node `at`, the held tool, translation, rotation and scale
// channels in the animation called name: out of sight but for where the
// hands put it, in every key.
func place(m *gltf.Model, node int, name string, g Grip) error {
	for i := range m.Animations {
		a := &m.Animations[i]
		if a.Name != name || len(a.Channels) == 0 {
			continue
		}
		times := a.Channels[0].Times
		chans := [3]gltf.Channel{{Path: "translation"}, {Path: "rotation"}, {Path: "scale"}}
		for k := range times {
			at, q := handle(m, a, g, k)
			if n := len(chans[1].Values); n > 0 && dot4(chans[1].Values[n-1], q) < 0 {
				q = quat{-q[0], -q[1], -q[2], -q[3]}
			}
			chans[0].Values = append(chans[0].Values, [4]float32{float32(at[0]), float32(at[1]), float32(at[2])})
			chans[1].Values = append(chans[1].Values, [4]float32{float32(q[0]), float32(q[1]), float32(q[2]), float32(q[3])})
			chans[2].Values = append(chans[2].Values, [4]float32{1, 1, 1})
		}
		var added []gltf.Channel
		for _, ch := range chans {
			ch.Node, ch.Interpolation, ch.Times = node, "LINEAR", times
			added = append(added, ch)
		}
		a.Channels = append(a.Channels, added...)
		return nil
	}
	return fmt.Errorf("human: no animation %s to hold in", name)
}
