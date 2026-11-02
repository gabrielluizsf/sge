package human

import "image"

// recolour tints material mat: its colour texture, as a copy of its own,
// for others who wear the same keep theirs, or its colour, if it has no
// texture.
func (d *dresser) recolour(mat int, t tint) {
	m := &d.model.Materials[mat]
	if m.Texture < 0 {
		m.BaseColor = [4]float64{t.to[0], t.to[1], t.to[2], m.BaseColor[3]}
		return
	}
	src, ok := d.model.Images[m.Texture].(*image.NRGBA)
	if !ok {
		return
	}
	d.model.Images = append(d.model.Images, tinted(src, t))
	m.Texture = len(d.model.Images) - 1
}

// tintOf is the tint of what the look wears at i, the Wear name: hair
// and brows in the hair's colour, eyes in the eyes', clothes in their
// own; false if it is left as it is.
func (d *dresser) tintOf(i int) (tint, bool) {
	return tintFor(roleOf(d.look.Wear[i], i), d.look)
}

// tintTo is the tint to colour c, over mask, false if c asks for none.
func tintTo(c Colour, mask func(x, y float64) float64) (tint, bool) {
	return tint{to: c.linear(), mask: mask}, c.Set
}
