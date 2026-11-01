//go:build linux && amd64

package gl

import "github.com/gabrielluizsf/sge/render3d"

// state is what the thread has set as it runs a pass's commands.
type state struct {
	drawing    bool // the pass has a framebuffer to draw into
	pipe       *pipeline
	verts, idx render3d.Buffer
	laid       uint32 // the GL buffer the vertex arrays read
	laidStride int    // and its vertices' floats
	bind       render3d.Bindings
	bound      bool // bind is bound
	ok         bool // and can be drawn with
}

// begin starts a pass into p's framebuffer, cleared as p says, with
// nothing set.
func (d *Device) begin(p *render3d.PassDesc) {
	f, ok := d.frameOf(p)
	d.st = state{drawing: ok, laid: d.st.laid, laidStride: d.st.laidStride}
	call(d.gl.bindBuffer, n(glElementArrayBuffer), n(0))
	if !ok {
		return
	}
	g := &d.gl
	call(g.bindFramebuffer, n(glFramebuffer), n(f.fbo))
	call(g.viewport, n(0), n(0), n(f.w), n(f.h))
	var mask uint32
	if p.Clear && !p.DepthOnly {
		c := p.ClearColor
		callFloats(g.clearColor, c[0], c[1], c[2], c[3])
		mask |= glColorBufferBit
	}
	if p.ClearDepth {
		callFloats(g.clearDepthf, 1)
		mask |= glDepthBufferBit
	}
	if mask != 0 {
		call(g.clear, n(mask))
	}
}

// setPipeline uses pipeline pl's program, culling the faces it does not
// draw.
func (d *Device) setPipeline(pl render3d.Pipeline) {
	p := d.pipeline(pl)
	d.st.pipe = p
	if p == nil {
		return
	}
	g := &d.gl
	call(g.useProgram, n(p.program))
	if p.double {
		call(g.disable, n(glCullFace))
	} else {
		call(g.enable, n(glCullFace))
	}
}

// layout points the vertex arrays at buf's vertices, stride floats each:
// position, normal and texture coordinates, then joints and weights.
func (d *Device) layout(buf *buffer, stride int) {
	g := &d.gl
	call(g.bindBuffer, n(glArrayBuffer), n(buf.name))
	sizes := [...]int{3, 3, 2, 4, 4}
	used := 3
	if stride == render3d.SkinnedFloats {
		used = 5
	}
	at := 0
	for i, size := range sizes {
		if i >= used {
			call(g.disableVertexAttribArray, n(i))
			continue
		}
		call(g.enableVertexAttribArray, n(i))
		call(g.vertexAttribPointer, n(i), n(size), n(glFloat), n(0), n(4*stride), n(4*at))
		at += size
	}
	d.st.laid, d.st.laidStride = buf.name, stride
}

// drawElements draws dr with what is set, when all of it is there.
func (d *Device) drawElements(dr render3d.Draw) {
	s := &d.st
	if !s.drawing || s.pipe == nil || dr.First < 0 || dr.Count < 0 {
		return
	}
	vb, ib := d.buffer(s.verts), d.buffer(s.idx)
	if vb == nil || ib == nil || ib.kind != render3d.Indices {
		return
	}
	if !s.bound {
		s.ok, s.bound = d.bindAll(), true
	}
	if !s.ok {
		return
	}
	if s.laid != vb.name || s.laidStride != s.pipe.stride {
		d.layout(vb, s.pipe.stride)
	}
	count := (min(dr.First+dr.Count, ib.size) - dr.First) / 3 * 3
	if count <= 0 {
		return
	}
	call(d.gl.drawElementsInstanced, n(glTriangles), n(count), n(glUnsignedInt), n(4*dr.First), n(max(dr.Instances, 1)))
}
