package render3d

import (
	"cmp"
	"slices"
)

// drawItem is a draw of a frame: a part of a model, the Objects
// uniforms it reads, count matrices from objects on, drawn for so many
// instances, as deep as its nearest instance, and its sort key.
type drawItem struct {
	model                     *Model
	part                      int
	objects, count, instances int
	depth                     float64
	key                       uint64
}

// add adds d to the frame's draws, its key worked out: those cut by
// alpha after the opaque, then by pipeline, then by material.
func (r *Renderer) add(d drawItem) {
	p := d.model.parts[d.part]
	mt := d.model.materials[p.material]
	var key uint64
	if mt.cut {
		key |= 1 << 63
	}
	key |= uint64(kindOf(p, false)) << 58
	if mt.double {
		key |= 1 << 57
	}
	key |= uint64(uint32(mt.id)) << 16
	d.key = key
	r.draws = append(r.draws, d)
	r.Stats.Draws++
	r.Stats.Triangles += p.count / 3 * d.instances
}

// order sorts the frame's draws by key, and front to back within one.
func (r *Renderer) order() {
	slices.SortFunc(r.draws, func(a, b drawItem) int {
		if c := cmp.Compare(a.key, b.key); c != 0 {
			return c
		}
		return cmp.Compare(a.depth, b.depth)
	})
}

// kindOf is the pipeline kind drawing p, for its shadow or not.
func kindOf(p part, shadow bool) Kind {
	switch {
	case p.skin >= 0 && shadow:
		return ShadowSkinned
	case p.skin >= 0:
		return Skinned
	case shadow:
		return Shadow
	}
	return Static
}
