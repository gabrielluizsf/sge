package human

// browClearance is how far, in metres, an eyebrow stands off the skin: a
// hair has a thickness, and the body's facets are flat where the skin is
// not, so brows laid on it sink into the flesh without it.
const browClearance = 0.0016

// bodyNormals are the body's smooth normals at pos, one for each vertex:
// each face's, weighted by its area, summed on its corners.
func (m *mesh) bodyNormals(pos [][3]float64) [][3]float64 {
	out := make([][3]float64, len(pos))
	for _, f := range m.faces {
		if len(f.v) < 3 || m.groups[f.group] != "body" {
			continue
		}
		a := pos[f.v[0]]
		for i := 1; i+1 < len(f.v); i++ {
			n := cross(sub(pos[f.v[i]], a), sub(pos[f.v[i+1]], a))
			for _, v := range []int32{f.v[0], f.v[i], f.v[i+1]} {
				out[v] = add(out[v], n)
			}
		}
	}
	for i := range out {
		out[i] = normalize(out[i])
	}
	return out
}

// standOff moves the vertices at, placed by p, out along the skin until
// each is clear above it.
func (p *proxy) standOff(at, pos [][3]float64, normals [][3]float64) {
	for i, r := range p.refs {
		var on, n [3]float64
		for k := range 3 {
			on = add(on, scale(pos[r.v[k]], r.w[k]))
			n = add(n, scale(normals[r.v[k]], r.w[k]))
		}
		n = normalize(n)
		h := (at[i][0]-on[0])*n[0] + (at[i][1]-on[1])*n[1] + (at[i][2]-on[2])*n[2]
		if h < browClearance*metre {
			at[i] = add(at[i], scale(n, browClearance*metre-h))
		}
	}
}
