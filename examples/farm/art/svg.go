package art

import (
	"fmt"
	"strings"
)

// SVG is g as an SVG drawing, as many pixels wide as its rows and tall as
// it has them: a rect for every run of one colour along a row, edges
// crisp, transparent where g is '.'.
func SVG(g Grid) string {
	var b strings.Builder
	w, h := len([]rune(g[0])), len(g)
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" shape-rendering="crispEdges">`+"\n", w, h, w, h)
	for y, row := range g {
		r := []rune(row)
		for x := 0; x < len(r); {
			c := r[x]
			n := 1
			for x+n < len(r) && r[x+n] == c {
				n++
			}
			if c != '.' {
				fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="1" fill="%s"/>`+"\n", x, y, n, Palette[c])
			}
			x += n
		}
	}
	b.WriteString("</svg>\n")
	return b.String()
}
