package art

// The ground's tiles: grass in tufts that lean as the wind goes through,
// flowers among it, the sandy path and the soil of the plots, dry and wet.

// tufts are where the grass's blades grow on the cell, u and v.
var tufts = [][2]float64{{3, 2.5}, {9, 2}, {13.5, 4.5}, {6, 6.5}, {11, 8.5}, {2.5, 9.5}, {7.5, 11}, {13, 12.5}, {4, 14}, {10, 14.5}, {14.5, 9}, {1.5, 5.5}}

func grassTile(lean float64) string {
	return groundTile(grassLight, grassShade, func(d *doc) {
		specks(d, 7, 10, grassShade, 2.2)
		specks(d, 11, 8, grassLight, 1.6)
		blades(d, lean)
	})
}

// blades draws a tuft of three upright blades at each of tufts, their
// tips leant by lean.
func blades(d *doc, lean float64) {
	for i, t := range tufts {
		x, y := iso(t[0], t[1], 0)
		c := []string{grassDeep, grassShade, grassLight}[i%3]
		for _, b := range [][2]float64{{-0.8, 2.2}, {0, 3}, {0.8, 2}} {
			d.add(`<path d="M%.2f %.2f q%.2f %.2f %.2f %.2f" stroke="%s" stroke-width="0.55" stroke-linecap="round" fill="none"/>`,
				x+b[0], y, lean*0.6, -b[1]*0.6, b[0]*0.4+lean, -b[1], c)
		}
	}
}

func flowersTile(lean float64) string {
	return groundTile(grassLight, grassShade, func(d *doc) {
		specks(d, 7, 10, grassShade, 2.2)
		blades(d, lean)
		for i, f := range [][2]float64{{4, 4}, {11, 9}, {6, 12.5}, {12.5, 3}} {
			x, y := iso(f[0], f[1], 1.2)
			x += lean * 0.5
			petal := []string{pink, white, yellow, pink}[i]
			d.add(`<path d="M%.2f %.2f l0 1.4" stroke="%s" stroke-width="0.4"/>`, x, y, grassDeep)
			for _, p := range [][2]float64{{0, -0.8}, {0.8, -0.2}, {0.5, 0.7}, {-0.5, 0.7}, {-0.8, -0.2}} {
				d.add(`<circle cx="%.2f" cy="%.2f" r="0.62" fill="%s"/>`, x+p[0], y+p[1], petal)
			}
			d.add(`<circle cx="%.2f" cy="%.2f" r="0.45" fill="%s"/>`, x, y, orange)
		}
	})
}

func pathTile() string {
	return groundTile(sandLight, sandShade, func(d *doc) {
		specks(d, 3, 9, sandShade, 1.8)
		specks(d, 5, 6, sandLight, 1.4)
		d.add(`<g transform="%s">`, onCell)
		for _, p := range [][3]float64{{4, 5, 0.9}, {11, 3, 0.6}, {8, 11, 1}, {13, 13, 0.7}, {3, 13, 0.6}} {
			d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s" stroke="%s" stroke-width="0.25"/>`, p[0], p[1], p[2], p[2]*0.8, sandLight, sandShade)
		}
		d.add(`</g>`)
	})
}

// soilTile is a plot, light and base and shade its colours: furrows
// across it, each a dark trough under a lit ridge; wet soil glistens.
func soilTile(light, base, shade string, wet bool) string {
	return groundTile(base, shade, func(d *doc) {
		d.add(`<g transform="%s" stroke-linecap="round">`, onCell)
		for v := 2.0; v < 16; v += 4 {
			d.add(`<path d="M0.5 %g L15.5 %g" stroke="%s" stroke-width="1.5"/>`, v+1, v+1, shade)
			d.add(`<path d="M0.5 %g L15.5 %g" stroke="%s" stroke-width="0.8"/>`, v, v, light)
		}
		d.add(`</g>`)
		if wet {
			for _, p := range [][2]float64{{4, 6}, {10, 4}, {12, 11}, {6, 12}} {
				x, y := iso(p[0], p[1], 0)
				d.add(`<ellipse cx="%g" cy="%g" rx="1.4" ry="0.45" fill="%s" opacity="0.5"/>`, x, y, waterLight)
			}
		}
	})
}
