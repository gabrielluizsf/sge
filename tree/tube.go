package tree

import "math"

// limb is a trunk, branch or arm: a tube along path, r0 in radius at its
// start and r1 at its end, flared by flare at its foot, sides round; of
// material mat, its texture repeating every span metres, a whole number
// of times round it (tiles, or as its girth says); ribbed, if ribs is
// over 0, rib deep, as a cactus is; capped, with a rounded end, if capped.
type limb struct {
	path   []vec
	r0, r1 float64
	flare  float64
	sides  int
	mat    int
	span   float64
	tiles  int
	ribs   int
	rib    float64
	capped bool
}

// ring is a circle of a tube: its middle and radius.
type ring struct {
	at vec
	r  float64
}

// rings are the circles of limb l, and its rounded end if it has one.
func (l limb) rings() []ring {
	n := len(l.path)
	out := make([]ring, 0, n+3)
	for i, p := range l.path {
		t := float64(i) / float64(n-1)
		out = append(out, ring{p, (l.r0 + (l.r1-l.r0)*t) * (1 + l.flare*math.Pow(1-t, 8))})
	}
	if l.capped {
		_, dir := along(l.path, 1)
		for k := 1; k <= 3; k++ {
			th := float64(k) / 3 * math.Pi / 2
			out = append(out, ring{add(l.path[n-1], scale(dir, l.r1*math.Sin(th))), l.r1 * math.Cos(th)})
		}
	}
	return out
}

// tube adds limb l, its texture laid by its span.
func (g *grower) tube(l limb) {
	rings := l.rings()
	sides := max(l.sides, 5)
	if l.ribs > 0 {
		sides = max(sides, l.ribs*3)
	}
	span := l.span
	if span <= 0 {
		span = 1
	}
	tiles := float64(l.tiles)
	if l.tiles == 0 {
		tiles = math.Max(1, math.Round(2*math.Pi*l.r0/span))
	}
	p := g.prim(l.mat)
	first := uint32(len(p.Positions))
	u := perp(sub(rings[1].at, rings[0].at))
	run := 0.0
	for i, rg := range rings {
		t := normalize(sub(rings[min(i+1, len(rings)-1)].at, rings[max(i-1, 0)].at))
		if i > 0 {
			run += length(sub(rg.at, rings[i-1].at))
		}
		u = normalize(sub(u, scale(t, dot(u, t))))
		w := cross(t, u)
		for j := 0; j <= sides; j++ {
			a := 2 * math.Pi * float64(j) / float64(sides)
			n := add(scale(u, math.Cos(a)), scale(w, math.Sin(a)))
			k := 1.0
			if l.ribs > 0 {
				k = 1 - l.rib*(1-math.Abs(math.Cos(float64(l.ribs)*a/2)))
			}
			p.Positions = append(p.Positions, f32(add(rg.at, scale(n, rg.r*k))))
			p.Normals = append(p.Normals, f32(n))
			p.UVs = append(p.UVs, [2]float32{float32(float64(j) / float64(sides) * tiles), float32(run / span)})
			if i > 0 && j > 0 {
				q := first + uint32((i-1)*(sides+1)+j-1)
				row := uint32(sides + 1)
				p.Indices = append(p.Indices, q, q+1, q+row+1, q, q+row+1, q+row)
			}
		}
	}
}
