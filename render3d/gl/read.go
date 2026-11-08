//go:build linux && amd64

package gl

import (
	"fmt"
	"image"
	"math"
	"unsafe"

	"github.com/gabrielluizsf/sge/render3d"
)

// Read is t's pixels, top row first: its colours as they are stored,
// sRGB for the screen, or, for a texture of depth, its depths as greys,
// white farthest.
func (d *Device) Read(t render3d.Texture) (*image.NRGBA, error) {
	var img *image.NRGBA
	var err error
	d.do(func() { img, err = d.read(t) })
	return img, err
}

// read is Read, on the thread.
func (d *Device) read(t render3d.Texture) (*image.NRGBA, error) {
	f, depth := frame{0, d.w, d.h}, false
	if t != render3d.Screen || !d.window {
		tx := d.texture(t)
		if tx == nil {
			return nil, fmt.Errorf("gl: no texture %d", t)
		}
		key := [2]render3d.Texture{t, 0}
		if tx.depth {
			key = [2]render3d.Texture{render3d.Screen, t}
		}
		fbo, ok := d.fbo(key)
		if !ok {
			return nil, fmt.Errorf("gl: texture %d cannot be read", t)
		}
		f, depth = frame{fbo, tx.w, tx.h}, tx.depth
	}
	g := &d.gl
	call(g.bindFramebuffer, n(glFramebuffer), n(f.fbo))
	img := image.NewNRGBA(image.Rect(0, 0, f.w, f.h))
	if depth {
		z := make([]float32, f.w*f.h)
		call(g.readPixels, n(0), n(0), n(f.w), n(f.h), n(glDepthComponent), n(glFloat), ptr(unsafe.Pointer(&z[0])))
		for i, v := range z {
			grey := byte(math.Round(float64(v) * 255))
			copy(img.Pix[4*i:], []byte{grey, grey, grey, 255})
		}
	} else {
		call(g.readPixels, n(0), n(0), n(f.w), n(f.h), n(glRGBA), n(glUnsignedByte), ptr(unsafe.Pointer(&img.Pix[0])))
	}
	flip(img)
	return img, nil
}

// flip turns img upside down: GL's rows run up, from the bottom.
func flip(img *image.NRGBA) {
	h, stride := img.Rect.Dy(), img.Stride
	row := make([]byte, stride)
	for y := range h / 2 {
		a, b := img.Pix[y*stride:(y+1)*stride], img.Pix[(h-1-y)*stride:(h-y)*stride]
		copy(row, a)
		copy(a, b)
		copy(b, row)
	}
}
