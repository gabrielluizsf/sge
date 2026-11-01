//go:build linux && amd64

package gl

import (
	"errors"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/render3d"
)

// Device is an OpenGL 4.2 device: its context, on a thread of its own,
// and the resources it holds, by number.
type Device struct {
	th     *thread
	egl    *egl
	gl     procs
	w, h   int
	window bool // drawing into a window's own framebuffer, not screen's

	buffers   []buffer // by number; 0 is none
	textures  []texture
	pipelines []pipeline
	programs  [render3d.Unlit + 1]uint32

	screen   texture // offscreen, the screen: a colour and a depth target
	white    uint32  // the texture of a draw with none
	far      uint32  // the shadow map of a draw with none: all farthest
	material uint32  // the uniforms of a draw with no material
	vao      uint32
	fbos     map[[2]render3d.Texture]uint32

	queue  []cmd     // what is recorded, for the thread
	floats []float32 // what Writes recorded write
	run    func()    // runs the queue, on the thread
	st     state     // what the thread has set, as it runs the queue
}

var _ render3d.Device = (*Device)(nil)

// Offscreen is a device drawing into a screen of its own, w by h, in
// memory: for tests, and pictures taken without a window.
func Offscreen(w, h int) (*Device, error) {
	return open(func(d *Device) error {
		if err := d.egl.offscreen(); err != nil {
			return err
		}
		return d.setup()
	}, w, h)
}

// Window is a device drawing into win, which it takes the surface of
// (win's canvas no longer copied in): on X11 only. Present shows each
// frame drawn, before win.End.
func Window(win *antui.Window) (*Device, error) {
	h, ok := win.NativeHandle()
	if !ok || h.System != antui.NativeX11 {
		return nil, errors.New("gl: the window is not an X11 one")
	}
	d, err := open(func(d *Device) error {
		xdpy, err := openX11(h.Display)
		if err != nil {
			return err
		}
		if err := d.egl.window(xdpy, h.Window); err != nil {
			return err
		}
		d.window = true
		return d.setup()
	}, win.Width(), win.Height())
	if err == nil {
		win.SetExternal(true)
	}
	return d, err
}

// open is a device whose context start makes current, on its thread,
// its screen w by h.
func open(start func(d *Device) error, w, h int) (*Device, error) {
	e, err := openEGL()
	if err != nil {
		return nil, err
	}
	d := &Device{th: startThread(), egl: e, buffers: []buffer{{}}, textures: []texture{{}}, pipelines: []pipeline{{}},
		fbos: map[[2]render3d.Texture]uint32{}}
	d.run = d.runQueue
	d.w, d.h = max(w, 1), max(h, 1)
	d.th.do(func() { err = start(d) })
	if err != nil {
		d.th.stop()
		return nil, err
	}
	return d, nil
}

// Present shows the frame drawn on the window, or, offscreen, waits for
// the GPU to finish it.
func (d *Device) Present() {
	d.do(func() {
		if d.window {
			call(d.egl.swapBuffers, n(d.egl.dpy), n(d.egl.surf))
			return
		}
		call(d.gl.finish)
	})
}

// Close lets go of the device: its context, and all it holds.
func (d *Device) Close() {
	d.do(d.egl.close)
	d.th.stop()
}

// do runs what is recorded, then job, on the device's thread. A job may
// bind textures of its own: the next draw binds its own again.
func (d *Device) do(job func()) {
	d.th.do(func() {
		d.runQueue()
		job()
		d.st.bound = false
	})
}

// Size is the screen's size.
func (d *Device) Size() (w, h int) { return d.w, d.h }

// Resize makes the screen w by h: offscreen, cleared; a window's follows
// the window, and passes into it are drawn that size.
func (d *Device) Resize(w, h int) {
	d.w, d.h = max(w, 1), max(h, 1)
	if !d.window {
		d.do(d.newScreen)
	}
}
