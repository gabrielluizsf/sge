package building

import "math"

// column is a round solid, sides sided, from a to b, its radius ra at a
// and rb at b: a post, a log, a bucket, a cone nearly to a point.
type column struct {
	a, b   vec
	ra, rb float64
	sides  int
	paint  paint
}

// turn adds column c, its grain along it.
func (b *builder) turn(c column) {
	axis := normalize(sub(c.b, c.a))
	u := cross(axis, vec{1, 0, 0})
	if length(u) < 0.5 {
		u = cross(axis, vec{0, 0, 1})
	}
	u = normalize(u)
	w := cross(axis, u)
	n := max(c.sides, 3)
	s := solid{paint: c.paint, grain: axis}
	bottom, top := make([]int, n), make([]int, n)
	for i := range n {
		a := 2 * math.Pi * (float64(i) + 0.5) / float64(n)
		round := add(scale(u, math.Cos(a)), scale(w, math.Sin(a)))
		s.points = append(s.points, add(c.a, scale(round, c.ra)), add(c.b, scale(round, c.rb)))
		bottom[i], top[i] = 2*i, 2*i+1
	}
	s.faces = append(s.faces, bottom, top)
	for i := range n {
		j := (i + 1) % n
		s.faces = append(s.faces, []int{2 * i, 2 * j, 2*j + 1, 2*i + 1})
	}
	b.convex(s)
}

// ball is a round solid about c, radius r, squashed to squash of its
// width tall: a fruit, a stone.
type ball struct {
	c      vec
	r      float64
	squash float64
	paint  paint
}

// roll adds ball x, of rings rings from pole to pole.
func (b *builder) roll(x ball) {
	const rings, segs = 6, 10
	s := solid{paint: x.paint}
	s.points = append(s.points, add(x.c, vec{0, -x.r * x.squash, 0}))
	for r := 1; r < rings; r++ {
		th := math.Pi * float64(r) / rings
		for k := range segs {
			ph := 2 * math.Pi * float64(k) / segs
			s.points = append(s.points, add(x.c, vec{x.r * math.Sin(th) * math.Cos(ph), -x.r * x.squash * math.Cos(th), x.r * math.Sin(th) * math.Sin(ph)}))
		}
	}
	top := len(s.points)
	s.points = append(s.points, add(x.c, vec{0, x.r * x.squash, 0}))
	at := func(r, k int) int { return 1 + (r-1)*segs + k%segs }
	for k := range segs {
		s.faces = append(s.faces, []int{0, at(1, k), at(1, k+1)}, []int{top, at(rings-1, k+1), at(rings-1, k)})
		for r := 1; r+1 < rings; r++ {
			s.faces = append(s.faces, []int{at(r, k), at(r, k+1), at(r+1, k+1), at(r+1, k)})
		}
	}
	b.convex(s)
}
