package art

import "fmt"

// The inside of the houses: a plank floor, a rug, the doormat, the back
// walls in each kind of house's colours, and the furniture, each a tile
// standing on its cell as the farm's things do.

// floorTile is a cell of plank floor, the planks running across.
func floorTile() string {
	return groundTile(woodLight, woodShade, func(d *doc) {
		d.add(`<g transform="%s" stroke="%s" stroke-width="0.35">`, onCell, woodShade)
		for v := 4.0; v < Size; v += 4 {
			d.add(`<path d="M0 %g L16 %g"/>`, v, v)
		}
		for i, u := range []float64{5, 11, 8, 3} {
			d.add(`<path d="M%g %g L%g %g"/>`, u, float64(i)*4, u, float64(i)*4+4)
		}
		d.add(`</g>`)
	})
}

// rugTile is the floor under a woven rug of colour c.
func rugTile(c [3]string) string {
	return groundTile(woodLight, woodShade, func(d *doc) {
		d.add(`<g transform="%s">`, onCell)
		d.add(`<rect x="1.5" y="1.5" width="13" height="13" rx="1" fill="%s" stroke="%s" stroke-width="0.4"/>`, c[1], line)
		d.add(`<rect x="3.5" y="3.5" width="9" height="9" fill="none" stroke="%s" stroke-width="0.8"/>`, c[0])
		d.add(`<rect x="6" y="6" width="4" height="4" fill="%s"/>`, c[2])
		d.add(`</g>`)
	})
}

// matTile is the floor at the door, a doormat on it, the way out.
func matTile() string {
	return groundTile(woodLight, woodShade, func(d *doc) {
		d.add(`<g transform="%s">`, onCell)
		d.add(`<rect x="3" y="2" width="10" height="8" rx="1" fill="%s" stroke="%s" stroke-width="0.4"/>`, strawBase, line)
		for u := 4.0; u < 13; u += 1.5 {
			d.add(`<path d="M%g 2.6 L%g 9.4" stroke="%s" stroke-width="0.35"/>`, u, u, strawShade)
		}
		d.add(`</g>`)
	})
}

// wallH is how tall the back walls stand.
const wallH = 26

// wallTile is the back wall along a cell's far edges, north, west or
// both, in walls' colours, with a window when window is set.
func wallTile(walls [3]string, sides string, window bool) string {
	d := newDoc(IsoW, IsoH)
	d.origin = pt{16, 32}
	face := func(a, b [2]float64, c string) {
		d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`,
			d.at(a[0], a[1], 0), d.at(b[0], b[1], 0), d.at(b[0], b[1], wallH), d.at(a[0], a[1], wallH), c, line)
		d.add(`<path d="M%s L%s" stroke="%s" stroke-width="1.6"/>`, d.at(a[0], a[1], 1), d.at(b[0], b[1], 1), woodShade)
	}
	if sides == "west" || sides == "both" {
		face([2]float64{0, 0}, [2]float64{0, Size}, walls[2])
	}
	if sides == "north" || sides == "both" {
		face([2]float64{0, 0}, [2]float64{Size, 0}, walls[1])
		if window {
			d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="1" stroke-linejoin="round"/>`,
				d.at(4, 0, 10), d.at(12, 0, 10), d.at(12, 0, 19), d.at(4, 0, 19), d.down(skyLight, glassShade), white)
			d.add(`<path d="M%s L%s M%s L%s" stroke="%s" stroke-width="0.6"/>`, d.at(8, 0, 10), d.at(8, 0, 19), d.at(4, 0, 14.5), d.at(12, 0, 14.5), white)
		}
	}
	return d.String()
}

// furniture is a piece of furniture standing on its cell, drawn by paint
// with d's origin at the cell's far corner.
func furniture(paint func(d *doc)) string {
	d := newDoc(IsoW, IsoH)
	d.origin = pt{16, 32}
	shadow(d, 7, 3.2)
	paint(d)
	return d.String()
}

var wood3 = [3]string{woodLight, woodBase, woodShade}

func bedTile(blanket [3]string) string {
	return furniture(func(d *doc) {
		d.box([4]float64{2, 1, 14, 15}, [2]float64{0, 5}, wood3)
		d.box([4]float64{2.5, 1.5, 13.5, 14.5}, [2]float64{5, 7}, [3]string{white, shirtBase, shirtShade})
		d.box([4]float64{2.5, 6, 13.5, 14.5}, [2]float64{7, 7.6}, blanket)
		d.box([4]float64{4, 2, 12, 5}, [2]float64{7, 9}, [3]string{white, white, shirtShade})
		d.box([4]float64{1, 0, 15, 1.5}, [2]float64{0, 12}, wood3)
	})
}

func tableTile() string {
	return furniture(func(d *doc) {
		for _, c := range [][2]float64{{3, 3}, {12, 3}, {3, 12}, {12, 12}} {
			d.box([4]float64{c[0], c[1], c[0] + 1.2, c[1] + 1.2}, [2]float64{0, 8}, wood3)
		}
		d.box([4]float64{2, 2, 14, 14}, [2]float64{8, 9.5}, wood3)
		p := d.point(8, 8, 9.5)
		d.add(`<ellipse cx="%.2f" cy="%.2f" rx="2.4" ry="1.2" fill="%s" stroke="%s" stroke-width="0.4"/>`, p.x, p.y-0.6, white, line)
		d.add(`<circle cx="%.2f" cy="%.2f" r="0.9" fill="%s"/>`, p.x, p.y-1.2, red)
	})
}

func chairTile() string {
	return furniture(func(d *doc) {
		for _, c := range [][2]float64{{5, 5}, {10, 5}, {5, 10}, {10, 10}} {
			d.box([4]float64{c[0], c[1], c[0] + 1, c[1] + 1}, [2]float64{0, 5}, wood3)
		}
		d.box([4]float64{4.5, 4.5, 11.5, 11.5}, [2]float64{5, 6}, wood3)
		d.box([4]float64{4.5, 4.5, 11.5, 5.6}, [2]float64{6, 13}, wood3)
	})
}

func shelfTile() string {
	return furniture(func(d *doc) {
		d.box([4]float64{1, 1, 15, 5}, [2]float64{0, 22}, wood3)
		for i, h := range []float64{6, 12, 18} {
			for j := range 5 {
				c := []string{red, denimBase, leafBase, yellow, pink}[(i+j)%5]
				u := 2.5 + float64(j)*2.4
				d.box([4]float64{u, 4.6, u + 1.6, 5.6}, [2]float64{h, h + 4 - float64((i+j)%2)}, [3]string{c, c, line})
			}
		}
	})
}

func stoveTile() string {
	return furniture(func(d *doc) {
		d.box([4]float64{2, 2, 14, 14}, [2]float64{0, 9}, [3]string{slateLight, slateBase, slateShade})
		p := d.point(8, 14, 4)
		d.add(`<ellipse cx="%.2f" cy="%.2f" rx="2.6" ry="2.2" fill="%s" stroke="%s" stroke-width="0.4"/>`, p.x, p.y, orange, line)
		d.add(`<ellipse cx="%.2f" cy="%.2f" rx="1.2" ry="1.2" fill="%s"/>`, p.x, p.y+0.3, yellow)
		d.box([4]float64{6, 6, 10, 10}, [2]float64{9, 13}, [3]string{metal, metalShade, slateShade})
	})
}

func plantTile() string {
	return furniture(func(d *doc) {
		d.box([4]float64{5, 5, 11, 11}, [2]float64{0, 5}, [3]string{tileLight, tileBase, tileShade})
		p := d.point(8, 8, 5)
		for _, l := range [][2]float64{{-4, -6}, {4, -6}, {0, -9}, {-6, -2}, {6, -2}} {
			leaf(d, p.x, p.y, l)
		}
	})
}

func chestTile() string {
	return furniture(func(d *doc) {
		d.box([4]float64{3, 4, 13, 12}, [2]float64{0, 6}, wood3)
		d.box([4]float64{3, 4, 13, 12}, [2]float64{6, 8}, [3]string{woodLight, woodShade, woodShade})
		p := d.point(8, 12, 5.4)
		d.add(`<rect x="%.2f" y="%.2f" width="1.6" height="2" fill="%s" stroke="%s" stroke-width="0.3"/>`, p.x-0.8, p.y-1, brass, line)
	})
}

// forgeTile is the smith's furnace: a brick oven with a fire in its
// mouth, toward the viewer, and a chimney.
func forgeTile() string {
	return furniture(func(d *doc) {
		d.box([4]float64{1, 1, 15, 15}, [2]float64{0, 12}, [3]string{brickLight, brickBase, brickShade})
		p := d.point(8, 15, 5)
		d.add(`<path d="M%.2f %.2f a3.4 3.6 0 0 1 6.8 0 l0 2.6 l-6.8 0 Z" fill="%s" stroke="%s" stroke-width="0.5"/>`, p.x-3.4, p.y, line, line)
		d.add(`<path d="M%.2f %.2f q1.4 -4 2.6 -1 q0.8 -3.2 2 0 q0.8 2 -0.6 2.6 l-3.4 0 q-1.2 -0.6 -0.6 -1.6 Z" fill="%s"/>`, p.x-3, p.y+2.4, orange)
		d.add(`<path d="M%.2f %.2f q0.8 -2 1.6 0 Z" fill="%s"/>`, p.x-0.8, p.y+2.4, yellow)
		d.box([4]float64{9, 4, 13, 8}, [2]float64{12, 26}, [3]string{brickLight, brickBase, brickShade})
		for i := range 3 {
			s := d.point(11, 6, 28+float64(i)*3)
			d.add(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s" opacity="0.5"/>`, s.x+float64(i)*0.8, s.y, 1.2+float64(i)*0.5, metal)
		}
	})
}

// Interior is every tile of the houses' insides, by name: the floors,
// each kind of house's walls, and the furniture.
func Interior() map[string]string {
	s := map[string]string{
		"floor": floorTile(), "mat": matTile(),
		"bed": bedTile([3]string{red, redShade, redShade}), "bed_blue": bedTile([3]string{denimLight, denimBase, denimShade}),
		"table": tableTile(), "chair": chairTile(), "shelf": shelfTile(), "stove": stoveTile(),
		"plant": plantTile(), "chest": chestTile(), "forge": forgeTile(),
	}
	for i, c := range Shirts[:4] {
		s[fmt.Sprint("rug", i)] = rugTile(c)
	}
	for k, st := range houseStyles {
		walls := st.walls
		s[fmt.Sprint("wall_n", k)] = wallTile(walls, "north", false)
		s[fmt.Sprint("wall_nwin", k)] = wallTile(walls, "north", true)
		s[fmt.Sprint("wall_w", k)] = wallTile(walls, "west", false)
		s[fmt.Sprint("wall_nw", k)] = wallTile(walls, "both", false)
	}
	return s
}
