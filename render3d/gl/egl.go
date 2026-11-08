//go:build linux && amd64

package gl

import (
	"fmt"

	"github.com/gabrielluizsf/sge/native"
)

// EGL's enums, those the device uses.
const (
	eglNone                = 0x3038
	eglSurfaceType         = 0x3033
	eglRenderableType      = 0x3040
	eglWindowBit           = 0x4
	eglOpenGLBit           = 0x8
	eglRedSize             = 0x3024
	eglGreenSize           = 0x3023
	eglBlueSize            = 0x3022
	eglAlphaSize           = 0x3021
	eglDepthSize           = 0x3025
	eglOpenGLAPI           = 0x30A2
	eglContextMajor        = 0x3098
	eglContextMinor        = 0x30FB
	eglContextProfile      = 0x30FD
	eglCoreProfile         = 0x1
	eglPlatformX11         = 0x31D5
	eglPlatformSurfaceless = 0x31DD // EGL_PLATFORM_SURFACELESS_MESA
	eglColorspace          = 0x309D // EGL_GL_COLORSPACE_KHR
	eglColorspaceSRGB      = 0x3089
)

// egl is libEGL's functions, and the display, configuration, context and
// surface a device draws with.
type egl struct {
	getDisplay, getPlatformDisplay, initialize, bindAPI native.Proc
	chooseConfig, createContext, createWindowSurface    native.Proc
	makeCurrent, swapBuffers, swapInterval, getProcAddr native.Proc
	destroySurface, destroyContext, getError            native.Proc
	dpy, cfg, ctx, surf                                 uintptr
}

// openEGL is libEGL, its functions looked up.
func openEGL() (*egl, error) {
	lib, err := native.Open("libEGL.so.1")
	if err != nil {
		return nil, err
	}
	e := &egl{}
	for _, f := range []named{
		{"eglGetDisplay", &e.getDisplay}, {"eglGetPlatformDisplay", &e.getPlatformDisplay},
		{"eglInitialize", &e.initialize}, {"eglBindAPI", &e.bindAPI}, {"eglChooseConfig", &e.chooseConfig},
		{"eglCreateContext", &e.createContext}, {"eglCreateWindowSurface", &e.createWindowSurface},
		{"eglMakeCurrent", &e.makeCurrent}, {"eglSwapBuffers", &e.swapBuffers}, {"eglSwapInterval", &e.swapInterval},
		{"eglGetProcAddress", &e.getProcAddr}, {"eglDestroySurface", &e.destroySurface},
		{"eglDestroyContext", &e.destroyContext}, {"eglGetError", &e.getError},
	} {
		if *f.p, err = lib.Symbol(f.name); err != nil {
			return nil, err
		}
	}
	return e, nil
}

// fail is an error saying what failed, with EGL's own error.
func (e *egl) fail(what string) error {
	return fmt.Errorf("gl: %s failed (EGL error %#x)", what, call(e.getError))
}

// attribs is an EGL attribute list, ended.
func attribs(v ...int32) []int32 { return append(v, eglNone) }

// start initialises e's display and makes a GL 4.2 core context on it,
// of a configuration for windows or for none.
func (e *egl) start(window bool) error {
	var major, minor int32
	if call(e.initialize, n(e.dpy), ptr(pointer(&major)), ptr(pointer(&minor))) != 1 {
		return e.fail("eglInitialize")
	}
	if call(e.bindAPI, n(eglOpenGLAPI)) != 1 {
		return e.fail("eglBindAPI")
	}
	want := attribs(eglRenderableType, eglOpenGLBit, eglSurfaceType, 0, eglRedSize, 8, eglGreenSize, 8, eglBlueSize, 8, eglAlphaSize, 8)
	if window {
		want = attribs(eglRenderableType, eglOpenGLBit, eglSurfaceType, eglWindowBit, eglRedSize, 8, eglGreenSize, 8, eglBlueSize, 8, eglAlphaSize, 8, eglDepthSize, 24)
	}
	var count int32
	if call(e.chooseConfig, n(e.dpy), ptr(pointer(&want[0])), ptr(pointer(&e.cfg)), n(1), ptr(pointer(&count))) != 1 || count == 0 {
		return e.fail("eglChooseConfig")
	}
	ctx := attribs(eglContextMajor, 4, eglContextMinor, 2, eglContextProfile, eglCoreProfile)
	if e.ctx = call(e.createContext, n(e.dpy), n(e.cfg), n(0), ptr(pointer(&ctx[0]))); e.ctx == 0 {
		return e.fail("eglCreateContext (OpenGL 4.2 core)")
	}
	return nil
}

// current makes e's context current on this thread, drawing into its
// surface, or into none.
func (e *egl) current() error {
	if call(e.makeCurrent, n(e.dpy), n(e.surf), n(e.surf), n(e.ctx)) != 1 {
		return e.fail("eglMakeCurrent")
	}
	return nil
}

// offscreen makes e a context drawing into no surface, on Mesa's
// surfaceless platform or the default display.
func (e *egl) offscreen() error {
	e.dpy = call(e.getPlatformDisplay, n(eglPlatformSurfaceless), n(0), n(0))
	if e.dpy == 0 || e.start(false) != nil {
		e.dpy = call(e.getDisplay, n(0))
		if e.dpy == 0 {
			return e.fail("eglGetDisplay")
		}
		if err := e.start(false); err != nil {
			return err
		}
	}
	return e.current()
}

// window makes e a context drawing into the X11 window xid of the
// display xdpy, its colours sRGB when EGL can.
func (e *egl) window(xdpy, xid uintptr) error {
	if e.dpy = call(e.getPlatformDisplay, n(eglPlatformX11), n(xdpy), n(0)); e.dpy == 0 {
		return e.fail("eglGetPlatformDisplay (X11)")
	}
	if err := e.start(true); err != nil {
		return err
	}
	srgb := attribs(eglColorspace, eglColorspaceSRGB)
	if e.surf = call(e.createWindowSurface, n(e.dpy), n(e.cfg), n(xid), ptr(pointer(&srgb[0]))); e.surf == 0 {
		e.surf = call(e.createWindowSurface, n(e.dpy), n(e.cfg), n(xid), n(0))
	}
	if e.surf == 0 {
		return e.fail("eglCreateWindowSurface")
	}
	if err := e.current(); err != nil {
		return err
	}
	call(e.swapInterval, n(e.dpy), n(1))
	return nil
}

// close lets go of e's context and surface.
func (e *egl) close() {
	call(e.makeCurrent, n(e.dpy), n(0), n(0), n(0))
	if e.surf != 0 {
		call(e.destroySurface, n(e.dpy), n(e.surf))
	}
	call(e.destroyContext, n(e.dpy), n(e.ctx))
}
