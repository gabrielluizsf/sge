package soft

import "github.com/gabrielluizsf/sge/render3d"

// pass is a pass drawing into its colour (nil for depth only) and depth,
// with what is set.
type pass struct {
	d     *Device
	w, h  int
	color []float32
	depth []float32
	srgb  bool

	pipe  *render3d.PipelineDesc
	verts *buffer
	idx   *buffer
	bind  render3d.Bindings
}

// Begin starts a pass into p's targets, cleared as p says; a pass whose
// targets do not exist draws nothing.
func (d *Device) Begin(p render3d.PassDesc) render3d.Pass {
	ps := &pass{d: d}
	if p.DepthOnly {
		t := d.texture(p.Depth)
		if t == nil || t.depth == nil {
			return ps
		}
		ps.w, ps.h, ps.depth = t.w, t.h, t.depth
	} else {
		t := d.texture(p.Target)
		if t == nil {
			return ps
		}
		t.target()
		ps.w, ps.h, ps.color, ps.depth, ps.srgb = t.w, t.h, t.color, t.depth, t.srgb
		if dt := d.texture(p.Depth); p.Depth != render3d.Screen && dt != nil && len(dt.depth) == t.w*t.h {
			ps.depth = dt.depth
		}
	}
	if p.Clear && ps.color != nil {
		for i := 0; i < len(ps.color); i += 4 {
			copy(ps.color[i:i+4], p.ClearColor[:])
		}
	}
	if p.ClearDepth {
		for i := range ps.depth {
			ps.depth[i] = 1
		}
	}
	return ps
}

func (p *pass) SetPipeline(pl render3d.Pipeline) {
	p.pipe = nil
	if int(pl) < len(p.d.pipelines) {
		p.pipe = p.d.pipelines[pl]
	}
}

func (p *pass) SetVertices(b render3d.Buffer)   { p.verts = p.d.buffer(b) }
func (p *pass) SetIndices(b render3d.Buffer)    { p.idx = p.d.buffer(b) }
func (p *pass) SetBindings(b render3d.Bindings) { p.bind = b }
func (p *pass) End()                            {}

// uniforms are the floats of r, or nil when r is not inside its buffer.
func (p *pass) uniforms(r render3d.Range) []float32 {
	b := p.d.buffer(r.Buffer)
	if b == nil || r.At < 0 || r.Count < 0 || r.At+r.Count > len(b.floats) {
		return nil
	}
	return b.floats[r.At : r.At+r.Count]
}
