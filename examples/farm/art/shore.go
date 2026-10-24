package art

import (
	"math"

	"github.com/gabrielluizsf/sge/terrain"
)

// corners are the diagonal bits of a mask, and the two sides beside each.
var corners = [4][3]int{
	{terrain.NorthEast, terrain.North, terrain.East},
	{terrain.SouthEast, terrain.South, terrain.East},
	{terrain.SouthWest, terrain.South, terrain.West},
	{terrain.NorthWest, terrain.North, terrain.West},
}

// shore is water f steps into its ripple in a tile whose neighbours are
// water as terrain.Mask says. A side facing land ends two pixels in,
// with an outline and a line of foam, and transparent beyond, so the
// grass under the water shows; where two such sides meet the bank turns
// round. A corner of land between two sides of water, a diagonal bit
// unset beside two set, bends the bank of the tiles beside it round the
// corner.
func shore(mask, f int) Grid {
	g := water(f)
	const inset, r = 2.0, 4.0
	// The water's edge on each side: at the tile's edge where water goes
	// on, inset where land begins.
	lo := func(bit int) float64 {
		if mask&bit != 0 {
			return -100
		}
		return inset
	}
	left, top := lo(terrain.West), lo(terrain.North)
	right, bottom := float64(Size)-lo(terrain.East), float64(Size)-lo(terrain.South)
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			px, py := float64(x)+0.5, float64(y)+0.5
			// How far in from the bank the pixel's centre is.
			d := min(px-left, right-px, py-top, bottom-py)
			for _, c := range [][4]float64{{left, top, 1, 1}, {right, top, -1, 1}, {left, bottom, 1, -1}, {right, bottom, -1, -1}} {
				if c[0] < 0 || c[0] > Size || c[1] < 0 || c[1] > Size {
					continue // a side of this corner goes on as water
				}
				cx, cy := c[0]+c[2]*r, c[1]+c[3]*r
				if (px-cx)*c[2] < 0 && (py-cy)*c[3] < 0 {
					d = min(d, r-math.Hypot(px-cx, py-cy))
				}
			}
			for _, c := range corners {
				if mask&c[1] == 0 || mask&c[2] == 0 || mask&c[0] != 0 {
					continue
				}
				// The land is the corner square beyond the banks of the
				// two tiles beside: d is how far the pixel is from it.
				kx, ky := float64(Size)-inset, inset // north-east
				if c[2] == terrain.West {
					kx = inset
				}
				if c[1] == terrain.South {
					ky = float64(Size) - inset
				}
				// ox and oy are how far the pixel is past the corner, towards
				// the land on each axis: it is land where both are.
				sx, sy := 1.0, -1.0
				if c[2] == terrain.West {
					sx = -1
				}
				if c[1] == terrain.South {
					sy = 1
				}
				ox, oy := (px-kx)*sx, (py-ky)*sy
				switch {
				case ox > 0 && oy > 0:
					d = min(d, -min(ox, oy))
				case ox > 0:
					d = min(d, -oy)
				case oy > 0:
					d = min(d, -ox)
				default:
					d = min(d, math.Hypot(ox, oy))
				}
			}
			switch {
			case d < 0:
				g = g.set(x, y, '.')
			case d < 1:
				g = g.set(x, y, 'k')
			case d < 2:
				g = g.set(x, y, 'x')
			}
		}
	}
	return g
}
