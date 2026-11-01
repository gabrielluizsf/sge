//go:build linux && amd64

package gl

import (
	"fmt"
	"unsafe"

	"github.com/gabrielluizsf/sge/render3d"
)

// buffer is a buffer's GL name, its kind, and how many floats, or
// indices, it holds.
type buffer struct {
	name uint32
	kind render3d.BufferKind
	size int
}

// targets are the GL target each kind of buffer is made on.
var targets = [...]uint32{render3d.Vertices: glArrayBuffer, render3d.Indices: glElementArrayBuffer, render3d.Uniforms: glUniformBuffer}

// NewBuffer is a buffer holding a copy of b's data, or b.Size floats of
// zeros.
func (d *Device) NewBuffer(b render3d.BufferDesc) (render3d.Buffer, error) {
	if b.Kind > render3d.Uniforms {
		return 0, fmt.Errorf("gl: no buffer kind %d", b.Kind)
	}
	buf := buffer{kind: b.Kind}
	var data unsafe.Pointer
	switch {
	case b.Kind == render3d.Indices:
		buf.size = len(b.Indices)
		if buf.size > 0 {
			data = unsafe.Pointer(&b.Indices[0])
		}
	case b.Floats != nil:
		buf.size = len(b.Floats)
		if buf.size > 0 {
			data = unsafe.Pointer(&b.Floats[0])
		}
	default:
		if b.Size < 0 {
			return 0, fmt.Errorf("gl: a buffer of %d floats", b.Size)
		}
		buf.size = b.Size
	}
	usage := uint32(glStaticDraw)
	if b.Kind == render3d.Uniforms {
		usage = glDynamicDraw
	}
	d.do(func() {
		g := &d.gl
		buf.name = gen(g.genBuffers)
		t := targets[b.Kind]
		call(g.bindBuffer, n(t), n(buf.name))
		// GL makes a buffer of no bytes, but not one bound to a range.
		call(g.bufferData, n(t), n(4*max(buf.size, 1)), ptr(data), n(usage))
		if b.Kind == render3d.Indices {
			// Indices bind to the vertex array: leave it as it was.
			call(g.bindBuffer, n(t), n(d.bufferName(d.st.idx)))
		}
		d.buffers = append(d.buffers, buf)
	})
	return render3d.Buffer(len(d.buffers) - 1), nil
}

// Write copies data into b's floats from at on, as many as fit, when the
// passes before it are drawn: data may be changed as soon as it returns.
func (d *Device) Write(b render3d.Buffer, at int, data []float32) {
	from := len(d.floats)
	d.floats = append(d.floats, data...)
	d.queue = append(d.queue, cmd{op: opWrite, buf: b, at: at, from: from, count: len(data)})
}

// write is Write, on the thread: what Write recorded from c.from on.
func (d *Device) write(c *cmd) {
	buf := d.buffer(c.buf)
	if buf == nil || c.at < 0 || c.at > buf.size || buf.kind == render3d.Indices {
		return
	}
	count := min(c.count, buf.size-c.at)
	if count == 0 {
		return
	}
	g := &d.gl
	t := targets[buf.kind]
	call(g.bindBuffer, n(t), n(buf.name))
	call(g.bufferSubData, n(t), n(4*c.at), n(4*count), ptr(unsafe.Pointer(&d.floats[c.from])))
}

// buffer is buffer b, or nil.
func (d *Device) buffer(b render3d.Buffer) *buffer {
	if int(b) >= len(d.buffers) || d.buffers[b].name == 0 {
		return nil
	}
	return &d.buffers[b]
}

// bufferName is b's GL name, 0 for none.
func (d *Device) bufferName(b render3d.Buffer) uint32 {
	if buf := d.buffer(b); buf != nil {
		return buf.name
	}
	return 0
}
