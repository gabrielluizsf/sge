package pose

import "github.com/gabrielluizsf/sge/gltf"

// Poser poses one model, as many times a frame as it is shown, reusing
// its memory: a frame's poses cost no allocations.
type Poser struct {
	m       *gltf.Model
	order   []int // the nodes reached from the roots, parents first
	parent  []int // each node's parent, -1 for a root
	rest    []trs // each node's place as the model had it when the poser was made
	local   []trs
	world   []Mat
	lengths []float64
	tracks  [][]track
}

// New is a poser of m.
func New(m *gltf.Model) *Poser {
	p := &Poser{m: m, parent: make([]int, len(m.Nodes)), local: make([]trs, len(m.Nodes)), world: make([]Mat, len(m.Nodes))}
	p.rest = make([]trs, len(m.Nodes))
	for i, n := range m.Nodes {
		p.rest[i] = trs{n.Translation, n.Rotation, n.Scale}
	}
	for i := range p.parent {
		p.parent[i] = -1
	}
	seen := make([]bool, len(m.Nodes))
	queue := []int{}
	for _, r := range m.Roots {
		if r >= 0 && r < len(m.Nodes) && !seen[r] {
			seen[r] = true
			queue = append(queue, r)
		}
	}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		p.order = append(p.order, n)
		for _, c := range m.Nodes[n].Children {
			if c >= 0 && c < len(m.Nodes) && !seen[c] {
				seen[c], p.parent[c] = true, n
				queue = append(queue, c)
			}
		}
	}
	for i := range m.Animations {
		p.lengths = append(p.lengths, Length(&m.Animations[i]))
		p.tracks = append(p.tracks, tracksOf(&m.Animations[i], len(m.Nodes)))
	}
	return p
}

// Animation is the index of m's animation called name, or -1.
func (p *Poser) Animation(name string) int {
	for i := range p.m.Animations {
		if name != "" && p.m.Animations[i].Name == name {
			return i
		}
	}
	return -1
}

// At is every node's matrix with base before the model's own, as
// animation anim (-1 for none) moves them t seconds in, looped round its
// length. The matrices are the poser's, good until its next pose.
func (p *Poser) At(anim int, t float64, base Mat) []Mat {
	copy(p.local, p.rest)
	if anim >= 0 && anim < len(p.m.Animations) {
		p.animate(anim, t)
	}
	for _, i := range p.order {
		l := &p.local[i]
		own := TRS(l.t, l.r, l.s)
		up := &base
		if p.parent[i] >= 0 {
			up = &p.world[p.parent[i]]
		}
		MulAffine(&p.world[i], up, &own)
	}
	return p.world
}

// trs is a node's place in its parent: translation, rotation, scale.
type trs struct {
	t [3]float64
	r [4]float64
	s [3]float64
}

// animate sets the nodes' local places to animation anim's at t.
func (p *Poser) animate(anim int, t float64) {
	if l := p.lengths[anim]; l > 0 {
		t -= l * float64(int(t/l))
		if t < 0 {
			t += l
		}
	}
	a := &p.m.Animations[anim]
	tracks := p.tracks[anim]
	for c := range tracks {
		tr := &tracks[c]
		if tr.path == none {
			continue
		}
		ch := &a.Channels[c]
		v := tr.sample(ch, t)
		l := &p.local[ch.Node]
		switch tr.path {
		case translation:
			l.t = [3]float64{v[0], v[1], v[2]}
		case rotation:
			l.r = v
		case scale:
			l.s = [3]float64{v[0], v[1], v[2]}
		}
	}
}

// Joints fills out with skin's joints' matrices as the last pose placed
// them, each its node's matrix times its inverse bind matrix, and is
// out's first as many as the skin has joints.
func (p *Poser) Joints(skin int, out []Mat) []Mat {
	sk := &p.m.Skins[skin]
	out = out[:0]
	for j, node := range sk.Joints {
		out = append(out, p.world[node])
		if j < len(sk.InverseBind) {
			MulAffine(&out[j], &p.world[node], (*Mat)(&sk.InverseBind[j]))
		}
	}
	return out
}
