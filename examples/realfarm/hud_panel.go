package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/sprite"
)

// Panel is a HUD panel of w by h pixels of art: dark, a little
// translucent, edged in a thin warm gold, which the slot or line in use
// shows brighter and warmer.
func (a *realArt) Panel(w, h int, selected bool) *canvas.Canvas {
	key := fmt.Sprint("p", w, "x", h, selected)
	if cv, ok := a.hud[key]; ok {
		return cv
	}
	s, err := sprite.FromSVG(panelSVG(w, h, selected), w*farm.Scale, h*farm.Scale)
	if err != nil {
		panic(err)
	}
	a.hud[key] = s.Image
	return s.Image
}

// panelSVG is the SVG of a panel w by h.
func panelSVG(w, h int, selected bool) string {
	fill, edge, glow := "#15110f", "#b8924a", "0.0"
	if selected {
		fill, edge, glow = "#3b2a12", "#f0c55e", "0.22"
	}
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %[1]d %[2]d" width="%[1]d" height="%[2]d">
<rect x="0.5" y="0.5" width="%[3]g" height="%[4]g" rx="2.5" fill="%[5]s" opacity="0.84"/>
<rect x="1.6" y="1.6" width="%[6]g" height="%[7]g" rx="1.6" fill="#f0c55e" opacity="%[8]s"/>
<rect x="0.5" y="0.5" width="%[3]g" height="%[4]g" rx="2.5" fill="none" stroke="%[9]s" stroke-width="0.8"/>
</svg>
`, w, h, float64(w)-1, float64(h)-1, fill, float64(w)-3.2, float64(h)-3.2, glow, edge)
}
