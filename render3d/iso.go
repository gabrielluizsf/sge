package render3d

import (
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/pose"
)

// IsoDesc is the isometric game's camera: its Projection, the world
// point at the middle of the screen, Center, in the game's world units,
// a screen Width by Height pixels, and Depth, how far before and behind
// Center, along the line of sight, it sees.
type IsoDesc struct {
	Projection    iso.Projection
	Center        [3]float64
	Width, Height float64
	Depth         float64
}

// Iso is the orthographic camera that puts every world point exactly
// where d's Projection puts it, so that 3D models stand on the game's
// diamonds: the line of sight is the world's (1, 1, 1), along which the
// projection sees points one over another, and nearer is up it.
func Iso(d IsoDesc) Camera {
	p := d.Projection
	hx, hy := p.Width/2/p.Cell, p.Height/2/p.Cell
	sx, sy := 2/d.Width, 2/d.Height
	c := d.Center
	dz := 1 / (2 * d.Depth)
	vp := pose.Mat{
		hx * sx, -hy * sy, -dz, 0,
		-hx * sx, -hy * sy, -dz, 0,
		0, 2 * hy * sy, -dz, 0,
		-(c[0] - c[1]) * hx * sx, ((c[0]+c[1])*hy - 2*hy*c[2]) * sy, 0.5 + (c[0]+c[1]+c[2])*dz, 1,
	}
	return Camera{ViewProj: vp, Eye: [3]float64{c[0] + d.Depth, c[1] + d.Depth, c[2] + d.Depth}}
}
