//go:build linux && amd64

package screenpick

import (
	"errors"
	"unsafe"

	"github.com/gabrielluizsf/sge/native"
)

// Offsets in C's XImage (Xlib.h) on 64-bit Linux.
const (
	imgWidth    = 0
	imgHeight   = 4
	imgBitsPP   = 48
	imgRedMask  = 56
	imgGreenMsk = 64
	imgBlueMask = 72
)

// cfield is the C integer of type T at offset off in the C struct at p.
func cfield[T int32 | uint64](p uintptr, off uintptr) T {
	base := *(*unsafe.Pointer)(unsafe.Pointer(&p))
	return *(*T)(unsafe.Add(base, off))
}

// pixel is the colour of the pixel at x, y of the root window. It asks X
// for an image of one pixel, never more, reads that pixel, and frees it.
func (x *x11) pixel(px, py int) (r, g, b uint8, err error) {
	if px < 0 || py < 0 || px > 1<<16 || py > 1<<16 {
		return 0, 0, 0, errors.New("screenpick: the pointer is off the screen")
	}
	img := new(native.Frame).Int(x.dpy).Int(x.root).Int(uintptr(px)).Int(uintptr(py)).
		Int(1).Int(1).Int(^uintptr(0)).Int(zPixmap).Call(x.getImage)
	if img == 0 {
		return 0, 0, 0, errors.New("screenpick: the screen could not be read")
	}
	defer new(native.Frame).Int(img).Call(x.destroyImage)
	if cfield[int32](img, imgWidth) != 1 || cfield[int32](img, imgHeight) != 1 || cfield[int32](img, imgBitsPP) < 16 {
		return 0, 0, 0, errors.New("screenpick: an image of an unexpected shape")
	}
	px32 := uint64(new(native.Frame).Int(img).Int(0).Int(0).Call(x.getPixel))
	return channelOf(px32, cfield[uint64](img, imgRedMask)), channelOf(px32, cfield[uint64](img, imgGreenMsk)),
		channelOf(px32, cfield[uint64](img, imgBlueMask)), nil
}
