package tree

// crown is how a broadleaf tree branches: how many generations of
// boughs, how many children a bough has, and the trunk, the angle they leave it at,
// how long and thick they are for their parent, how they wobble and
// rise (or sag) along their length, what they are made of, and the
// foliage on the last.
type crown struct {
	depth   int
	kids    [2]int
	fork    [2]int // the children of the trunk
	spread  [2]float64
	shrink  float64
	thin    float64
	wobble  float64
	lift    float64
	flare   float64
	sides   int
	bark    int
	span    float64
	foliage foliage
}

// foliage is how a twig is covered: the material of its cards, how many
// there are to a twig, how long they are, how wide for their length, the
// share of the twig's end they sit on, and the middle of the crown, from
// which the cards are lit as facing outward.
type foliage struct {
	mat    int
	n      [2]int
	size   [2]float64
	aspect float64
	reach  float64
	centre vec
}

// bough is a bough to grow: from where, heading where, how long and how
// thick at its start, and its generation.
type bough struct {
	from, dir vec
	length, r float64
	depth     int
}

// grow grows bough b and, on and on, its children.
func (g *grower) grow(c *crown, b bough) {
	path := g.course(course{b.from, b.dir, b.length, 6, c.wobble, c.lift})
	flare := 0.0
	if b.depth == 0 {
		flare = c.flare
	}
	g.tube(limb{path: path, r0: b.r, r1: b.r * 0.5, flare: flare, sides: c.sides, mat: c.bark, span: c.span})
	if b.depth >= c.depth {
		g.foliate(c.foliage, path)
		return
	}
	kids := c.kids
	if b.depth == 0 {
		kids = c.fork
	}
	n := g.count(kids[0], kids[1])
	turn := g.between(0, 6.28)
	for i := range n {
		t := g.between(0.45, 1)
		if b.depth == 0 {
			t = g.between(0.85, 1)
		}
		at, dir := along(path, t)
		d := leave(cone{dir, g.between(c.spread[0], c.spread[1]), turn + 6.2832*float64(i)/float64(n) + g.between(-0.4, 0.4)})
		g.grow(c, bough{at, d, b.length * c.shrink * g.between(0.8, 1.1), b.r * (1 - 0.5*t) * c.thin, b.depth + 1})
	}
}

// foliate covers the end of twig path with cards of f.
func (g *grower) foliate(f foliage, path []vec) {
	for range g.count(f.n[0], f.n[1]) {
		at, dir := along(path, g.between(1-f.reach, 1))
		heading := normalize(add(dir, scale(g.anywhere(), 0.9)))
		l := g.between(f.size[0], f.size[1])
		g.card(card{at: at, along: heading, across: normalize(cross(heading, g.anywhere())), length: l, width: l * f.aspect, mat: f.mat,
			normal: add(sub(at, f.centre), scale(up, 0.3))})
	}
}
