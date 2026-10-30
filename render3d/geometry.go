package render3d

import (
	"fmt"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/pose"
)

// geometry is a model's vertices and indices gathered for upload, and
// its points in the rest pose, for its bounds.
type geometry struct {
	static, skinned []float32
	indices         []uint32
	points          [][3]float64
}

// gather lays out the primitives of every node of m that shows a mesh,
// in m's parts.
func (g *geometry) gather(m *Model) error {
	src := m.src
	palettes := make([][]pose.Mat, len(src.Skins))
	for s, sk := range src.Skins {
		if len(sk.Joints) > MaxObjects {
			return fmt.Errorf("render3d: skin %d has %d joints, more than %d", s, len(sk.Joints), MaxObjects)
		}
		for _, j := range sk.Joints {
			if j < 0 || j >= len(src.Nodes) {
				return fmt.Errorf("render3d: skin %d has no node %d", s, j)
			}
		}
		palettes[s] = m.poser.Joints(s, nil)
	}
	for ni, n := range src.Nodes {
		if n.Mesh < 0 || n.Mesh >= len(src.Meshes) || m.rest[ni][15] != 1 {
			continue // no mesh, or not in the scene
		}
		for pi := range src.Meshes[n.Mesh].Primitives {
			p := &src.Meshes[n.Mesh].Primitives[pi]
			if err := g.add(m, primAt{p, ni, palettes}); err != nil {
				return fmt.Errorf("render3d: mesh %d primitive %d: %w", n.Mesh, pi, err)
			}
		}
	}
	return nil
}

// primAt is a primitive of node, with the skins' rest palettes.
type primAt struct {
	p        *gltf.Primitive
	node     int
	palettes [][]pose.Mat
}

// add lays out one primitive as a part of m.
func (g *geometry) add(m *Model, at primAt) error {
	p := at.p
	nv := len(p.Positions)
	for _, i := range p.Indices {
		if int(i) >= nv {
			return fmt.Errorf("index %d of %d vertices", i, nv)
		}
	}
	if len(p.Indices) < 3 || nv == 0 {
		return nil
	}
	skin := m.src.Nodes[at.node].Skin
	if skin < 0 || skin >= len(at.palettes) || len(at.palettes[skin]) == 0 || len(p.Joints) != nv || len(p.Weights) != nv {
		skin = -1
	}
	pr := prim{p: p, normals: normalsOf(p)}
	dst := &g.static
	stride := StaticFloats
	if skin >= 0 {
		pr.joints, dst, stride = len(at.palettes[skin]), &g.skinned, SkinnedFloats
	}
	base := uint32(len(*dst) / stride)
	*dst = appendVertices(*dst, pr)
	first := len(g.indices)
	for _, i := range p.Indices[:len(p.Indices)/3*3] {
		g.indices = append(g.indices, base+i)
	}
	mat := 0
	if p.Material >= 0 && p.Material < len(m.src.Materials) {
		mat = p.Material + 1
	}
	m.parts = append(m.parts, part{node: at.node, skin: skin, first: first, count: len(g.indices) - first, material: mat})
	for i, pos := range p.Positions {
		place := m.rest[at.node]
		if skin >= 0 {
			place = blend(at.palettes[skin], p.Joints[i], p.Weights[i])
		}
		g.points = append(g.points, place.Point(f64(pos)))
	}
	return nil
}

// blend is the joints' matrices mixed by a vertex's weights.
func blend(joints []pose.Mat, js [4]uint16, ws [4]float32) pose.Mat {
	j, w := weighed(js, ws, len(joints))
	var m pose.Mat
	for k := range 4 {
		for i := range m {
			m[i] += joints[int(j[k])][i] * float64(w[k])
		}
	}
	return m
}

// bounds is the ball round the points: the middle of their box, out to
// the farthest, grown by margin.
func (g *geometry) bounds(margin float64) Sphere {
	if len(g.points) == 0 {
		return Sphere{}
	}
	lo, hi := g.points[0], g.points[0]
	for _, p := range g.points {
		for k := range 3 {
			lo[k], hi[k] = min(lo[k], p[k]), max(hi[k], p[k])
		}
	}
	c := [3]float64{(lo[0] + hi[0]) / 2, (lo[1] + hi[1]) / 2, (lo[2] + hi[2]) / 2}
	var r float64
	for _, p := range g.points {
		d := sub(p, c)
		r = max(r, dot(d, d))
	}
	return Sphere{Center: c, Radius: math.Sqrt(r) * margin}
}
