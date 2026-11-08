//go:build linux && amd64

package gl

import (
	"fmt"
	"unsafe"

	"github.com/gabrielluizsf/sge/native"
	"github.com/gabrielluizsf/sge/render3d"
)

// setup readies the context just made current: its functions, the state
// every pass shares, and what a draw binds when given nothing.
func (d *Device) setup() error {
	if err := d.gl.load(d.egl.getProcAddr); err != nil {
		return err
	}
	g := &d.gl
	d.vao = gen(g.genVertexArrays)
	call(g.bindVertexArray, n(d.vao))
	call(g.enable, n(glDepthTest))
	call(g.depthFunc, n(glLess))
	call(g.frontFace, n(glCCW))
	call(g.cullFace, n(glBack))
	call(g.enable, n(glFramebufferSRGB))
	call(g.pixelStorei, n(glUnpackAlignment), n(1))
	call(g.pixelStorei, n(glPackAlignment), n(1))
	d.white = d.newTexture(texImage{w: 1, h: 1, format: glRGBA8, kind: glRGBA, typ: glUnsignedByte, px: unsafe.Pointer(&[4]byte{255, 255, 255, 255})})
	far := float32(1)
	d.far = d.newTexture(texImage{w: 1, h: 1, format: glDepthComponent32F, kind: glDepthComponent, typ: glFloat, px: unsafe.Pointer(&far)})
	m := [render3d.MaterialFloats]float32{1, 1, 1, 1, 0, 1, 0, 0}
	d.material = gen(g.genBuffers)
	call(g.bindBuffer, n(glUniformBuffer), n(d.material))
	call(g.bufferData, n(glUniformBuffer), n(4*len(m)), ptr(unsafe.Pointer(&m)), n(glStaticDraw))
	if !d.window {
		d.newScreen()
	}
	if e := call(g.getError); e != 0 {
		return fmt.Errorf("gl: setting up %s: error %#x", native.GoString(call(g.getString, n(glVersion))), e)
	}
	return nil
}

// gen is a new object of fn's, a glGen function.
func gen(fn native.Proc) uint32 {
	var id uint32
	call(fn, n(1), ptr(pointer(&id)))
	return id
}

// del deletes the object id by fn, a glDelete function.
func del(fn native.Proc, id uint32) {
	if id != 0 {
		call(fn, n(1), ptr(pointer(&id)))
	}
}
