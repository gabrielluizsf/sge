package human

import "math"

// skinOf is a worn garment's surface, as fitted on the body: its
// triangles and, at each vertex, which way is out of the body.
type skinOf struct {
	at      [][3]float64
	tris    [][3]int32
	normals [][3]float64
}

// layerReach is how near, in the body's units, a vertex must be to the
// garment under it to be pushed out of it: farther, it is another limb's.
const layerReach = 0.04 * metre

// layerClear is how far out of the garment under it a vertex is left.
const layerClear = 0.0012 * metre

// underClear is how far under the garment over it a covered vertex is drawn:
// the garment's faces run straight between its vertices, so what lies just
// under it would still come through between them.
const underClear = 0.006 * metre

// coverSlant is how far to the side a surface may be and still cover a
// vertex: the square of the distance over the square of the height.
const coverSlant = 4

// grid finds the triangles of a skinOf near a point.
type grid struct {
	s     *skinOf
	cell  float64
	cells map[int64][]int32
}

// newGrid is the grid of s, each triangle in every cell it is within
// layerReach of.
func newGrid(s *skinOf) *grid {
	g := &grid{s: s, cell: 2 * layerReach, cells: map[int64][]int32{}}
	for ti, t := range s.tris {
		lo, hi := s.at[t[0]], s.at[t[0]]
		for _, v := range t[1:] {
			for k := range 3 {
				lo[k], hi[k] = min(lo[k], s.at[v][k]), max(hi[k], s.at[v][k])
			}
		}
		var a, b [3]int32
		for k := range 3 {
			a[k], b[k] = g.index(lo[k]-layerReach), g.index(hi[k]+layerReach)
		}
		for x := a[0]; x <= b[0]; x++ {
			for y := a[1]; y <= b[1]; y++ {
				for z := a[2]; z <= b[2]; z++ {
					k := cellKey(x, y, z)
					g.cells[k] = append(g.cells[k], int32(ti))
				}
			}
		}
	}
	return g
}

// cellKey is the cell x, y, z as one number.
func cellKey(x, y, z int32) int64 {
	return int64(x&0x1fffff)<<42 | int64(y&0x1fffff)<<21 | int64(z&0x1fffff)
}

func (g *grid) index(v float64) int32 { return int32(math.Floor(v / g.cell)) }

// closest is the point of s nearest q, the way out there, and how far, if
// it is within layerReach.
func (g *grid) closest(q [3]float64) (on, out [3]float64, ok bool) {
	best := layerReach * layerReach
	for _, ti := range g.cells[cellKey(g.index(q[0]), g.index(q[1]), g.index(q[2]))] {
		t := g.s.tris[ti]
		p, u, v := nearestOnTriangle(q, g.s.at[t[0]], g.s.at[t[1]], g.s.at[t[2]])
		d := sub(q, p)
		if d2 := d[0]*d[0] + d[1]*d[1] + d[2]*d[2]; d2 < best {
			best, ok, on = d2, true, p
			w := 1 - u - v
			out = add(add(scale(g.s.normals[t[0]], w), scale(g.s.normals[t[1]], u)), scale(g.s.normals[t[2]], v))
		}
	}
	return on, normalize(out), ok
}

func dot(a, b [3]float64) float64 { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }

// nearestOnTriangle is the point of triangle abc nearest p, and its
// weights u of b and v of c (Ericson, Real-Time Collision Detection).
func nearestOnTriangle(p, a, b, c [3]float64) (pt [3]float64, u, v float64) {
	ab, ac, ap := sub(b, a), sub(c, a), sub(p, a)
	d1, d2 := dot(ab, ap), dot(ac, ap)
	if d1 <= 0 && d2 <= 0 {
		return a, 0, 0
	}
	bp := sub(p, b)
	d3, d4 := dot(ab, bp), dot(ac, bp)
	if d3 >= 0 && d4 <= d3 {
		return b, 1, 0
	}
	if vc := d1*d4 - d3*d2; vc <= 0 && d1 >= 0 && d3 <= 0 {
		t := d1 / (d1 - d3)
		return add(a, scale(ab, t)), t, 0
	}
	cp := sub(p, c)
	d5, d6 := dot(ab, cp), dot(ac, cp)
	if d6 >= 0 && d5 <= d6 {
		return c, 0, 1
	}
	if vb := d5*d2 - d1*d6; vb <= 0 && d2 >= 0 && d6 <= 0 {
		t := d2 / (d2 - d6)
		return add(a, scale(ac, t)), 0, t
	}
	if va := d3*d6 - d5*d4; va <= 0 && d4-d3 >= 0 && d5-d6 >= 0 {
		t := (d4 - d3) / ((d4 - d3) + (d5 - d6))
		return add(b, scale(sub(c, b), t)), 1 - t, t
	}
	vc := d1*d4 - d3*d2
	vb := d5*d2 - d1*d6
	va := d3*d6 - d5*d4
	s := 1 / (va + vb + vc)
	v, w := vb*s, vc*s
	return add(a, add(scale(ab, v), scale(ac, w))), v, w
}

// pushOut moves the vertices at, of a garment outside s, out of it where
// they are inside it or less than layerClear from it.
func (g *grid) pushOut(at [][3]float64) {
	for i, q := range at {
		// A push can bring a vertex to another triangle, whose way out is
		// another: it is looked at again, a few times.
		for range 4 {
			on, out, ok := g.closest(q)
			if !ok {
				break
			}
			h := dot(sub(q, on), out)
			if h >= layerClear {
				break
			}
			q = add(q, scale(out, layerClear-h))
		}
		at[i] = q
	}
}

// pullIn draws under this surface each of the vertices at that it covers and
// that lie out through it or close under it, so what is under a garment stays
// under it. Those beyond its edge, seen past its hem, are left where they are.
func (g *grid) pullIn(at [][3]float64) {
	for i, q := range at {
		// A pull can bring a vertex under another face, whose way out is
		// another: it is looked at again, a few times.
		for range 4 {
			on, out, ok := g.closest(q)
			if !ok || !covers(q, on, out) {
				break
			}
			h := dot(sub(q, on), out)
			if h <= -underClear {
				break
			}
			q = add(q, scale(out, -underClear-h))
		}
		at[i] = q
	}
}

// covers reports whether the surface's nearest point on to q, whose way out
// is out, lies over or under q rather than off to its side, past an edge.
func covers(q, on, out [3]float64) bool {
	d := sub(q, on)
	h := dot(d, out)
	return dot(d, d) <= h*h*coverSlant+1e-9
}

// near is whether the vertices a and b come within layerReach of each
// other's bounds at all.
func near(a, b [][3]float64) bool {
	box := func(v [][3]float64) (lo, hi [3]float64) {
		lo, hi = v[0], v[0]
		for _, p := range v {
			for k := range 3 {
				lo[k], hi[k] = min(lo[k], p[k]), max(hi[k], p[k])
			}
		}
		return
	}
	alo, ahi := box(a)
	blo, bhi := box(b)
	for k := range 3 {
		if ahi[k]+layerReach < blo[k] || bhi[k]+layerReach < alo[k] {
			return false
		}
	}
	return true
}
