package art

import (
	"fmt"
	"strings"
	"unicode"
)

// TextWidth is the width in pixels Text draws s at: six a character, a
// pixel of shadow and a pixel apart.
func TextWidth(s string) int {
	return len([]rune(s))*6 + 1
}

// Text is the SVG of s in the HUD's letters, in the palette colour ink,
// with a shadow a pixel down and right in the outline colour so it reads
// on any ground. Lower case is written upper; a character the font lacks
// is left blank.
func Text(s string, ink rune) string {
	w, h := TextWidth(s), 8
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" shape-rendering="crispEdges">`+"\n", w, h, w, h)
	for _, pass := range []struct {
		dx, dy int
		c      string
	}{{1, 1, Palette['k']}, {0, 0, Palette[ink]}} {
		for i, r := range []rune(s) {
			if r != 'x' {
				r = unicode.ToUpper(r)
			}
			g, ok := font[r]
			if !ok {
				continue
			}
			for y, row := range g {
				for x, c := range row {
					if c == '#' {
						fmt.Fprintf(&b, `<rect x="%d" y="%d" width="1" height="1" fill="%s"/>`+"\n", i*6+x+pass.dx, y+pass.dy, pass.c)
					}
				}
			}
		}
	}
	b.WriteString("</svg>\n")
	return b.String()
}
