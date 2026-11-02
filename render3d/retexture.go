package render3d

import "fmt"

// Retexture makes material i of m (an index of its glTF model's
// Materials) draw as the model now says: its base colour, and its colour
// texture, whose pixels may have changed in place. Only that material is
// touched, its texture made again on the device, for a model whose
// colours change while the mesh does not.
func (r *Renderer) Retexture(m *Model, i int) error {
	at := i + 1 // the default material comes first
	if i < 0 || at >= len(m.materials) {
		return fmt.Errorf("render3d: no material %d", i)
	}
	g := m.src.Materials[i]
	mt := &m.materials[at]
	var tex Texture
	if g.Texture >= 0 && g.Texture < len(m.src.Images) && m.src.Images[g.Texture] != nil {
		w, h, px := rgba8(m.src.Images[g.Texture])
		var err error
		if tex, err = r.dev.NewTexture(TextureDesc{Width: w, Height: h, Format: RGBA8, Pixels: px, SRGB: true, Mipmaps: true}); err != nil {
			return err
		}
	}
	if mt.owned != 0 {
		r.dev.Release(mt.owned)
	}
	mt.texture, mt.owned = tex, tex
	r.dev.Write(m.uniforms, mt.at, []float32{float32(g.BaseColor[0]), float32(g.BaseColor[1]), float32(g.BaseColor[2]), float32(g.BaseColor[3])})
	return nil
}
