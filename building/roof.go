package building

import "math"

// thickness is how thick each roof covering is: tiles on their battens,
// slates, a deep coat of thatch.
var thickness = map[Material]float64{Tiles: 0.14, Slates: 0.09, Thatch: 0.38}

// pitched is the shape of a roof over walls: their width along X and
// depth along Z, the height of the eaves, the slope's tangent and the
// overhang past the walls.
type pitched struct {
	w, d, eaves float64
	slope, over float64
}

// ridge is how high the underside of the roof's ridge is, over the
// middle of walls, depth deep, rising from their tops.
func (p pitched) ridge() float64 { return p.eaves + p.d/2*p.slope }

// roof is the house's roof, hipped or gabled.
func (r *raising) roof() {
	h := r.h.Roof
	p := pitched{w: r.h.Width, d: r.h.Depth, eaves: r.eaves, slope: math.Tan(h.Pitch * math.Pi / 180), over: h.Overhang}
	if h.Hipped {
		r.b.hip(cover{p, paint{of: h.Cover}})
		return
	}
	r.b.gable(cover{p, paint{of: h.Cover}})
	r.b.gableEnds(cover{p, r.wall()})
}

// cover is a roof's shape and what it is covered with.
type cover struct {
	pitched
	paint paint
}

// gable adds a gabled roof's two slopes, meeting at the ridge, along X,
// reaching past the walls' ends and eaves, and a row of ridge tiles
// capping them.
func (b *builder) gable(c cover) {
	p := c.pitched
	t := thickness[c.paint.of] * math.Sqrt(1+p.slope*p.slope) // its thickness upright
	x := p.w/2 + p.over*0.7
	ridge := p.ridge()
	low := p.eaves - p.over*p.slope
	for _, z := range []float64{1, -1} {
		edge := z * (p.d/2 + p.over)
		var c8 [8]vec
		for k, sx := range []float64{-1, 1} {
			c8[2*k] = vec{sx * x, ridge, 0}
			c8[2*k+1] = vec{sx * x, low, edge}
			c8[4+2*k] = vec{sx * x, ridge + t, 0}
			c8[4+2*k+1] = vec{sx * x, low + t, edge}
		}
		// The loops must run round: ridge, eave, eave, ridge.
		c8[2], c8[3] = c8[3], c8[2]
		c8[6], c8[7] = c8[7], c8[6]
		b.convex(hexa(c8, c.paint))
	}
	capR := map[Material]float64{Thatch: 0.2, Slates: 0.08}[c.paint.of]
	if capR == 0 {
		capR = 0.11
	}
	b.turn(column{a: vec{-x, ridge + t - capR*0.4, 0}, b: vec{x, ridge + t - capR*0.4, 0}, ra: capR, rb: capR, sides: 8, paint: c.paint})
}

// gableEnds are the triangles of wall, of c's paint, under the ends of
// a gabled roof of c's shape.
func (b *builder) gableEnds(c cover) {
	p := c.pitched
	w, d := p.w/2, p.d/2
	s := solid{paint: c.paint, points: []vec{
		{-w, p.eaves, -d}, {-w, p.eaves, d}, {-w, p.ridge(), 0},
		{w, p.eaves, -d}, {w, p.eaves, d}, {w, p.ridge(), 0},
	}, faces: [][]int{{0, 1, 2}, {3, 4, 5}, {0, 1, 4, 3}, {1, 2, 5, 4}, {2, 0, 3, 5}}}
	b.convex(s)
}

// hip adds a hipped roof: its eaves a band all round, past the walls,
// its four slopes rising from them to a ridge along the longer side, or
// to a point over a square.
func (b *builder) hip(c cover) {
	p := c.pitched
	t := thickness[c.paint.of] * math.Sqrt(1+p.slope*p.slope)
	hx, hz := p.w/2+p.over, p.d/2+p.over
	low := p.eaves - p.over*p.slope
	top := low + t
	run := min(hx, hz)
	rise := top + run*p.slope
	s := solid{paint: c.paint}
	for _, y := range []float64{low, top} {
		s.points = append(s.points, vec{-hx, y, -hz}, vec{hx, y, -hz}, vec{hx, y, hz}, vec{-hx, y, hz})
	}
	s.faces = [][]int{{0, 1, 2, 3}, {0, 1, 5, 4}, {1, 2, 6, 5}, {2, 3, 7, 6}, {3, 0, 4, 7}}
	rx, rz := hx-run, hz-run
	if rx < 0.01 && rz < 0.01 {
		s.points = append(s.points, vec{0, rise, 0})
		s.faces = append(s.faces, []int{4, 5, 8}, []int{5, 6, 8}, []int{6, 7, 8}, []int{7, 4, 8})
		b.convex(s)
		return
	}
	// The ridge's ends, 8 and 9: along X, 8 to the left; along Z, 8 at the back.
	s.points = append(s.points, vec{-rx, rise, -rz}, vec{rx, rise, rz})
	if rx >= rz {
		s.faces = append(s.faces, []int{4, 5, 9, 8}, []int{5, 6, 9}, []int{6, 7, 8, 9}, []int{7, 4, 8})
	} else {
		s.faces = append(s.faces, []int{4, 5, 8}, []int{5, 6, 9, 8}, []int{6, 7, 9}, []int{7, 4, 8, 9})
	}
	b.convex(s)
}
