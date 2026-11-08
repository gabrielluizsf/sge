package animal

import (
	"fmt"

	"github.com/gabrielluizsf/sge/gltf"
)

// Rig is a model and the bones of its first skin, with where each rests
// in the scene.
type Rig struct {
	M      *gltf.Model
	rest   []placed
	parent []int
	joint  []bool
	byName map[string]int
}

// NewRig is the rig of m, whose skin 0 holds the bones.
func NewRig(m *gltf.Model) (*Rig, error) {
	if len(m.Skins) == 0 {
		return nil, fmt.Errorf("animal: the model has no skin")
	}
	r := &Rig{M: m, joint: make([]bool, len(m.Nodes)), byName: map[string]int{}}
	r.rest, r.parent = worlds(m)
	for _, j := range m.Skins[0].Joints {
		r.joint[j] = true
		r.byName[m.Nodes[j].Name] = j
	}
	return r, nil
}

// Bone is the node of the bone called name; it panics for a name the model
// does not have, a mistake in the code that names it.
func (r *Rig) Bone(name string) int {
	n, ok := r.byName[name]
	if !ok {
		panic("animal: no bone " + name)
	}
	return n
}

// Head is where bone b rests in the scene.
func (r *Rig) Head(b int) Vec { return r.rest[b].pos }

// Root is the first bone with no bone above it.
func (r *Rig) Root() int {
	for _, j := range r.M.Skins[0].Joints {
		if p := r.parent[j]; p < 0 || !r.joint[p] {
			return j
		}
	}
	return r.M.Skins[0].Joints[0]
}

// Toward is the direction from bone a's head to bone b's, along the ground.
func (r *Rig) Toward(a, b int) Vec {
	d := sub(r.Head(b), r.Head(a))
	d[1] = 0
	return d
}
