//go:build linux && amd64

package gl

import "github.com/gabrielluizsf/sge/render3d"

// The bytes each uniform block declares.
const (
	frameBytes    = 4 * render3d.FrameFloats
	materialBytes = 4 * render3d.MaterialFloats
	objectsBytes  = 4 * render3d.ObjectFloats * render3d.MaxObjects
)

// block is a range of a uniform buffer to bind: its GL name, and where
// it starts and how long it is, in bytes.
type block struct {
	name       uint32
	at, length int
}

// blockOf is the range r as a block of at least need floats, as much of
// declared bytes as its buffer holds from r.At on; false when r is not
// inside its buffer, or shorter than need.
func (d *Device) blockOf(r render3d.Range, need, declared int) (block, bool) {
	buf := d.buffer(r.Buffer)
	if buf == nil || buf.kind != render3d.Uniforms || r.At < 0 || r.Count < need || r.At+r.Count > buf.size {
		return block{}, false
	}
	return block{buf.name, 4 * r.At, max(4*r.Count, min(declared, 4*(buf.size-r.At)))}, true
}

// bindAll binds what the set bindings say a draw reads, and is whether
// it can draw: not without its frame's and objects' uniforms. A draw
// with no material is white; with no texture, untextured; with no
// shadow map, all in the sun.
func (d *Device) bindAll() bool {
	b := &d.st.bind
	frame, ok := d.blockOf(b.Frame, render3d.FrameFloats, frameBytes)
	if !ok {
		return false
	}
	objects, ok := d.blockOf(b.Objects, render3d.ObjectFloats, objectsBytes)
	if !ok {
		return false
	}
	material, ok := d.blockOf(b.Material, render3d.MaterialFloats, materialBytes)
	if !ok {
		material = block{d.material, 0, materialBytes}
	}
	g := &d.gl
	for i, bl := range [...]block{frame, material, objects} {
		call(g.bindBufferRange, n(glUniformBuffer), n(i), n(bl.name), n(bl.at), n(bl.length))
	}
	base := d.white
	if t := d.texture(b.Texture); b.Texture != render3d.Screen && t != nil && !t.depth {
		base = t.name
	}
	shadow := d.far
	if t := d.texture(b.Shadow); b.Shadow != render3d.Screen && t != nil && t.depth {
		shadow = t.name
	}
	call(g.activeTexture, n(glTexture0+1))
	call(g.bindTexture, n(glTexture2D), n(shadow))
	call(g.activeTexture, n(glTexture0))
	call(g.bindTexture, n(glTexture2D), n(base))
	return true
}
