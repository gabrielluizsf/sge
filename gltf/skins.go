package gltf

// skins are the model's skins, their joints checked to be nodes and their
// inverse bind matrices, one a joint.
func skins(doc *document, v views) ([]Skin, error) {
	var out []Skin
	for i, ds := range doc.Skins {
		s := Skin{Joints: ds.Joints}
		for _, j := range ds.Joints {
			if j < 0 || j >= len(doc.Nodes) {
				return nil, malformed("skin %d's joint %d", i, j)
			}
		}
		s.InverseBind = make([][16]float64, len(ds.Joints))
		for j := range s.InverseBind {
			s.InverseBind[j] = [16]float64{0: 1, 5: 1, 10: 1, 15: 1}
		}
		if ds.InverseBindMatrices != nil {
			a, err := v.accessor(*ds.InverseBindMatrices, "MAT4")
			if err != nil {
				return nil, err
			}
			if a.count != len(ds.Joints) {
				return nil, malformed("skin %d has %d matrices for %d joints", i, a.count, len(ds.Joints))
			}
			for j := range s.InverseBind {
				for k := range 16 {
					s.InverseBind[j][k] = a.float(j, k)
				}
			}
		}
		out = append(out, s)
	}
	return out, nil
}
