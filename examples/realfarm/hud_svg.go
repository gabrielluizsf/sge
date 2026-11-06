package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
)

// The HUD's own pictures, for what no model is of, in the colours of its
// panels: dark glass edged in gold, and a little light on what is lit.

// heartSVG is a heart of fondness: lit crimson, or an empty one, dark glass
// edged in gold.
func heartSVG(full bool) string {
	fill, edge, op := "url(#g)", "#4a0e12", "1"
	if !full {
		fill, edge, op = "#15110f", "#b8924a", "0.7"
	}
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
<defs><radialGradient id="g" cx="0.38" cy="0.3" r="0.8"><stop offset="0" stop-color="#ff8a78"/><stop offset="0.45" stop-color="#d23b3c"/><stop offset="1" stop-color="#7d1820"/></radialGradient></defs>
<path d="M8 13.6 C3 10 1.6 7.6 1.6 5.4 C1.6 3.4 3.1 2.2 4.8 2.2 C6.2 2.2 7.4 3 8 4.2 C8.6 3 9.8 2.2 11.2 2.2 C12.9 2.2 14.4 3.4 14.4 5.4 C14.4 7.6 13 10 8 13.6 Z" fill="%s" fill-opacity="%s" stroke="%s" stroke-width="0.9" stroke-linejoin="round"/>
</svg>`, fill, op, edge)
}

// markerSVG is the quiet marker over someone who wants something: a small
// dark glass disc, gold-edged, with a gold lozenge lit in it.
const markerSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
<defs><linearGradient id="m" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#ffe49a"/><stop offset="1" stop-color="#c08a2c"/></linearGradient></defs>
<circle cx="8" cy="8" r="6.6" fill="#15110f" fill-opacity="0.84" stroke="#f0c55e" stroke-width="0.8"/>
<path d="M8 3.8 L11.2 8 L8 12.2 L4.8 8 Z" fill="url(#m)" stroke="#7a5518" stroke-width="0.4"/>
</svg>`

// unbakedSVG is a faint ring, in an icon's place until it is baked.
const unbakedSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
<circle cx="8" cy="8" r="5" fill="none" stroke="#b8924a" stroke-opacity="0.5" stroke-width="0.8"/>
</svg>`

// cursorSVG is the outline of the cell being worked: its diamond in thin
// warm gold over a darker line, so it reads on any ground, barely filled.
const cursorSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 32" width="64" height="32">
<path d="M32 1.6 L62.4 16 L32 30.4 L1.6 16 Z" fill="#f0c55e" fill-opacity="0.07" stroke="#14100f" stroke-opacity="0.55" stroke-width="2.6" stroke-linejoin="round"/>
<path d="M32 1.6 L62.4 16 L32 30.4 L1.6 16 Z" fill="none" stroke="#f0c55e" stroke-opacity="0.95" stroke-width="1.1" stroke-linejoin="round"/>
</svg>`

// Cursor is the outline of a cell's diamond, as the farm marks the cell
// being worked.
func (a *realArt) Cursor() *canvas.Canvas {
	return a.hudPicture(drawing{"cursor", cursorSVG, farm.TileW, farm.Cell})
}

// foodSVG is the food bar's apple: lit crimson-gold, edged dark, like the
// heart beside it.
const foodSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
<defs><radialGradient id="a" cx="0.36" cy="0.32" r="0.8"><stop offset="0" stop-color="#f0a860"/><stop offset="0.5" stop-color="#b8402a"/><stop offset="1" stop-color="#6e1d16"/></radialGradient></defs>
<path d="M8 4.8 C5.2 3 2.4 5.2 3 8.6 C3.6 11.6 5.6 13.8 8 13 C10.4 13.8 12.4 11.6 13 8.6 C13.6 5.2 10.8 3 8 4.8 Z" fill="url(#a)" stroke="#3a0e0a" stroke-width="0.9" stroke-linejoin="round"/>
<path d="M8 4.8 Q8 3 9 2" fill="none" stroke="#3a2a14" stroke-width="0.9" stroke-linecap="round"/>
<path d="M9.2 3.6 Q11 2 12.6 3.2 Q11.2 4.6 9.2 3.6 Z" fill="#5a7a2c" stroke="#2a3a14" stroke-width="0.5"/>
</svg>`
