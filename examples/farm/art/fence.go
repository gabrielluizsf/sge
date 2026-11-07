package art

import "fmt"

// fenceTile is a run of fence along a cell's row, seen from the corner:
// two rails and a post in the middle of the cell.
func fenceTile() string { return fenceArms(west, east) }

// fenceDownTile is the same fence turned a quarter, along a cell's column.
func fenceDownTile() string { return fenceArms(north, south) }

// fenceCornerTile is a fence turning on its cell: its post in the middle and
// its rails running from it to the edges a and b.
func fenceCornerTile(a, b pt) string { return fenceArms(a, b) }

// The middles of a cell's edges on its ground.
var (
	north = pt{Size / 2, 0}
	south = pt{Size / 2, Size}
	west  = pt{0, Size / 2}
	east  = pt{Size, Size / 2}
)

// fenceArms is a piece of fence with rails from its post, in the middle of
// the cell's ground, to each of the points a and b on it.
func fenceArms(a, b pt) string {
	d := newDoc(IsoW, IsoH)
	for _, end := range []pt{a, b} {
		for _, h := range []float64{3.5, 8} {
			ax, ay := iso(Size/2, Size/2, h)
			bx, by := iso(end.x, end.y, h)
			d.add(`<path d="M%g %g L%g %g" stroke="%s" stroke-width="2.2" stroke-linecap="round"/>`, ax, ay, bx, by, line)
			d.add(`<path d="M%g %g L%g %g" stroke="%s" stroke-width="1.2" stroke-linecap="round"/>`, ax, ay, bx, by, woodBase)
			d.add(`<path d="M%g %g L%g %g" stroke="%s" stroke-width="0.4" stroke-linecap="round"/>`, ax, ay-0.35, bx, by-0.35, woodLight)
		}
	}
	x, y := iso(Size/2, Size/2, 0)
	post := fmt.Sprintf(`M%g %g l0 -11 l1.4 -0.9 l1.4 0.9 l0 11 l-1.4 0.8 Z`, x-1.4, y)
	d.add(`<path d="%s" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`, post, woodBase, line)
	d.add(`<path d="M%g %g l0 -11 l1.4 0.9 l0 10.9 Z" fill="%s"/>`, x, y+0.8, woodShade)
	d.add(`<path d="M%g %g l1.4 -0.9 l1.4 0.9 l-1.4 0.8 Z" fill="%s" stroke="%s" stroke-width="0.4"/>`, x-1.4, y-11, woodLight, line)
	return d.String()
}
