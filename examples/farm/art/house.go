package art

// A house's picture: HouseW by HouseH, its two by two cells' far corner at
// houseX, houseY, so the middle of its cells is at HousePivot.
const (
	HouseW         = 4*Size + 8
	HouseH         = 5 * Size
	houseX, houseY = HouseW / 2, 3 * Size
	// Houses is how many kinds of houses there are.
	Houses = 4
)

// HousePivot is the point of a house's picture, as fractions of its width
// and height, that stands on the middle of its cells.
var HousePivot = [2]float64{houseX / float64(HouseW), (houseY + Size) / float64(HouseH)}

// houseStyle is a kind of house: its walls' and roof's light, base and
// shade, its trim, and whether it has a chimney and a flower box.
type houseStyle struct {
	walls, roof [3]string
	trim        string
	chimney     bool
	flowers     bool
	bricks      bool
}

var houseStyles = [Houses]houseStyle{
	{walls: [3]string{woodLight, woodBase, woodShade}, roof: [3]string{redShade, "#8e2f2a", "#6a2420"}, trim: barkShade, chimney: true},
	{walls: [3]string{brickLight, brickBase, brickShade}, roof: [3]string{slateLight, slateBase, slateShade}, trim: white, bricks: true},
	{walls: [3]string{skyLight, skyBase, skyShade}, roof: [3]string{tileLight, tileBase, tileShade}, trim: white, flowers: true},
	{walls: [3]string{plasterLight, plasterBase, plasterSh}, roof: [3]string{mossLight, mossBase, mossShade}, trim: woodShade, chimney: true, flowers: true},
}

// House is a house of the given kind on its two by two cells: walls under
// a gabled roof, windows, and the door on the side facing the viewer when
// front is set, on the far side otherwise, unseen.
func House(kind int, front bool) string {
	s := houseStyles[(kind%Houses+Houses)%Houses]
	d := newDoc(HouseW, HouseH)
	d.origin = pt{houseX, houseY}
	c := d.point(Size, Size, 0)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="30" ry="15" fill="%s" opacity="0.3"/>`, c.x, c.y, shadowInk)
	const lo, hi, wall, ridge = 2.0, 30.0, 18.0, 31.0
	d.box([4]float64{lo, lo, hi, hi}, [2]float64{0, wall}, [3]string{s.walls[0], s.walls[1], s.walls[2]})
	if s.bricks {
		for h := 3.0; h < wall; h += 3 {
			d.add(`<path d="M%s L%s L%s" stroke="%s" stroke-width="0.35" fill="none" opacity="0.7"/>`, d.at(lo, hi, h), d.at(hi, hi, h), d.at(hi, lo, h), s.walls[2])
		}
	} else {
		for h := 4.0; h < wall; h += 4 {
			d.add(`<path d="M%s L%s L%s" stroke="%s" stroke-width="0.35" fill="none" opacity="0.6"/>`, d.at(lo, hi, h), d.at(hi, hi, h), d.at(hi, lo, h), s.walls[2])
		}
	}
	// The gable at the near end, under the roof's ridge.
	d.add(`<path d="M%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`, d.at(hi, lo, wall), d.at(hi, hi, wall), d.at(hi, 16, ridge), s.walls[2], line)
	window(d, s, [3]float64{hi, 16, 21})
	// Windows on both faces, and the door, or a third window.
	window(d, s, [3]float64{hi, 9, 10})
	window(d, s, [3]float64{hi, 23, 10})
	window(d, s, [3]float64{22, hi, 10})
	if front {
		door(d)
	} else {
		window(d, s, [3]float64{9, hi, 10})
	}
	if s.chimney {
		d.box([4]float64{6, 6, 10, 10}, [2]float64{wall, ridge + 3}, [3]string{brickLight, brickBase, brickShade})
	}
	// The roof's slopes: the far one first, the near one over it.
	roof := d.linear(grad{0, 0, 0, 1, []stop{{0, s.roof[0]}, {1, s.roof[1]}}})
	d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, d.at(0, 16, ridge), d.at(32, 16, ridge), d.at(32, -1, wall-2), d.at(0, -1, wall-2), s.roof[2], line)
	d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, d.at(0, 16, ridge), d.at(32, 16, ridge), d.at(32, 33, wall-2), d.at(0, 33, wall-2), roof, line)
	for u := 4.0; u < 32; u += 4 {
		d.add(`<path d="M%s L%s" stroke="%s" stroke-width="0.4" opacity="0.8"/>`, d.at(u, 17, ridge-0.6), d.at(u, 32.4, wall-1.6), s.roof[2])
	}
	d.add(`<path d="M%s L%s" stroke="%s" stroke-width="0.8" stroke-linecap="round"/>`, d.at(0, 16, ridge), d.at(32, 16, ridge), s.roof[0])
	return d.String()
}

// window is a window on a face of the house, its middle at u, v, h: on
// the near end when u is the wall's, on the front otherwise.
func window(d *doc, s houseStyle, at [3]float64) {
	u, v, h := at[0], at[1], at[2]
	var pts [4]string
	if u >= 30 {
		pts = [4]string{d.at(u, v-2.6, h-2.6), d.at(u, v+2.6, h-2.6), d.at(u, v+2.6, h+2.6), d.at(u, v-2.6, h+2.6)}
	} else {
		pts = [4]string{d.at(u-2.6, v, h-2.6), d.at(u+2.6, v, h-2.6), d.at(u+2.6, v, h+2.6), d.at(u-2.6, v, h+2.6)}
	}
	d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.9" stroke-linejoin="round"/>`, pts[0], pts[1], pts[2], pts[3], d.down(glass, glassShade), s.trim)
	d.add(`<path d="M%s L%s L%s L%s Z" fill="none" stroke="%s" stroke-width="0.35"/>`, pts[0], pts[1], pts[2], pts[3], line)
	if s.flowers && h < 15 {
		for i := range 3 {
			c := []string{pink, red, yellow}[i]
			f := (float64(i) + 0.5) / 3
			var p pt
			if u >= 30 {
				p = d.point(u, v-2.6+5.2*f, h-2.6)
			} else {
				p = d.point(u-2.6+5.2*f, v, h-2.6)
			}
			d.add(`<circle cx="%.2f" cy="%.2f" r="0.9" fill="%s" stroke="%s" stroke-width="0.3"/>`, p.x, p.y-0.6, c, line)
		}
	}
}

// door is the front door, with a step and a knob.
func door(d *doc) {
	const v = 30
	d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, d.at(8, v, 0), d.at(13, v, 0), d.at(13, v, 11), d.at(8, v, 11), d.down(barkBase, barkShade), line)
	k := d.point(12, v, 5)
	d.add(`<circle cx="%.2f" cy="%.2f" r="0.5" fill="%s"/>`, k.x, k.y, brass)
	d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.4"/>`, d.at(7, v, 0), d.at(14, v, 0), d.at(14, v+2, 0), d.at(7, v+2, 0), metalShade, line)
}
