package tree

// card is a flat quad of foliage, flowers or needles: its foot at at, its
// texture's bottom middle, running along for length, width wide across;
// of material mat. It is lit as facing normal, which, for a cluster in a
// crown, is outward from the crown, as a leafy mass is lit, or zero for
// the card's own face.
type card struct {
	at, along, across vec
	length, width     float64
	mat               int
	normal            vec
}

// card adds c.
func (g *grower) card(c card) {
	n := c.normal
	if n == (vec{}) {
		n = normalize(cross(c.across, c.along))
	}
	n = normalize(n)
	p := g.prim(c.mat)
	first := uint32(len(p.Positions))
	tip := add(c.at, scale(c.along, c.length))
	half := scale(c.across, c.width/2)
	for _, q := range []struct {
		at vec
		uv [2]float32
	}{{sub(c.at, half), [2]float32{0, 1}}, {add(c.at, half), [2]float32{1, 1}}, {sub(tip, half), [2]float32{0, 0}}, {add(tip, half), [2]float32{1, 0}}} {
		p.Positions = append(p.Positions, f32(q.at))
		p.Normals = append(p.Normals, f32(n))
		p.UVs = append(p.UVs, q.uv)
	}
	p.Indices = append(p.Indices, first, first+1, first+3, first, first+3, first+2)
}

// ribbon is a long card bent along path, as a palm frond arches: width
// wide, its two edges folded down fold of its width, as a frond's
// leaflets hang; of material mat.
type ribbon struct {
	path  []vec
	width float64
	fold  float64
	mat   int
}

// ribbon adds r, its texture's bottom at its first point.
func (g *grower) ribbon(r ribbon) {
	p := g.prim(r.mat)
	first := uint32(len(p.Positions))
	n := len(r.path)
	for i, at := range r.path {
		_, t := along(r.path, float64(i)/float64(n-1))
		side := cross(t, up)
		if length(side) < 0.1 {
			side = perp(t)
		}
		side = normalize(side)
		face := normalize(cross(side, t))
		drop := vec{0, -r.fold * r.width / 2, 0}
		v := float32(1 - float64(i)/float64(n-1))
		for k, e := range []vec{add(sub(at, scale(side, r.width/2)), drop), at, add(add(at, scale(side, r.width/2)), drop)} {
			p.Positions = append(p.Positions, f32(e))
			p.Normals = append(p.Normals, f32(face))
			p.UVs = append(p.UVs, [2]float32{float32(k) / 2, v})
		}
		if i > 0 {
			b := first + uint32(3*(i-1))
			p.Indices = append(p.Indices, b, b+1, b+4, b, b+4, b+3, b+1, b+2, b+5, b+1, b+5, b+4)
		}
	}
}
