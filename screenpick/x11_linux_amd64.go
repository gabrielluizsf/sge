//go:build linux && amd64

package screenpick

import (
	"errors"
	"os"
	"unsafe"

	"github.com/gabrielluizsf/sge/native"
)

// X11 constants: the left button in the pointer mask, all planes, ZPixmap.
const (
	button1Mask = 1 << 8
	zPixmap     = 2
)

// x11 reads the root window of an X display of its own.
type x11 struct {
	dpy, root                                                    uintptr
	queryPointer, getImage, getPixel, destroyImage, closeDisplay native.Proc
}

// Open connects to the X display. It is ErrUnavailable in a Wayland
// session or with no display, and an error if libX11 is not there.
func Open() (Source, error) {
	if err := usable(os.Getenv); err != nil {
		return nil, err
	}
	lib, err := native.Open("libX11.so.6")
	if err != nil {
		return nil, ErrUnavailable
	}
	x := &x11{}
	syms := map[string]*native.Proc{"XQueryPointer": &x.queryPointer, "XGetImage": &x.getImage,
		"XGetPixel": &x.getPixel, "XDestroyImage": &x.destroyImage, "XCloseDisplay": &x.closeDisplay}
	for name, p := range syms {
		if *p, err = lib.Symbol(name); err != nil {
			return nil, ErrUnavailable
		}
	}
	open, err := lib.Symbol("XOpenDisplay")
	if err != nil {
		return nil, ErrUnavailable
	}
	rootOf, err := lib.Symbol("XDefaultRootWindow")
	if err != nil {
		return nil, ErrUnavailable
	}
	if x.dpy = new(native.Frame).Int(0).Call(open); x.dpy == 0 {
		return nil, ErrUnavailable
	}
	x.root = new(native.Frame).Int(x.dpy).Call(rootOf)
	return x, nil
}

// Close ends the connection.
func (x *x11) Close() {
	if x.dpy != 0 {
		new(native.Frame).Int(x.dpy).Call(x.closeDisplay)
		x.dpy = 0
	}
}

// Read asks where the pointer is and takes the one pixel under it.
func (x *x11) Read() (Reading, error) {
	if x.dpy == 0 {
		return Reading{}, ErrUnavailable
	}
	var root, child uintptr
	var rx, ry, wx, wy int32
	var mask uint32
	ok := new(native.Frame).Int(x.dpy).Int(x.root).
		Ptr(unsafe.Pointer(&root)).Ptr(unsafe.Pointer(&child)).
		Ptr(unsafe.Pointer(&rx)).Ptr(unsafe.Pointer(&ry)).
		Ptr(unsafe.Pointer(&wx)).Ptr(unsafe.Pointer(&wy)).
		Ptr(unsafe.Pointer(&mask)).Call(x.queryPointer)
	if int32(ok) == 0 {
		return Reading{}, errors.New("screenpick: the pointer is on another screen")
	}
	r := Reading{X: int(rx), Y: int(ry), Down: mask&button1Mask != 0}
	var err error
	r.R, r.G, r.B, err = x.pixel(int(rx), int(ry))
	return r, err
}
