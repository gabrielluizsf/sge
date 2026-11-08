package animal

import (
	"fmt"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// Bone is a bone to give a model: its name and its parent's ("" for the
// root), where it begins and where it ends, in the model's own axes. It
// moves the skin nearest to the line between them.
type Bone struct {
	Name, Parent string
	Head, Tail   Vec
}

// Skeleton gives m, whose meshes have no skin, the bones as skin 0, each
// vertex bound to the nearest bones, and returns it as a rig. The bones
// rest unturned, so the model looks as it did.
func Skeleton(m *gltf.Model, bones []Bone) (*Rig, error) {
	at := map[string]int{}
	var sk gltf.Skin
	for i, b := range bones {
		at[b.Name] = i
		node := len(m.Nodes)
		parent := Vec{}
		if b.Parent != "" {
			j, ok := at[b.Parent]
			if !ok {
				return nil, fmt.Errorf("animal: bone %s comes before its parent %s", b.Name, b.Parent)
			}
			parent = bones[j].Head
			up := sk.Joints[j]
			m.Nodes[up].Children = append(m.Nodes[up].Children, node)
		} else {
			m.Roots = append(m.Roots, node)
		}
		m.Nodes = append(m.Nodes, gltf.Node{
			Name: b.Name, Mesh: -1, Skin: -1, Translation: sub(b.Head, parent),
			Rotation: identity, Scale: [3]float64{1, 1, 1},
		})
		sk.Joints = append(sk.Joints, node)
		sk.InverseBind = append(sk.InverseBind, [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, -b.Head[0], -b.Head[1], -b.Head[2], 1})
	}
	m.Skins = append(m.Skins, sk)
	each(m, func(n int, p *gltf.Primitive) {
		m.Nodes[n].Skin = len(m.Skins) - 1
		bind(p, bones)
	})
	return NewRig(m)
}

// bind gives each vertex of p the three nearest bones, nearer ones weighing
// more.
func bind(p *gltf.Primitive, bones []Bone) {
	p.Joints = make([][4]uint16, len(p.Positions))
	p.Weights = make([][4]float32, len(p.Positions))
	for i, v := range p.Positions {
		q := Vec{float64(v[0]), float64(v[1]), float64(v[2])}
		var best [3]int
		var near [3]float64
		for k := range near {
			near[k] = math.Inf(1)
		}
		for j, b := range bones {
			d := segment(q, b.Head, b.Tail)
			for k := range near {
				if d < near[k] {
					copy(near[k+1:], near[k:len(near)-1])
					copy(best[k+1:], best[k:len(best)-1])
					near[k], best[k] = d, j
					break
				}
			}
		}
		var sum float64
		var w [3]float64
		for k := range w {
			w[k] = 1 / math.Pow(near[k]+1e-3, 4)
			sum += w[k]
		}
		for k := range w {
			p.Joints[i][k] = uint16(best[k])
			p.Weights[i][k] = float32(w[k] / sum)
		}
	}
}

// segment is how far q is from the line from a to b.
func segment(q, a, b Vec) float64 {
	ab, aq := sub(b, a), sub(q, a)
	l := ab[0]*ab[0] + ab[1]*ab[1] + ab[2]*ab[2]
	t := 0.0
	if l > 0 {
		t = min(max((aq[0]*ab[0]+aq[1]*ab[1]+aq[2]*ab[2])/l, 0), 1)
	}
	return dist(q, add(a, scale(ab, t)))
}
