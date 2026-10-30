package render3d

import "github.com/gabrielluizsf/sge/pose"

// group lays out the draws of instances of one model, nearest first:
// each rigid part drawn once for all of them, MaxObjects at a time, its
// matrices side by side; each skinned part once an instance, with its
// skin's palette.
func (r *Renderer) group(g []seen) {
	m := g[0].model
	r.partAt = fit(r.partAt, len(m.parts))
	for pi, p := range m.parts {
		if p.skin >= 0 {
			continue
		}
		at := r.reserve(len(g) * ObjectFloats)
		r.partAt[pi] = at
		for c := 0; c < len(g); c += MaxObjects {
			n := min(MaxObjects, len(g)-c)
			r.add(drawItem{model: m, part: pi, objects: at + c*ObjectFloats, count: n, instances: n, depth: g[c].depth})
		}
	}
	for i, v := range g {
		r.pose(v, i)
	}
}

// pose poses the i-th instance of its group, v, and writes its rigid
// parts' matrices and its skins' palettes.
func (r *Renderer) pose(v seen, i int) {
	m, in := v.model, v.in
	anim := in.animation(m)
	var worlds []pose.Mat
	if anim >= 0 || m.skins > 0 {
		worlds = m.poser.At(anim, in.Time, in.Place.Affine(up))
	}
	r.skinAt = fit(r.skinAt, m.skins)
	for k := range r.skinAt {
		r.skinAt[k] = -1
	}
	for pi, p := range m.parts {
		if p.skin < 0 {
			var w pose.Mat
			if worlds != nil {
				w = worlds[p.node]
			} else {
				w = in.Place.Affine(m.rest[p.node])
			}
			r.put(r.partAt[pi]+i*ObjectFloats, &w)
			continue
		}
		joints := len(m.src.Skins[p.skin].Joints)
		at := r.skinAt[p.skin]
		if at < 0 {
			m.joints = m.poser.Joints(p.skin, m.joints)
			at = r.reserve(joints * ObjectFloats)
			for j := range m.joints {
				r.put(at+j*ObjectFloats, &m.joints[j])
			}
			r.skinAt[p.skin] = at
		}
		r.add(drawItem{model: m, part: pi, objects: at, count: joints, instances: 1, depth: v.depth})
	}
}

// fit is s with room for n, reusing its memory.
func fit(s []int, n int) []int {
	if cap(s) < n {
		return make([]int, n)
	}
	return s[:n]
}
