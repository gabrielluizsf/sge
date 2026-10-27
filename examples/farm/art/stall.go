package art

import "fmt"

// The stall's picture: StallW by StallH, the far corner of its first
// cell at stallX, stallY, so the point it stands on, the middle of its two
// cells along a row, is StallPivot.
const (
	StallW         = 3 * Size
	StallH         = 3*Size + Size/2
	stallX, stallY = Size, 2 * Size
)

// StallPivot is the point of the stall's picture, as fractions of its
// width and height, that stands on the middle of its two cells.
var StallPivot = [2]float64{(stallX + Size - Size/2) / float64(StallW), (stallY + (Size+Size/2)/2) / float64(StallH)}

// stallAt is where the point u, v of the stall's two cells, h up, is
// drawn; sp is that point written for a path.
func stallAt(u, v, h float64) (x, y float64) { return stallX + u - v, stallY + (u+v)/2 - h }

func sp(u, v, h float64) string {
	x, y := stallAt(u, v, h)
	return fmt.Sprintf("%.2f %.2f", x, y)
}

// box is a block from u0, v0 to u1, v1, h0 up to h1: its two faces to the
// screen in left and right, its top in top, outlined.
func (d *doc) box(c [4]float64, h [2]float64, paint [3]string) {
	u0, v0, u1, v1 := c[0], c[1], c[2], c[3]
	d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`, sp(u0, v1, h[0]), sp(u1, v1, h[0]), sp(u1, v1, h[1]), sp(u0, v1, h[1]), paint[1], line)
	d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`, sp(u1, v0, h[0]), sp(u1, v1, h[0]), sp(u1, v1, h[1]), sp(u1, v0, h[1]), paint[2], line)
	d.add(`<path d="M%s L%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`, sp(u0, v0, h[1]), sp(u1, v0, h[1]), sp(u1, v1, h[1]), sp(u0, v1, h[1]), paint[0], line)
}

// Stall is the shop's stall: a wooden counter with baskets of the farm's
// produce, under a striped awning with a scalloped edge on four poles.
func Stall() string {
	d := newDoc(StallW, StallH)
	d.add(`<ellipse cx="%g" cy="%g" rx="20" ry="7" fill="%s" opacity="0.3"/>`, StallPivot[0]*StallW, StallPivot[1]*StallH, shadowInk)
	pole := func(u, v float64) {
		d.box([4]float64{u, v, u + 1.4, v + 1.4}, [2]float64{0, 26}, [3]string{woodLight, woodBase, woodShade})
	}
	pole(2, 3)
	pole(28.6, 3)
	d.box([4]float64{2, 5, 30, 14}, [2]float64{0, 10}, [3]string{woodLight, woodBase, woodShade})
	for i := 1.0; i < 4; i++ { // planks
		d.add(`<path d="M%s L%s" stroke="%s" stroke-width="0.4"/>`, sp(2+i*7, 14, 0.4), sp(2+i*7, 14, 9.6), woodShade)
	}
	for i, c := range [][2]string{{orange, orangeShade}, {red, redShade}, {yellow, sandShade}} {
		u := 6 + float64(i)*8
		d.box([4]float64{u, 7, u + 5, 11}, [2]float64{10, 12}, [3]string{woodShade, woodBase, woodShade})
		for _, o := range [][2]float64{{1.3, 1}, {3.4, 1.3}, {2.3, 2.8}} {
			x, y := stallAt(u+o[0], 7+o[1], 13)
			d.add(`<circle cx="%.2f" cy="%.2f" r="1.15" fill="%s" stroke="%s" stroke-width="0.3"/>`, x, y, c[0], line)
			d.add(`<circle cx="%.2f" cy="%.2f" r="0.4" fill="%s" opacity="0.7"/>`, x-0.4, y-0.4, white)
		}
	}
	pole(2, 13)
	pole(28.6, 13)
	// The awning: its top in stripes, its front hem scalloped.
	d.box([4]float64{0, 1, 32, 16}, [2]float64{26, 27.5}, [3]string{white, redShade, redShade})
	for u := 0.0; u < 32; u += 4 {
		d.add(`<path d="M%s L%s L%s L%s Z" fill="%s"/>`, sp(u, 1, 27.5), sp(u+2, 1, 27.5), sp(u+2, 16, 27.5), sp(u, 16, 27.5), red)
	}
	for u := 0.0; u < 32; u += 2 {
		c := map[bool]string{true: red, false: white}[int(u/2)%2 == 0]
		d.add(`<path d="M%s Q%s %s Z" fill="%s" stroke="%s" stroke-width="0.4"/>`, sp(u, 16, 26), sp(u+1, 16, 24), sp(u+2, 16, 26), c, line)
	}
	return d.String()
}
