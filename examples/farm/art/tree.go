package art

// treeTile is a tree standing on its cell, its crown leant lean units in
// the wind: a round, many-leafed crown, lit from the top left, on a trunk
// that stays put, over its shadow.
func treeTile(lean float64) string {
	d := newDoc(IsoW, IsoH)
	shadow(d, 9, 4)
	bark := d.linear(grad{0, 0, 1, 0, []stop{{0, barkLight}, {0.55, barkBase}, {1, barkShade}}})
	d.add(`<path d="M14.2 41 C14.8 36 14.6 30 15 24 L17 24 C17.4 30 17.2 36 17.8 41 C16.6 41.8 15.4 41.8 14.2 41 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, bark, line)
	d.add(`<path d="M15.4 30 q-2 -2 -3.2 -2.6 M16.8 27 q2 -1.4 3 -1.6" stroke="%s" stroke-width="0.9" fill="none" stroke-linecap="round"/>`, barkBase)
	d.add(`<g transform="translate(%g 0)">`, lean)
	leaves := d.linear(grad{0.2, 0, 0.8, 1, []stop{{0, leafLight}, {0.5, leafBase}, {1, leafShade}}})
	for _, c := range [][3]float64{{10, 20, 5.6}, {22, 20, 5.4}, {16, 13, 6.4}, {11.5, 13.5, 4.6}, {20.5, 13, 4.8}, {16, 21.5, 6}} {
		d.add(`<circle cx="%g" cy="%g" r="%g" fill="%s" stroke="%s" stroke-width="0.6"/>`, c[0], c[1], c[2], leaves, line)
	}
	// Fill the seams between the leaves' outlines, then light and shade.
	d.add(`<ellipse cx="16" cy="17.5" rx="8.4" ry="6.6" fill="%s"/>`, leafBase)
	d.add(`<path d="M8 22 Q16 28 24.5 21.5 Q20 25.5 16 25.6 Q11 25.4 8 22 Z" fill="%s" opacity="0.8"/>`, leafShade)
	for _, g := range [][3]float64{{12.5, 11.5, 2.4}, {9, 17.5, 1.6}, {17, 9.5, 1.5}, {19.5, 15.5, 1.2}} {
		d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s" opacity="0.9"/>`, g[0], g[1], g[2], g[2]*0.7, leafGlow)
	}
	d.add(`</g>`)
	return d.String()
}

// shadow is a soft shadow on the ground at the cell's middle, rx by ry.
func shadow(d *doc, rx, ry float64) {
	x, y := iso(Size/2, Size/2, 0)
	d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s" opacity="0.35"/>`, x, y, rx, ry, shadowInk)
}
