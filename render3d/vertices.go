package render3d

import (
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// prim is a primitive to lay out: its normals, its own or worked out, and
// how many joints its skin has, 0 for a rigid one.
type prim struct {
	p       *gltf.Primitive
	normals [][3]float32
	joints  int
}

// appendVertices is dst with pr's vertices after it, in Static's layout,
// or Skinned's when pr has joints.
func appendVertices(dst []float32, pr prim) []float32 {
	p := pr.p
	for i, pos := range p.Positions {
		n := pr.normals[i]
		var uv [2]float32
		if i < len(p.UVs) {
			uv = p.UVs[i]
		}
		dst = append(dst, pos[0], pos[1], pos[2], n[0], n[1], n[2], uv[0], uv[1])
		if pr.joints > 0 {
			js, ws := weighed(p.Joints[i], p.Weights[i], pr.joints)
			dst = append(dst, js[0], js[1], js[2], js[3], ws[0], ws[1], ws[2], ws[3])
		}
	}
	return dst
}

// weighed is a vertex's joints and weights made sound: a joint past the
// skin's, or a weight below 0, left out, and the rest adding up to 1; a
// vertex with none follows the first joint.
func weighed(js [4]uint16, ws [4]float32, joints int) (j, w [4]float32) {
	var total float32
	for k := range 4 {
		if int(js[k]) < joints && ws[k] > 0 && !math.IsNaN(float64(ws[k])) {
			j[k], w[k] = float32(js[k]), ws[k]
			total += ws[k]
		}
	}
	if total == 0 {
		return [4]float32{}, [4]float32{1}
	}
	for k := range w {
		w[k] /= total
	}
	return j, w
}

// normalsOf are p's normals, or, when it has none, each vertex's worked
// out from the triangles round it, larger ones counting more.
func normalsOf(p *gltf.Primitive) [][3]float32 {
	if len(p.Normals) == len(p.Positions) {
		return p.Normals
	}
	sum := make([][3]float64, len(p.Positions))
	for t := 0; t+2 < len(p.Indices); t += 3 {
		a, b, c := f64(p.Positions[p.Indices[t]]), f64(p.Positions[p.Indices[t+1]]), f64(p.Positions[p.Indices[t+2]])
		n := cross(sub(b, a), sub(c, a))
		for k := range 3 {
			v := &sum[p.Indices[t+k]]
			v[0], v[1], v[2] = v[0]+n[0], v[1]+n[1], v[2]+n[2]
		}
	}
	out := make([][3]float32, len(sum))
	for i, n := range sum {
		n = unit(n)
		out[i] = [3]float32{float32(n[0]), float32(n[1]), float32(n[2])}
	}
	return out
}

func f64(v [3]float32) [3]float64 { return [3]float64{float64(v[0]), float64(v[1]), float64(v[2])} }
