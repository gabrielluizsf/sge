package human

import (
	"fmt"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// Expression is a face, as MakeHuman's expression mixer makes one: how
// far, from 0 to 1, each of its face pose units is taken, by name.
type Expression map[string]float64

// Lids is the face with both upper eyelids shut as far as closed, from 0,
// open, to 1, shut: a blink is taken through them.
func Lids(closed float64) Expression {
	return Expression{"LeftUpperLidClosed": closed, "RightUpperLidClosed": closed}
}

// Mouth is the face speaking: the jaw dropped as far as open, and the
// lips pursed as far as round, for an o or a u; both from 0 to 1.
func Mouth(open, round float64) Expression {
	return Expression{"JawDrop": open, "LipsKiss": round}
}

// Express is m, a person the kit made, with face e: the bones of the head
// turned as e's units, each part of the way its weight says, every node
// else as m has it. m itself is not changed.
func (k *Kit) Express(m *gltf.Model, e Expression) (*gltf.Model, error) {
	units, err := k.faceUnits()
	if err != nil {
		return nil, err
	}
	turns := map[string]quat{}
	for name, w := range e {
		unit, ok := units[name]
		if !ok {
			return nil, fmt.Errorf("human: no face pose unit %q", name)
		}
		w = max(0, min(1, w))
		for _, b := range unit {
			t, ok := turns[b.bone]
			if !ok {
				t = still
			}
			turns[b.bone] = t.then(partway(b.turn, w))
		}
	}
	c := *m
	c.Nodes = append([]gltf.Node{}, m.Nodes...)
	for i, n := range c.Nodes {
		if t, ok := turns[n.Name]; ok {
			c.Nodes[i].Rotation = [4]float64(quat(n.Rotation).then(t))
		}
	}
	return &c, nil
}

// partway is w of the way from no turn to q.
func partway(q quat, w float64) quat {
	half := math.Acos(max(-1, min(1, q[3])))
	s := math.Sin(half)
	if s < 1e-12 {
		return still
	}
	k := math.Sin(w*half) / s
	return quat{q[0] * k, q[1] * k, q[2] * k, math.Cos(w * half)}
}
