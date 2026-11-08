package gltf

// materials are the model's materials, their textures checked to be
// images it has.
func materials(doc *document, nImages int) ([]Material, error) {
	var out []Material
	for i, dm := range doc.Materials {
		m := Material{Name: dm.Name, BaseColor: [4]float64{1, 1, 1, 1}, Texture: -1, MetalRough: -1, Normal: -1, Occlusion: -1,
			NormalScale: 1, OcclusionStrength: 1, Metallic: 1, Roughness: 1, AlphaMode: "OPAQUE", AlphaCutoff: 0.5, DoubleSided: dm.DoubleSided}
		if len(dm.PBR.BaseColorFactor) == 4 {
			copy(m.BaseColor[:], dm.PBR.BaseColorFactor)
		}
		image := func(texture int) (int, error) {
			if texture < 0 || texture >= len(doc.Textures) || doc.Textures[texture].Source == nil {
				return -1, malformed("material %d's texture %d", i, texture)
			}
			src := *doc.Textures[texture].Source
			if src < 0 || src >= nImages {
				return -1, malformed("material %d's image %d", i, src)
			}
			return src, nil
		}
		var err error
		if t := dm.PBR.BaseColorTexture; t != nil {
			if m.Texture, err = image(t.Index); err != nil {
				return nil, err
			}
		}
		if t := dm.PBR.MetallicRoughnessTexture; t != nil {
			if m.MetalRough, err = image(t.Index); err != nil {
				return nil, err
			}
		}
		if t := dm.NormalTexture; t != nil {
			if m.Normal, err = image(t.Index); err != nil {
				return nil, err
			}
			if t.Scale != nil {
				m.NormalScale = *t.Scale
			}
		}
		if t := dm.OcclusionTexture; t != nil {
			if m.Occlusion, err = image(t.Index); err != nil {
				return nil, err
			}
			if t.Strength != nil {
				m.OcclusionStrength = *t.Strength
			}
		}
		if f := dm.PBR.MetallicFactor; f != nil {
			m.Metallic = *f
		}
		if f := dm.PBR.RoughnessFactor; f != nil {
			m.Roughness = *f
		}
		switch dm.AlphaMode {
		case "", "OPAQUE":
		case "MASK", "BLEND":
			m.AlphaMode = dm.AlphaMode
		default:
			return nil, malformed("material %d's alpha mode %q", i, dm.AlphaMode)
		}
		if dm.AlphaCutoff != nil {
			m.AlphaCutoff = *dm.AlphaCutoff
		}
		out = append(out, m)
	}
	return out, nil
}
