package human

import "image"

// track notes material mat, as it is before it is tinted, to tint again
// later by role r; it must be called before recolour and greyOut.
func (d *dresser) track(mat int, r role, greys bool) {
	if r.kind == roleNone {
		return
	}
	m := d.model.Materials[mat]
	rt := retint{open: true, mat: mat, role: r, grey: greys, orig: -1, own: -1, color: m.BaseColor}
	if m.Texture >= 0 {
		if img, ok := d.model.Images[m.Texture].(*image.NRGBA); ok {
			rt.src, rt.orig = img, m.Texture
		}
	}
	d.retints = append(d.retints, rt)
}

// settle notes, for each tracked material, what recolour and greyOut made
// of it: the tint it wears now, and the texture of its own, if any.
func (d *dresser) settle() {
	for i := range d.retints {
		r := &d.retints[i]
		if !r.open {
			continue
		}
		r.open = false
		if m := d.model.Materials[r.mat]; r.src != nil && m.Texture != r.orig {
			r.own = m.Texture
		}
		r.wantSrc = r.src
		t, ok := tintFor(r.role, d.look)
		r.want = keyOf(t, ok)
		if r.grey {
			r.want.greying = d.grey
		}
		if r.role.kind == roleEyes {
			r.want.eyes = d.look.Eyes
		}
	}
}
