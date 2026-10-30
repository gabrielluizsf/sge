package render3d

import "github.com/gabrielluizsf/sge/gltf"

// fallback is glTF's default material: white, rough, opaque.
var fallback = gltf.Material{BaseColor: [4]float64{1, 1, 1, 1}, Texture: -1, MetalRough: -1, Normal: -1, Occlusion: -1,
	Roughness: 1, AlphaMode: "OPAQUE", AlphaCutoff: 0.5}

// uploadMaterials makes m's materials' uniforms, after the default one,
// each at its own aligned place in one buffer, and their textures.
func (r *Renderer) uploadMaterials(m *Model) error {
	all := append([]gltf.Material{fallback}, m.src.Materials...)
	data := make([]float32, len(all)*UniformAlign)
	m.textures = make([]Texture, len(m.src.Images))
	for i, g := range all {
		at := i * UniformAlign
		cutoff := 0.0
		switch g.AlphaMode {
		case "MASK":
			cutoff = g.AlphaCutoff
		case "BLEND":
			cutoff = 0.5 // drawn cut, not blended
		}
		vals := [MaterialFloats]float64{g.BaseColor[0], g.BaseColor[1], g.BaseColor[2], g.BaseColor[3], g.Metallic, g.Roughness, cutoff}
		for k, v := range vals {
			data[at+k] = float32(v)
		}
		tex, err := r.textureOf(m, g.Texture)
		if err != nil {
			return err
		}
		r.materials++
		m.materials = append(m.materials, material{at: at, texture: tex, cut: cutoff > 0, double: g.DoubleSided, id: r.materials})
	}
	var err error
	m.uniforms, err = r.dev.NewBuffer(BufferDesc{Kind: Uniforms, Floats: data})
	return err
}

// textureOf is m's image img as an sRGB texture, made once, or 0 for none.
func (r *Renderer) textureOf(m *Model, img int) (Texture, error) {
	if img < 0 || img >= len(m.src.Images) || m.src.Images[img] == nil {
		return 0, nil
	}
	if m.textures[img] != 0 {
		return m.textures[img], nil
	}
	w, h, px := rgba8(m.src.Images[img])
	t, err := r.dev.NewTexture(TextureDesc{Width: w, Height: h, Format: RGBA8, Pixels: px, SRGB: true, Mipmaps: true})
	m.textures[img] = t
	return t, err
}
