package human

// surface is faces to make a part of: their texture places, where their
// vertices are, how bones move each, and their material.
type surface struct {
	faces []face
	uv    [][2]float64
	at    [][3]float64
	infl  func(int32) []influence
	mat   int
}

// split is a part of s's faces: each corner a vertex of its own wherever
// its position has two texture places, faces cut into triangles.
func split(s surface) part {
	faces, uv, at, infl := s.faces, s.uv, s.at, s.infl
	p := part{mat: s.mat}
	made := map[[2]int32]int32{}
	corner := func(f face, k int) int32 {
		key := [2]int32{f.v[k], -1}
		if len(f.t) > 0 {
			key[1] = f.t[k]
		}
		if i, ok := made[key]; ok {
			return i
		}
		i := int32(len(p.pos))
		made[key] = i
		p.pos = append(p.pos, at[f.v[k]])
		t := [2]float64{}
		if key[1] >= 0 {
			t = uv[key[1]]
		}
		p.uv = append(p.uv, [2]float64{t[0], 1 - t[1]})
		p.infl = append(p.infl, infl(f.v[k]))
		p.smooth = append(p.smooth, f.v[k])
		return i
	}
	for _, f := range faces {
		first := corner(f, 0)
		for k := 1; k+1 < len(f.v); k++ {
			p.tris = append(p.tris, first, corner(f, k), corner(f, k+1))
		}
	}
	return p
}

// normals are the part's vertices' normals, smooth: each the sum of its
// faces' areas' normals, one for all the vertices of one position.
func (p part) normals() [][3]float64 {
	by := map[int32][3]float64{}
	for t := 0; t+2 < len(p.tris); t += 3 {
		a, b, c := p.pos[p.tris[t]], p.pos[p.tris[t+1]], p.pos[p.tris[t+2]]
		n := cross(sub(b, a), sub(c, a))
		for k := range 3 {
			s := p.smooth[p.tris[t+k]]
			by[s] = add(by[s], n)
		}
	}
	out := make([][3]float64, len(p.pos))
	for i := range out {
		out[i] = normalize(by[p.smooth[i]])
	}
	return out
}
