package art

import "fmt"

// fenceTile is a run of fence along a cell's row, seen from the corner:
// two rails from side to side, so a row of them joins up, and a post in
// the middle, wood lit from the top left.
func fenceTile() string {
	d := newDoc(IsoW, IsoH)
	for _, h := range []float64{3.5, 8} {
		ax, ay := iso(0, Size/2, h)
		bx, by := iso(Size, Size/2, h)
		d.add(`<path d="M%g %g L%g %g" stroke="%s" stroke-width="2.2" stroke-linecap="round"/>`, ax, ay, bx, by, line)
		d.add(`<path d="M%g %g L%g %g" stroke="%s" stroke-width="1.2" stroke-linecap="round"/>`, ax, ay, bx, by, woodBase)
		d.add(`<path d="M%g %g L%g %g" stroke="%s" stroke-width="0.4" stroke-linecap="round"/>`, ax, ay-0.35, bx, by-0.35, woodLight)
	}
	x, y := iso(Size/2, Size/2, 0)
	post := fmt.Sprintf(`M%g %g l0 -11 l1.4 -0.9 l1.4 0.9 l0 11 l-1.4 0.8 Z`, x-1.4, y)
	d.add(`<path d="%s" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`, post, woodBase, line)
	d.add(`<path d="M%g %g l0 -11 l1.4 0.9 l0 10.9 Z" fill="%s"/>`, x, y+0.8, woodShade)
	d.add(`<path d="M%g %g l1.4 -0.9 l1.4 0.9 l-1.4 0.8 Z" fill="%s" stroke="%s" stroke-width="0.4"/>`, x-1.4, y-11, woodLight, line)
	return d.String()
}
