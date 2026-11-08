//go:build linux && amd64

package gl

import "github.com/gabrielluizsf/sge/render3d"

// Release lets go of a buffer, texture or pipeline: its number is not
// given again, and using it draws nothing.
func (d *Device) Release(h render3d.Handle) {
	d.do(func() {
		g := &d.gl
		switch h := h.(type) {
		case render3d.Buffer:
			if b := d.buffer(h); b != nil {
				if b.name == d.st.laid {
					d.st.laid = 0
				}
				del(g.deleteBuffers, b.name)
				*b = buffer{}
			}
		case render3d.Texture:
			if t := d.texture(h); t != nil && h != render3d.Screen {
				d.forget(h)
				del(g.deleteTextures, t.name)
				del(g.deleteRenderbuffers, t.rb)
				*t = texture{}
			}
		case render3d.Pipeline:
			if p := d.pipeline(h); p != nil {
				*p = pipeline{} // its program is its kind's, kept
			}
		}
	})
}
