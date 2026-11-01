//go:build linux && amd64

package gl

import "github.com/gabrielluizsf/sge/render3d"

// texture is texture t, the offscreen screen for Screen, or nil.
func (d *Device) texture(t render3d.Texture) *texture {
	if t == render3d.Screen {
		if d.window {
			return nil
		}
		return &d.screen
	}
	if int(t) >= len(d.textures) || d.textures[t].name == 0 {
		return nil
	}
	return &d.textures[t]
}

// fbo is the framebuffer drawing into the colour texture key[0] (none
// when it is Screen and key[1] is not, the screen when both are) and the
// depth texture key[1] (when 0, the colour texture's own renderbuffer),
// made the first time; false when GL cannot draw into them.
func (d *Device) fbo(key [2]render3d.Texture) (uint32, bool) {
	if f, ok := d.fbos[key]; ok {
		return f, true
	}
	var color, depth *texture
	if key[0] != render3d.Screen || key[1] == 0 {
		if color = d.texture(key[0]); color == nil {
			return 0, false
		}
	}
	if key[1] != 0 {
		if depth = d.texture(key[1]); depth == nil {
			return 0, false
		}
	}
	g := &d.gl
	f := gen(g.genFramebuffers)
	call(g.bindFramebuffer, n(glFramebuffer), n(f))
	if color != nil {
		call(g.framebufferTexture2D, n(glFramebuffer), n(glColorAttachment0), n(glTexture2D), n(color.name), n(0))
		call(g.drawBuffer, n(glColorAttachment0))
		call(g.readBuffer, n(glColorAttachment0))
	} else {
		call(g.drawBuffer, n(glNone))
		call(g.readBuffer, n(glNone))
	}
	if depth != nil {
		call(g.framebufferTexture2D, n(glFramebuffer), n(glDepthAttachment), n(glTexture2D), n(depth.name), n(0))
	} else {
		if color.rb == 0 {
			color.rb = gen(g.genRenderbuffers)
			call(g.bindRenderbuffer, n(glRenderbuffer), n(color.rb))
			call(g.renderbufferStorage, n(glRenderbuffer), n(glDepthComponent24), n(color.w), n(color.h))
		}
		call(g.framebufferRenderbuffer, n(glFramebuffer), n(glDepthAttachment), n(glRenderbuffer), n(color.rb))
	}
	if call(g.checkFramebufferStatus, n(glFramebuffer)) != glFramebufferDone {
		del(g.deleteFramebuffers, f)
		return 0, false
	}
	d.fbos[key] = f
	return f, true
}

// forget deletes the framebuffers drawing into t: for Screen, the
// screen's own.
func (d *Device) forget(t render3d.Texture) {
	for key, f := range d.fbos {
		if t == render3d.Screen && key == [2]render3d.Texture{} || t != render3d.Screen && (key[0] == t || key[1] == t) {
			del(d.gl.deleteFramebuffers, f)
			delete(d.fbos, key)
		}
	}
}

// frame is the framebuffer a pass draws into, and its size.
type frame struct {
	fbo  uint32
	w, h int
}

// frameOf is the framebuffer p draws into, and whether there is one: not
// when its textures are released, or not of the kinds it needs.
func (d *Device) frameOf(p *render3d.PassDesc) (frame, bool) {
	switch {
	case p.DepthOnly:
		t := d.texture(p.Depth)
		if p.Depth == render3d.Screen || t == nil || !t.depth {
			return frame{}, false
		}
		f, ok := d.fbo([2]render3d.Texture{render3d.Screen, p.Depth})
		return frame{f, t.w, t.h}, ok
	case p.Target == render3d.Screen && d.window:
		return frame{0, d.w, d.h}, true
	case p.Target == render3d.Screen:
		f, ok := d.fbo([2]render3d.Texture{})
		return frame{f, d.w, d.h}, ok
	}
	t := d.texture(p.Target)
	if t == nil || t.depth {
		return frame{}, false
	}
	depth := p.Depth
	if dt := d.texture(depth); depth == render3d.Screen || dt == nil || !dt.depth || dt.w != t.w || dt.h != t.h {
		depth = 0
	}
	f, ok := d.fbo([2]render3d.Texture{p.Target, depth})
	return frame{f, t.w, t.h}, ok
}
