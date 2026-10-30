package gltf

// materials are the model's materials, their textures checked to be
// images it has.
func materials(doc *document, nImages int) ([]Material, error) {
	var out []Material
	for i, dm := range doc.Materials {
		m := Material{Name: dm.Name, BaseColor: [4]float64{1, 1, 1, 1}, Texture: -1, Metallic: 1, Roughness: 1,
			AlphaMode: "OPAQUE", AlphaCutoff: 0.5, DoubleSided: dm.DoubleSided}
		if len(dm.PBR.BaseColorFactor) == 4 {
			copy(m.BaseColor[:], dm.PBR.BaseColorFactor)
		}
		if t := dm.PBR.BaseColorTexture; t != nil {
			if t.Index < 0 || t.Index >= len(doc.Textures) || doc.Textures[t.Index].Source == nil {
				return nil, malformed("material %d's texture %d", i, t.Index)
			}
			src := *doc.Textures[t.Index].Source
			if src < 0 || src >= nImages {
				return nil, malformed("material %d's image %d", i, src)
			}
			m.Texture = src
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
