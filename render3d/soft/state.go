package soft

import "github.com/gabrielluizsf/sge/render3d"

// mat is a 4 by 4 matrix, column-major.
type mat [16]float64

// state is what a draw reads, gathered from the pass's bindings.
type state struct {
	kind    render3d.Kind
	double  bool
	stride  int
	vp      mat // the camera's, or the sun's for a shadow
	light   mat
	toward  [3]float64
	bias    float64
	sun     [3]float64
	ambient [3]float64
	color   [4]float64
	cutoff  float64
	objects []float32
	tex     *texture
	shadow  *texture
}

// state is the draw's state, and whether it can draw: a pipeline, its
// buffers and its uniforms all set.
func (p *pass) state() (*state, bool) {
	frame := p.uniforms(p.bind.Frame)
	objects := p.uniforms(p.bind.Objects)
	if p.pipe == nil || p.verts == nil || p.idx == nil || len(frame) < render3d.FrameFloats || len(objects) < render3d.ObjectFloats {
		return nil, false
	}
	s := &state{kind: p.pipe.Kind, double: p.pipe.DoubleSided, stride: render3d.StaticFloats, objects: objects, color: [4]float64{1, 1, 1, 1}}
	if s.kind == render3d.Skinned || s.kind == render3d.ShadowSkinned {
		s.stride = render3d.SkinnedFloats
	}
	s.vp, s.light = matAt(frame, render3d.FrameViewProj), matAt(frame, render3d.FrameLightViewProj)
	if s.kind == render3d.Shadow || s.kind == render3d.ShadowSkinned {
		s.vp = s.light
	}
	f := func(i int) float64 { return float64(frame[i]) }
	s.toward = [3]float64{f(render3d.FrameSun), f(render3d.FrameSun + 1), f(render3d.FrameSun + 2)}
	s.bias = f(render3d.FrameSun + 3)
	s.sun = [3]float64{f(render3d.FrameSunColor), f(render3d.FrameSunColor + 1), f(render3d.FrameSunColor + 2)}
	s.ambient = [3]float64{f(render3d.FrameAmbient), f(render3d.FrameAmbient + 1), f(render3d.FrameAmbient + 2)}
	if m := p.uniforms(p.bind.Material); len(m) >= render3d.MaterialFloats {
		for k := range 4 {
			s.color[k] = float64(m[render3d.MaterialColor+k])
		}
		s.cutoff = float64(m[render3d.MaterialParams+2])
	}
	if t := p.d.texture(p.bind.Texture); p.bind.Texture != render3d.Screen && t != nil && t.color != nil {
		s.tex = t
	}
	if t := p.d.texture(p.bind.Shadow); p.bind.Shadow != render3d.Screen && t != nil && t.depth != nil {
		s.shadow = t
	}
	return s, true
}

// matAt is the matrix of the 16 floats of f from at on.
func matAt(f []float32, at int) mat {
	var m mat
	for i := range m {
		m[i] = float64(f[at+i])
	}
	return m
}

// apply is m times the point (or, with w 0, direction) p.
func (m *mat) apply(p [3]float64, w float64) [4]float64 {
	return [4]float64{
		m[0]*p[0] + m[4]*p[1] + m[8]*p[2] + m[12]*w,
		m[1]*p[0] + m[5]*p[1] + m[9]*p[2] + m[13]*w,
		m[2]*p[0] + m[6]*p[1] + m[10]*p[2] + m[14]*w,
		m[3]*p[0] + m[7]*p[1] + m[11]*p[2] + m[15]*w,
	}
}
