package art

// Panel is the SVG of a HUD panel w by h pixels: a rounded card of warm
// cream, outlined, its top lit; selected draws it gold, for the slot in
// use.
func Panel(w, h int, selected bool) string {
	d := newDoc(float64(w), float64(h))
	light, dark := shirtLight, "#e7d3a8"
	if selected {
		light, dark = "#fff3b0", "#f2c14e"
	}
	d.add(`<rect x="0.6" y="0.6" width="%g" height="%g" rx="2.6" fill="%s" stroke="%s" stroke-width="1.2"/>`, float64(w)-1.2, float64(h)-1.2, d.down(light, dark), woodShade)
	d.add(`<rect x="2" y="1.8" width="%g" height="1" rx="0.5" fill="%s" opacity="0.8"/>`, float64(w)-4, white)
	return d.String()
}
