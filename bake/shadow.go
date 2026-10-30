package bake

import "math"

// shadowMap is the scene seen from the sun: for each texel, how near the
// sun the nearest surface there is, so a point farther than it is shaded.
type shadowMap struct {
	sun, u, v  [3]float64 // toward the sun, and across its view
	minU, minV float64
	scale      float64 // texels a world unit
	size       int
	depth      []float64
	bias       float64
}

const shadowSize = 1024

// newShadowMap renders tris' depth from the sun.
func newShadowMap(tris []triangle, sun [3]float64) *shadowMap {
	s := &shadowMap{sun: sun, size: shadowSize}
	up := [3]float64{0, 0, 1}
	if math.Abs(dot(sun, up)) > 0.99 {
		up = [3]float64{0, 1, 0}
	}
	s.u = normalize(cross(up, sun))
	s.v = cross(sun, s.u)
	s.minU, s.minV = math.Inf(1), math.Inf(1)
	maxU, maxV := math.Inf(-1), math.Inf(-1)
	for _, t := range tris {
		for _, v := range t.v {
			a, b := dot(v.p, s.u), dot(v.p, s.v)
			s.minU, s.minV, maxU, maxV = min(s.minU, a), min(s.minV, b), max(maxU, a), max(maxV, b)
		}
	}
	extent := max(maxU-s.minU, maxV-s.minV, 1e-6)
	s.minU -= extent * 0.02
	s.minV -= extent * 0.02
	s.scale = float64(s.size-1) / (extent * 1.04)
	s.bias = 2 / s.scale
	s.depth = make([]float64, s.size*s.size)
	for i := range s.depth {
		s.depth[i] = math.Inf(-1)
	}
	for _, t := range tris {
		var p [3][3]float64
		for k, v := range t.v {
			p[k] = s.project(v.p)
		}
		fill(p, s.size, s.size, func(x, y int, w [3]float64) {
			d := w[0]*p[0][2] + w[1]*p[1][2] + w[2]*p[2][2]
			if i := y*s.size + x; d > s.depth[i] {
				s.depth[i] = d
			}
		})
	}
	return s
}

// project is p in the map's texels, and its depth toward the sun.
func (s *shadowMap) project(p [3]float64) [3]float64 {
	return [3]float64{(dot(p, s.u) - s.minU) * s.scale, (dot(p, s.v) - s.minV) * s.scale, dot(p, s.sun)}
}

// lit is how much of the sun reaches p, from 0 to 1: the share of the
// nine texels around it with nothing nearer the sun, so shadows' edges
// are soft.
func (s *shadowMap) lit(p [3]float64) float64 {
	q := s.project(p)
	x, y := int(q[0]), int(q[1])
	open := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			xi, yi := x+dx, y+dy
			if xi < 0 || yi < 0 || xi >= s.size || yi >= s.size || s.depth[yi*s.size+xi] <= q[2]+s.bias {
				open++
			}
		}
	}
	return float64(open) / 9
}
