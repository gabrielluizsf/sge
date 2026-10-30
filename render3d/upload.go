package render3d

import (
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/pose"
)

// boundsMargin grows the bounds of a model with animations, whose limbs
// may reach past where they rest.
const boundsMargin = 1.25

// Upload puts src on the renderer's device, once, for its instances to
// draw: its vertices and indices, its materials, and their textures.
func (r *Renderer) Upload(src *gltf.Model) (*Model, error) {
	r.models++
	m := &Model{src: src, poser: pose.New(src), animated: len(src.Animations) > 0, skins: len(src.Skins), id: r.models}
	m.rest = append([]pose.Mat(nil), m.poser.At(-1, 0, up)...)
	var g geometry
	if err := g.gather(m); err != nil {
		return nil, err
	}
	margin := 1.0
	if m.animated {
		margin = boundsMargin
	}
	m.Bounds = g.bounds(margin)
	if err := r.uploadGeometry(m, &g); err != nil {
		m.Release(r.dev)
		return nil, err
	}
	if err := r.uploadMaterials(m); err != nil {
		m.Release(r.dev)
		return nil, err
	}
	return m, nil
}

// uploadGeometry makes m's vertex and index buffers of g's.
func (r *Renderer) uploadGeometry(m *Model, g *geometry) error {
	var err error
	if len(g.static) > 0 {
		if m.static, err = r.dev.NewBuffer(BufferDesc{Kind: Vertices, Floats: g.static}); err != nil {
			return err
		}
	}
	if len(g.skinned) > 0 {
		if m.skinned, err = r.dev.NewBuffer(BufferDesc{Kind: Vertices, Floats: g.skinned}); err != nil {
			return err
		}
	}
	if len(g.indices) > 0 {
		m.indices, err = r.dev.NewBuffer(BufferDesc{Kind: Indices, Indices: g.indices})
	}
	return err
}
