package soft

import (
	"math"

	"github.com/gabrielluizsf/sge/render3d"
)

// frag is a pixel of a triangle: the weights of its corners there, in
// perspective, and whether its back faces the camera.
type frag struct {
	tri  *[3]vout
	w    [3]float64
	back bool
}

// plot shades the pixel i as f, and is whether it is drawn: not when
// its alpha is under the material's cutoff.
func (p *pass) plot(s *state, f frag, i int) bool {
	var world, normal [3]float64
	var uv [2]float64
	for k, v := range f.tri {
		for c := range 3 {
			world[c] += v.world[c] * f.w[k]
			normal[c] += v.normal[c] * f.w[k]
		}
		uv[0] += v.uv[0] * f.w[k]
		uv[1] += v.uv[1] * f.w[k]
	}
	albedo := s.color
	if s.tex != nil {
		t := s.tex.at(uv)
		for k := range albedo {
			albedo[k] *= t[k]
		}
	}
	if s.cutoff > 0 && albedo[3] < s.cutoff {
		return false
	}
	if p.color == nil || s.kind == render3d.Shadow || s.kind == render3d.ShadowSkinned {
		return true
	}
	out := albedo
	if s.kind != render3d.Unlit {
		n := unit(normal)
		if f.back {
			n = [3]float64{-n[0], -n[1], -n[2]}
		}
		sun := math.Max(0, n[0]*s.toward[0]+n[1]*s.toward[1]+n[2]*s.toward[2]) * s.lit(world)
		for k := range 3 {
			out[k] = albedo[k] * (s.ambient[k] + s.sun[k]*sun)
		}
	}
	c := p.color[4*i : 4*i+4]
	c[0], c[1], c[2], c[3] = float32(out[0]), float32(out[1]), float32(out[2]), 1
	return true
}

// lit is 1 where the sun reaches world, 0 where its shadow map has
// something nearer the sun, by more than the bias.
func (s *state) lit(world [3]float64) float64 {
	if s.shadow == nil {
		return 1
	}
	c := s.light.apply(world, 1)
	x, y, z := c[0]/c[3], c[1]/c[3], c[2]/c[3]
	px := int(math.Floor((x + 1) / 2 * float64(s.shadow.w)))
	py := int(math.Floor((1 - y) / 2 * float64(s.shadow.h)))
	if px < 0 || py < 0 || px >= s.shadow.w || py >= s.shadow.h || z > 1 {
		return 1
	}
	if z-s.bias > float64(s.shadow.depth[py*s.shadow.w+px]) {
		return 0
	}
	return 1
}
