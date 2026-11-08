package tree

import "math"

// Pine is the pine of the floresta: a tall, straight trunk, whorls of
// branches all the way up it, shorter toward the top so the whole is a
// cone, carrying sprays of needles. Sapling 1.5 to 3 m, young 6 to 10,
// mature 15 to 20.
var Pine = Species{Name: "pine", Metres: [Sizes][2]float64{{1.5, 3}, {6, 10}, {15, 20}}, grow: growPine}

// pine is a pine being grown: its trunk and what it is made of.
type pine struct {
	g       *grower
	trunk   []vec
	r0, r1  float64
	bark    int
	needles int
}

// growPine grows a pine.
func growPine(g *grower) {
	h := g.size
	p := &pine{g: g, bark: g.bark(g.barks.Pine), needles: g.cutout("needles", needles(uint32(g.between(1, 9999)))), r0: 0.016*h + 0.012, r1: 0.002*h + 0.005}
	p.trunk = g.course(course{vec{}, vec{g.between(-0.03, 0.03), 1, g.between(-0.03, 0.03)}, h, 12, 0.012, 0})
	g.tube(limb{path: p.trunk, r0: p.r0, r1: p.r1, flare: 0.6, sides: 9, mat: p.bark, span: g.barks.Pine.Span})
	first := [Sizes]float64{0.05, 0.12, 0.2}[g.stage]
	n := max(5, min(int(h*0.9)+4, 22))
	for k := range n {
		p.whorl(first + (0.96-first)*float64(k)/float64(n-1))
	}
	for i := range 6 {
		at, _ := along(p.trunk, 1)
		d := leave(cone{up, 0.4, 2 * math.Pi * float64(i) / 6})
		g.card(card{at: at, along: d, across: perp(d), length: p.cardLength(), width: p.cardLength() * 0.8, mat: p.needles, normal: add(d, up)})
	}
}

// cardLength is how long a spray of needles is.
func (p *pine) cardLength() float64 { return 0.12 + 0.045*p.g.size }

// whorl grows a ring of branches at f of the trunk's height.
func (p *pine) whorl(f float64) {
	g := p.g
	first := [Sizes]float64{0.05, 0.12, 0.2}[g.stage]
	rel := (f - first) / (1 - first)
	reach := g.size * 0.19 * (1 - 0.93*rel)
	at, _ := along(p.trunk, f)
	n := g.count(4, 6)
	turn := g.between(0, 2*math.Pi)
	for i := range n {
		yaw := turn + 2*math.Pi*float64(i)/float64(n) + g.between(-0.3, 0.3)
		pitch := -0.25 + 0.7*rel + g.between(-0.1, 0.1)
		dir := vec{math.Cos(yaw) * math.Cos(pitch), math.Sin(pitch), math.Sin(yaw) * math.Cos(pitch)}
		path := g.course(course{at, dir, reach * g.between(0.85, 1.1), 5, 0.06, -0.035 + 0.045*rel})
		r := p.r0 + (p.r1-p.r0)*f
		g.tube(limb{path: path, r0: math.Max(0.006, r*0.32), r1: 0.003, sides: 5, mat: p.bark, span: g.barks.Pine.Span})
		p.spray(path)
	}
}

// spray covers branch path with sprays of needles, a flat one and an
// upright one at each place.
func (p *pine) spray(path []vec) {
	g := p.g
	n := 3 + int(length(sub(path[len(path)-1], path[0]))*2.2)
	for j := range n {
		at, d := along(path, 0.3+0.7*float64(j)/float64(n-1))
		l := p.cardLength() * g.between(0.8, 1.1)
		flat := normalize(cross(d, up))
		lit := add(normalize(vec{d[0], 0, d[2]}), scale(up, 0.6))
		g.card(card{at: at, along: d, across: flat, length: l, width: l * 0.9, mat: p.needles, normal: lit})
		g.card(card{at: at, along: d, across: normalize(cross(d, flat)), length: l, width: l * 0.9, mat: p.needles, normal: lit})
	}
}
