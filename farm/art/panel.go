package art

// Panel is the SVG of a HUD panel w by h pixels: a rounded card of warm
// parchment edged in wood, its top lit; selected draws it gold, for the
// slot or line in use.
func Panel(w, h int, selected bool) string {
	d := newDoc(float64(w), float64(h))
	light, dark, edge := "#fff8e8", "#f2ddb4", "#8a5a3b"
	if selected {
		light, dark, edge = "#fff1b8", "#f7c948", "#b0742a"
	}
	d.add(`<rect x="0.6" y="0.6" width="%g" height="%g" rx="3" fill="%s" stroke="%s" stroke-width="1.2"/>`, float64(w)-1.2, float64(h)-1.2, d.down(light, dark), edge)
	d.add(`<rect x="2.2" y="1.8" width="%g" height="1" rx="0.5" fill="%s" opacity="0.85"/>`, float64(w)-4.4, white)
	return d.String()
}
