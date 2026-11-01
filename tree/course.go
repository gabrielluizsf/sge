package tree

// course is a curve to grow: from from, heading dir, length long in
// steps straight pieces, each turning a random wobble radians and
// rising lift (falling if it is under 0) of a unit of direction.
type course struct {
	from, dir    vec
	length       float64
	steps        int
	wobble, lift float64
}

// course is the points of curve c.
func (g *grower) course(c course) []vec {
	path := []vec{c.from}
	dir := normalize(c.dir)
	step := c.length / float64(c.steps)
	for range c.steps {
		dir = add(dir, scale(g.anywhere(), c.wobble))
		dir[1] += c.lift
		dir = normalize(dir)
		path = append(path, add(path[len(path)-1], scale(dir, step)))
	}
	return path
}

// along is the point t of the way along path, 0 to 1 by its points, and
// the direction it runs there.
func along(path []vec, t float64) (vec, vec) {
	f := max(0, min(t, 1)) * float64(len(path)-1)
	i := min(int(f), len(path)-2)
	d := sub(path[i+1], path[i])
	return add(path[i], scale(d, f-float64(i))), normalize(d)
}
