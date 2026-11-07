package bakes

import (
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// barrelNode names the root of a rider's barrel.
const barrelNode = "barrel"

// rideOnBarrel gives m, a person already animated, the barrel that stands for
// the horse's body while they ride, get on or get down, if they do: a further
// root of the model for each.
func (b *Baker) rideOnBarrel(m *gltf.Model, motions []human.Motion) error {
	if !hasMotion(motions, human.Ride) {
		return nil
	}
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	for _, mo := range human.Riding {
		if err := human.Lay(m, building.Barrel.Model(palette), human.BedPlace{At: [3]float64{0, human.BarrelHeight, 0}, Scale: 1, During: mo.Name}); err != nil {
			return err
		}
		m.Nodes[m.Roots[len(m.Roots)-1]].Name = barrelNode
	}
	return nil
}

// rider is the picture of a rider as shot: baked round the barrel, so the far
// leg is hidden where the horse's body would hide it, and the barrel then
// taken away, with no shadow, which the horse's own picture has. It is drawn
// over the horse, on the same spot.
func rider(m *gltf.Model, sh bake.Shot) bake.Sprite {
	o := Options()
	o.Shadow = 0
	with := bake.Render(m, sh, o)
	return withoutBarrel(with, bake.Render(only(m, true), sh, o))
}

// only is m with its barrels alone, or with all but them.
func only(m *gltf.Model, barrels bool) *gltf.Model {
	out := *m
	out.Roots = nil
	for _, r := range m.Roots {
		if (m.Nodes[r].Name == barrelNode) == barrels {
			out.Roots = append(out.Roots, r)
		}
	}
	return &out
}

// withoutBarrel is the rider alone: what the barrel alone does not show.
// Where the picture is the barrel's, or close to it on an edge between the
// two, nothing is left.
func withoutBarrel(sp, bare bake.Sprite) bake.Sprite {
	return withoutBed(clearNear(sp, bare), bare)
}

// clearNear is sp with the pixels of mostly the barrel's colour made the
// barrel's own: the edge between rider and barrel is a mix of the two, which
// is not to stay as a pink fringe.
func clearNear(sp, bare bake.Sprite) bake.Sprite {
	out := *sp.Image
	out.Pix = append([]uint8(nil), sp.Image.Pix...)
	b := out.Rect
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := out.NRGBAAt(x, y)
			if c.A > 0 && pinkish(c.R, c.G, c.B) {
				out.SetNRGBA(x, y, bareAt(bare, sp, x, y))
			}
		}
	}
	return bake.Sprite{Image: &out, Pivot: sp.Pivot}
}

// pinkish reports whether a colour has more of the barrel's magenta than any
// rider's clothes or skin: red and blue both well over green.
func pinkish(r, g, b uint8) bool {
	return math.Min(float64(r), float64(b)) > float64(g)*1.35+18
}

// bareAt is bare's pixel under x, y of sp, whose pivots are one spot.
func bareAt(bare, sp bake.Sprite, x, y int) colorNRGBA {
	dx, dy := x-sp.Pivot.X+bare.Pivot.X, y-sp.Pivot.Y+bare.Pivot.Y
	if r := bare.Image.Rect; dx < r.Min.X || dy < r.Min.Y || dx >= r.Max.X || dy >= r.Max.Y {
		return colorNRGBA{}
	}
	return bare.Image.NRGBAAt(dx, dy)
}

// colorNRGBA is short for the pixels this file works in.
type colorNRGBA = color.NRGBA
