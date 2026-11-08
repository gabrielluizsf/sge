package building

// solid is a convex solid: its corners, its faces, each a loop of
// corners' indices in order round it, either way; its paint, and its
// grain, as a face's.
type solid struct {
	points []vec
	faces  [][]int
	paint  paint
	grain  vec
}

// convex adds solid s, closed, each face turned to face out of it: away
// from its middle, as every face of a convex solid does.
func (b *builder) convex(s solid) {
	mid := centroid(s.points)
	for _, loop := range s.faces {
		corners := make([]vec, len(loop))
		for k, i := range loop {
			corners[k] = s.points[i]
		}
		if dot(newell(corners), sub(centroid(corners), mid)) < 0 {
			for l, r := 0, len(corners)-1; l < r; l, r = l+1, r-1 {
				corners[l], corners[r] = corners[r], corners[l]
			}
		}
		b.face(face{corners: corners, paint: s.paint, grain: s.grain})
	}
}

// hexaFaces are a hexahedron's faces: its corners are a loop, 0 to 3,
// and the loop over it, 4 to 7, each over its own.
var hexaFaces = [][]int{{0, 1, 2, 3}, {4, 5, 6, 7}, {0, 1, 5, 4}, {1, 2, 6, 5}, {2, 3, 7, 6}, {3, 0, 4, 7}}

// hexa is a solid of six faces from its eight corners, as hexaFaces
// order them.
func hexa(corners [8]vec, p paint) solid {
	return solid{points: corners[:], faces: hexaFaces, paint: p}
}

// box is a solid from lo to hi, square to the axes.
type box struct {
	lo, hi vec
	paint  paint
	grain  vec
}

// block adds box x.
func (b *builder) block(x box) {
	lo, hi := x.lo, x.hi
	s := hexa([8]vec{
		{lo[0], lo[1], lo[2]}, {hi[0], lo[1], lo[2]}, {hi[0], lo[1], hi[2]}, {lo[0], lo[1], hi[2]},
		{lo[0], hi[1], lo[2]}, {hi[0], hi[1], lo[2]}, {hi[0], hi[1], hi[2]}, {lo[0], hi[1], hi[2]},
	}, x.paint)
	s.grain = x.grain
	b.convex(s)
}

// taper is a block narrowing or widening as it rises: from its foot's
// middle at, h tall, base and top half as wide along X and as deep
// along Z as each says.
type taper struct {
	at        vec
	h         float64
	base, top [2]float64
	paint     paint
}

// raiseTaper adds taper t.
func (b *builder) raiseTaper(t taper) {
	var c [8]vec
	for k, half := range [][2]float64{t.base, t.top} {
		y := t.at[1] + float64(k)*t.h
		x, z := t.at[0], t.at[2]
		c[4*k+0] = vec{x - half[0], y, z - half[1]}
		c[4*k+1] = vec{x + half[0], y, z - half[1]}
		c[4*k+2] = vec{x + half[0], y, z + half[1]}
		c[4*k+3] = vec{x - half[0], y, z + half[1]}
	}
	b.convex(hexa(c, t.paint))
}

// beam is a squared timber from a to b, w wide and h deep across it:
// its width level, its depth square to that, its grain along it.
type beam struct {
	a, b  vec
	w, h  float64
	paint paint
}

// lay adds beam x.
func (b *builder) lay(x beam) {
	axis := normalize(sub(x.b, x.a))
	side := cross(axis, up)
	if length(side) < 1e-6 {
		side = vec{1, 0, 0}
	}
	side = scale(normalize(side), x.w/2)
	rise := scale(normalize(cross(side, axis)), x.h/2)
	var c [8]vec
	for k, end := range []vec{x.a, x.b} {
		c[4*k+0] = add(add(end, side), rise)
		c[4*k+1] = sub(add(end, side), rise)
		c[4*k+2] = sub(sub(end, side), rise)
		c[4*k+3] = add(sub(end, side), rise)
	}
	s := hexa(c, x.paint)
	s.grain = axis
	b.convex(s)
}
