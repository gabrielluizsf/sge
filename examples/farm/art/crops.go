package art

// The crops, each in its stages from seed to ripe, standing on their
// cell: stage 0 the seed under a mound, 3 ripe, and 4 ripe leant in the
// wind, asking to be picked.

// leaf is a pointed leaf from x, y, its tip dx, dy away.
func leaf(d *doc, x, y float64, tip [2]float64) {
	mx, my := x+tip[0]/2, y+tip[1]/2
	nx, ny := -tip[1]*0.28, tip[0]*0.28
	d.add(`<path d="M%.2f %.2f Q%.2f %.2f %.2f %.2f Q%.2f %.2f %.2f %.2f Z" fill="%s" stroke="%s" stroke-width="0.4" stroke-linejoin="round"/>`,
		x, y, mx+nx, my+ny, x+tip[0], y+tip[1], mx-nx, my-ny, x, y, leafBase, line)
	d.add(`<path d="M%.2f %.2f L%.2f %.2f" stroke="%s" stroke-width="0.3" opacity="0.8"/>`, x, y, x+tip[0]*0.8, y+tip[1]*0.8, leafGlow)
}

// mound is the earth heaped over a seed at the cell's middle.
func mound(d *doc) {
	x, y := iso(Size/2, Size/2, 0)
	d.add(`<ellipse cx="%g" cy="%g" rx="3.4" ry="1.5" fill="%s" stroke="%s" stroke-width="0.4"/>`, x, y, soilLight, line)
	d.add(`<ellipse cx="%g" cy="%g" rx="1.6" ry="0.6" fill="%s"/>`, x-0.6, y-0.4, sandLight)
}

// sprout is a crop's tuft of n leaves of length size, leant by lean.
func sprout(d *doc, n int, size, lean float64) {
	x, y := iso(Size/2, Size/2, 0)
	for i := range n {
		a := (float64(i)/float64(max(n-1, 1)) - 0.5) * 1.6
		leaf(d, x, y, [2]float64{a*size*0.6 + lean, -size * (1 - a*a*0.25)})
	}
}

func carrotTile(stage int, lean float64) string {
	d := newDoc(IsoW, IsoH)
	shadow(d, 4, 1.6)
	mound(d)
	switch stage {
	case 1:
		sprout(d, 2, 3, lean)
	case 2:
		sprout(d, 4, 5, lean)
	case 3:
		x, y := iso(Size/2, Size/2, 0)
		d.add(`<path d="M%g %g q2.2 -0.4 2.2 0.6 q-0.8 3 -2.2 4.6 q-1.4 -1.6 -2.2 -4.6 q0 -1 2.2 -0.6 Z" fill="%s" stroke="%s" stroke-width="0.4" transform="translate(0 -1.6)"/>`, x, y, orange, line)
		d.add(`<path d="M%g %g l2 0" stroke="%s" stroke-width="0.35" transform="translate(-1 0.2)"/>`, x, y, orangeShade)
		sprout(d, 5, 7, lean)
	}
	return d.String()
}

func tomatoTile(stage int, lean float64) string {
	d := newDoc(IsoW, IsoH)
	shadow(d, 4, 1.6)
	mound(d)
	x, y := iso(Size/2, Size/2, 0)
	switch stage {
	case 1:
		sprout(d, 2, 3, lean)
	case 2, 3:
		tall := map[int]float64{2: 9, 3: 13}[stage]
		d.add(`<path d="M%g %g l0 %g" stroke="%s" stroke-width="0.7" stroke-linecap="round"/>`, x+2.2, y+0.6, -tall-2, woodBase)
		d.add(`<path d="M%g %g q%g %g %g %g" stroke="%s" stroke-width="0.8" fill="none" stroke-linecap="round"/>`, x, y, lean, -tall/2, lean*0.6, -tall, leafShade)
		for i := 1.0; i < tall; i += 3 {
			side := map[bool]float64{true: 1, false: -1}[int(i)%2 == 0]
			leaf(d, x+lean*i/tall, y-i, [2]float64{3 * side, -1.4})
		}
		if stage == 3 {
			for _, t := range [][2]float64{{-2.2, -5}, {1.8, -7.5}, {-1, -10}, {2.4, -3.2}} {
				d.add(`<circle cx="%g" cy="%g" r="1.6" fill="%s" stroke="%s" stroke-width="0.4"/>`, x+t[0]+lean*0.6, y+t[1], red, line)
				d.add(`<circle cx="%g" cy="%g" r="0.5" fill="%s" opacity="0.8"/>`, x+t[0]+lean*0.6-0.5, y+t[1]-0.5, white)
			}
		}
	}
	return d.String()
}
