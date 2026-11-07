package art

import (
	"fmt"
	"math"
	"strings"
)

// scale is how many pixels a unit of the drawings is drawn as in the
// game: the ground's diamonds are cut to that grid, a row of pixels at a
// time, so their edges are whole pixels, wholly drawn or not, and tiles
// meet without a seam, as an outline smoothed across the edge would not.
const scale = 2

// diamondRows is the cell's diamond as rectangles, a row of pixels each:
// the pixels whose middles lie on it.
func diamondRows() string {
	var b strings.Builder
	for y := 2 * Size * scale; y < IsoH*scale; y++ {
		cy := (float64(y)+0.5)/scale - 2*Size
		half := Size - 2*abs(cy-Size/2)
		x0 := int(math.Ceil((Size-half)*scale - 0.5))
		x1 := int(math.Floor((Size+half)*scale - 0.5))
		if x1 >= x0 {
			fmt.Fprintf(&b, `<rect x="%g" y="%g" width="%g" height="%g"/>`, float64(x0)/scale, float64(y)/scale, float64(x1-x0+1)/scale, 1.0/scale)
		}
	}
	return b.String()
}

// groundTile is a tile of ground: its cell's diamond filled light at the
// far corner to shade at the near, and what paint adds drawn inside it,
// clipped to it.
func groundTile(light, shade string, paint func(d *doc)) string {
	d := newDoc(IsoW, IsoH)
	rows := diamondRows()
	fmt.Fprintf(&d.defs, `<clipPath id="cell">%s</clipPath>`+"\n", rows)
	d.add(`<g fill="%s">%s</g>`, d.linear(grad{0, 0, 0.3, 1, []stop{{0, light}, {1, shade}}}), rows)
	d.add(`<g clip-path="url(#cell)">`)
	paint(d)
	d.add(`</g>`)
	return d.String()
}

// specks scatters n soft ovals of colour c over the cell, seed choosing
// where, size their radius.
func specks(d *doc, seed rng, n int, c string, size float64) {
	d.add(`<g transform="%s" fill="%s">`, onCell, c)
	for range n {
		u, v := 1+seed.next()*14, 1+seed.next()*14
		d.add(`<ellipse cx="%.2f" cy="%.2f" rx="%.2f" ry="%.2f" opacity="0.55"/>`, u, v, size*(0.6+seed.next()*0.8), size*(0.5+seed.next()*0.6))
	}
	d.add(`</g>`)
}
