package art

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/sge/terrain"
)

// The water's edge, in the cell's units: how far in from a side facing
// land the bank is, how round the bank turns where two such sides meet,
// and how far land reaches into the water at a corner of land between two
// sides of water.
const (
	bankIn    = 2.5
	bankRound = 4
	notch     = 2.6
	// past is how far beyond a side facing water the water is drawn,
	// clipped to the cell, so tiles of water meet edge to edge.
	past = 2
)

// shoreTile is water f steps into its ripple in a tile whose neighbours
// are water as terrain.Mask says: open water toward water, a bank of foam
// over shallows toward land, and nothing beyond it, so the grass under
// the water layer shows.
func shoreTile(mask, f int) string {
	d := newDoc(IsoW, IsoH)
	body, notches := waterShape(mask)
	fmt.Fprintf(&d.defs, `<clipPath id="cell">%s</clipPath>`+"\n", diamondRows())
	fmt.Fprintf(&d.defs, `<clipPath id="body"><path d="%s"/></clipPath>`+"\n", body)
	fmt.Fprintf(&d.defs, `<clipPath id="wet"><path d="%s M-9 -9 h34 v34 h-34 z" clip-rule="evenodd" transform="%s"/></clipPath>`+"\n", notches, onCell)
	d.add(`<g clip-path="url(#cell)"><g clip-path="url(#wet)">`)
	d.add(`<g transform="%s">`, onCell)
	d.add(`<path d="%s" fill="%s"/>`, body, waterBase)
	// Within the water: its glints, and along the bank the shallows and a
	// line of foam, their inner halves.
	d.add(`<g clip-path="url(#body)">`)
	ripples(d, f)
	d.add(`<path d="%s" fill="none" stroke="%s" stroke-width="4" opacity="0.6"/>`, body, waterLight)
	d.add(`<path d="%s" fill="none" stroke="%s" stroke-width="1.4"/>`, body, foam)
	d.add(`<path d="%s" fill="none" stroke="%s" stroke-width="1.4"/>`, notchRings(mask), foam)
	d.add(`</g>`)
	d.add(`</g></g></g>`)
	return d.String()
}

// ripples are light glints on the water, drifting with frame f; they
// repeat every cell, so tiles of water ripple as one sheet.
func ripples(d *doc, f int) {
	shift := float64(f) * 4
	for _, r := range [][2]float64{{2, 3}, {9, 6}, {5, 11}, {12, 13}, {14, 1}} {
		u := r[0] + shift
		for u >= 16 {
			u -= 16
		}
		d.add(`<path d="M%.2f %.2f q1.2 -0.9 2.4 0" stroke="%s" stroke-width="0.55" fill="none" stroke-linecap="round" opacity="0.85"/>`, u-1.2, r[1], waterLight)
	}
}

// waterShape is the water's outline on the cell, and the circles of land
// at its notched corners, as SVG paths in the cell's units.
func waterShape(mask int) (body, notches string) {
	at := func(bit int, near, far float64) float64 {
		if mask&bit != 0 {
			return far
		}
		return near
	}
	x0, y0 := at(terrain.West, bankIn, -past), at(terrain.North, bankIn, -past)
	x1, y1 := at(terrain.East, Size-bankIn, Size+past), at(terrain.South, Size-bankIn, Size+past)
	land := func(bit int) bool { return mask&bit == 0 }
	var b strings.Builder
	// Corners clockwise from the north-west, each round where both its
	// sides face land.
	type corner struct {
		x, y   float64
		a, b   int // its sides
		dx, dy float64
	}
	cs := []corner{{x0, y0, terrain.North, terrain.West, 1, 1}, {x1, y0, terrain.North, terrain.East, -1, 1}, {x1, y1, terrain.South, terrain.East, -1, -1}, {x0, y1, terrain.South, terrain.West, 1, -1}}
	for i, c := range cs {
		round := land(c.a) && land(c.b)
		cmd := "L"
		if i == 0 {
			cmd = "M"
		}
		if !round {
			fmt.Fprintf(&b, "%s%.2f %.2f ", cmd, c.x, c.y)
			continue
		}
		// In along the side before the corner, round it, out along the next.
		before := [4][2]float64{{c.x, c.y + bankRound}, {c.x - bankRound, c.y}, {c.x, c.y - bankRound}, {c.x + bankRound, c.y}}[i]
		after := [4][2]float64{{c.x + bankRound, c.y}, {c.x, c.y + bankRound}, {c.x - bankRound, c.y}, {c.x, c.y - bankRound}}[i]
		fmt.Fprintf(&b, "%s%.2f %.2f A%d %d 0 0 1 %.2f %.2f ", cmd, before[0], before[1], bankRound, bankRound, after[0], after[1])
	}
	b.WriteString("Z")
	return b.String(), notchCircles(mask)
}

// notchCorners are the cell's corners where land reaches in between two
// sides of water: the corner bit unset beside its two sides set.
func notchCorners(mask int) [][2]float64 {
	var out [][2]float64
	for _, c := range [4]struct {
		corner, a, b int
		u, v         float64
	}{{terrain.NorthEast, terrain.North, terrain.East, Size, 0}, {terrain.SouthEast, terrain.South, terrain.East, Size, Size},
		{terrain.SouthWest, terrain.South, terrain.West, 0, Size}, {terrain.NorthWest, terrain.North, terrain.West, 0, 0}} {
		if mask&c.a != 0 && mask&c.b != 0 && mask&c.corner == 0 {
			out = append(out, [2]float64{c.u, c.v})
		}
	}
	return out
}

// notchCircles are the circles of land at the notched corners, as a path.
func notchCircles(mask int) string {
	var b strings.Builder
	for _, c := range notchCorners(mask) {
		fmt.Fprintf(&b, "M%.2f %.2f a%g %g 0 1 0 %g 0 a%g %g 0 1 0 %g 0 ", c[0]-notch, c[1], notch, notch, 2*notch, notch, notch, -2*notch)
	}
	return b.String()
}

// notchRings are the foam round the notched corners, as a path.
func notchRings(mask int) string {
	if s := notchCircles(mask); s != "" {
		return s
	}
	return "M0 0"
}
