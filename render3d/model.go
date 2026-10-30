package render3d

import (
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/pose"
)

// Model is a glTF model uploaded to a device: its vertices, indices,
// materials and textures there, and what poses it here.
type Model struct {
	// Bounds is a ball round the model, in the game's axes, wide enough
	// for its animations to stay inside.
	Bounds Sphere

	src       *gltf.Model
	poser     *pose.Poser
	rest      []pose.Mat // the nodes in the model's rest pose
	animated  bool
	parts     []part
	materials []material
	skins     int
	id        int // its number among the renderer's, for grouping

	static, skinned, indices, uniforms Buffer
	textures                           []Texture
	joints                             []pose.Mat // scratch for a pose
}

// part is a primitive uploaded: its node, its skin (-1 for a rigid one),
// where its indices are, and its material, an index of materials.
type part struct {
	node, skin   int
	first, count int
	material     int
}

// material is a material uploaded: where its uniforms are, its base
// colour texture (0 for none), whether it cuts by alpha and is seen from
// both sides, and its number among all the renderer's, for sorting.
type material struct {
	at          int
	texture     Texture
	cut, double bool
	id          int
}

// up turns glTF's Y up onto the game's Z, and glTF's Z, toward a model's
// viewer, onto the game's -Y: a rotation, so nothing is mirrored.
var up = pose.Mat{1, 0, 0, 0, 0, 0, 1, 0, 0, -1, 0, 0, 0, 0, 0, 1}

// Animations are the model's animations' names.
func (m *Model) Animations() []string {
	var out []string
	for _, a := range m.src.Animations {
		out = append(out, a.Name)
	}
	return out
}

// Release lets go of the model's buffers and textures on dev.
func (m *Model) Release(dev Device) {
	for _, b := range []Buffer{m.static, m.skinned, m.indices, m.uniforms} {
		if b != 0 {
			dev.Release(b)
		}
	}
	for _, t := range m.textures {
		if t != 0 {
			dev.Release(t)
		}
	}
}
