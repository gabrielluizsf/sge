package art

import (
	"fmt"
	"strings"
)

// Panel is the SVG of a HUD panel w by h pixels: a fill with an outline,
// its corners cut a pixel, a light inner edge at the top and left and a
// shade at the bottom and right, as the sprites are lit. selected draws
// it in yellow, for the slot in use.
func Panel(w, h int, selected bool) string {
	fill, light, shade := Palette['c'], Palette['w'], Palette['b']
	if selected {
		fill, light, shade = Palette['y'], Palette['w'], Palette['Y']
	}
	ink := Palette['k']
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" shape-rendering="crispEdges">`+"\n", w, h, w, h)
	rect := func(x, y, rw, rh int, c string) {
		fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`+"\n", x, y, rw, rh, c)
	}
	rect(1, 0, w-2, h, ink)
	rect(0, 1, w, h-2, ink)
	rect(1, 1, w-2, h-2, fill)
	rect(1, 1, w-2, 1, light)
	rect(1, 1, 1, h-2, light)
	rect(1, h-2, w-2, 1, shade)
	rect(w-2, 1, 1, h-2, shade)
	b.WriteString("</svg>\n")
	return b.String()
}
