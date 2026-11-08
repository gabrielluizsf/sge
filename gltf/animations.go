package gltf

// animations are the model's animations: each channel's node and path
// checked, its keyframe times and values read, as many values as times,
// three a key for translation and scale and four for rotation.
func animations(doc *document, v views) ([]Animation, error) {
	var out []Animation
	for ai, da := range doc.Animations {
		an := Animation{Name: da.Name}
		for ci, dc := range da.Channels {
			if dc.Target.Node == nil || dc.Target.Path == "weights" {
				continue // morph targets are not read
			}
			ch, err := channel(doc, v, animChannel{da, ci})
			if err != nil {
				return nil, malformed("animation %d channel %d: %v", ai, ci, err)
			}
			an.Channels = append(an.Channels, ch)
		}
		out = append(out, an)
	}
	return out, nil
}

// animChannel is an animation's channel ci.
type animChannel struct {
	da docAnimation
	ci int
}

// channel reads an animation's channel: its node and path checked, its
// sampler's times and values.
func channel(doc *document, v views, at animChannel) (Channel, error) {
	da := at.da
	dc := da.Channels[at.ci]
	node := *dc.Target.Node
	if node < 0 || node >= len(doc.Nodes) {
		return Channel{}, malformed("node %d", node)
	}
	kind := map[string]string{"translation": "VEC3", "scale": "VEC3", "rotation": "VEC4"}[dc.Target.Path]
	if kind == "" {
		return Channel{}, malformed("path %q", dc.Target.Path)
	}
	if dc.Sampler < 0 || dc.Sampler >= len(da.Samplers) {
		return Channel{}, malformed("sampler %d", dc.Sampler)
	}
	s := da.Samplers[dc.Sampler]
	interp := s.Interpolation
	switch interp {
	case "":
		interp = "LINEAR"
	case "LINEAR", "STEP":
	default:
		return Channel{}, malformed("interpolation %q", interp)
	}
	in, err := v.accessor(s.Input, "SCALAR")
	if err != nil {
		return Channel{}, err
	}
	outA, err := v.accessor(s.Output, kind)
	if err != nil {
		return Channel{}, err
	}
	if outA.count != in.count {
		return Channel{}, malformed("%d values for %d times", outA.count, in.count)
	}
	ch := Channel{Node: node, Path: dc.Target.Path, Interpolation: interp, Times: make([]float32, in.count), Values: make([][4]float32, in.count)}
	for e := range in.count {
		ch.Times[e] = float32(in.float(e, 0))
		for c := range outA.comps {
			ch.Values[e][c] = float32(outA.float(e, c))
		}
	}
	return ch, nil
}
