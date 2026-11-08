//go:build linux && amd64

package gl

import (
	"fmt"
	"unsafe"

	"github.com/gabrielluizsf/sge/render3d"
)

// texture is a texture's GL name and what it is: its size, whether it
// holds depth or sRGB colours, and the depth renderbuffer it is drawn
// into with when a pass gives it none.
type texture struct {
	name  uint32
	w, h  int
	depth bool
	srgb  bool
	rb    uint32
}

// texImage is a texture to make: w by h, of GL's format, its pixels at
// px (nil for none) of kind and typ, with mipmaps or not.
type texImage struct {
	w, h              int
	format, kind, typ uint32
	px                unsafe.Pointer
	mipmaps           bool
}

// NewTexture is a texture of desc's size and pixels: sRGB ones decoded
// when sampled, mipmapped when asked.
func (d *Device) NewTexture(desc render3d.TextureDesc) (render3d.Texture, error) {
	w, h := desc.Width, desc.Height
	if w <= 0 || h <= 0 || w > 1<<14 || h > 1<<14 {
		return 0, fmt.Errorf("gl: a texture %d by %d", w, h)
	}
	if desc.Pixels != nil && len(desc.Pixels) < 4*w*h {
		return 0, fmt.Errorf("gl: %d bytes for a texture %d by %d", len(desc.Pixels), w, h)
	}
	img := texImage{w: w, h: h, format: glRGBA8, kind: glRGBA, typ: glUnsignedByte}
	t := texture{w: w, h: h, srgb: desc.SRGB}
	switch {
	case desc.Format == render3d.Depth:
		img.format, img.kind, img.typ, t.depth, t.srgb = glDepthComponent32F, glDepthComponent, glFloat, true, false
	case desc.Pixels != nil:
		img.px, img.mipmaps = unsafe.Pointer(&desc.Pixels[0]), desc.Mipmaps
	}
	if t.srgb {
		img.format = glSRGB8Alpha8
	}
	d.do(func() {
		t.name = d.newTexture(img)
		d.textures = append(d.textures, t)
	})
	return render3d.Texture(len(d.textures) - 1), nil
}

// newTexture is the GL texture img describes: colours repeated past its
// edges, and sampled at their nearest pixel up close; depths sampled
// as they are.
func (d *Device) newTexture(img texImage) uint32 {
	g := &d.gl
	t := gen(g.genTextures)
	call(g.bindTexture, n(glTexture2D), n(t))
	call(g.texImage2D, n(glTexture2D), n(0), n(img.format), n(img.w), n(img.h), n(0), n(img.kind), n(img.typ), ptr(img.px))
	minf, wrap := uint32(glNearest), uint32(glRepeat)
	if img.kind == glDepthComponent {
		wrap = glClampToEdge
	}
	if img.mipmaps {
		call(g.generateMipmap, n(glTexture2D))
		minf = glLinearMipmapLinear
	}
	call(g.texParameteri, n(glTexture2D), n(glTexMinFilter), n(minf))
	call(g.texParameteri, n(glTexture2D), n(glTexMagFilter), n(glNearest))
	call(g.texParameteri, n(glTexture2D), n(glTexWrapS), n(wrap))
	call(g.texParameteri, n(glTexture2D), n(glTexWrapT), n(wrap))
	return t
}

// newScreen makes the offscreen screen again, d.w by d.h: an sRGB colour
// texture and a depth renderbuffer.
func (d *Device) newScreen() {
	d.forget(render3d.Screen)
	del(d.gl.deleteTextures, d.screen.name)
	del(d.gl.deleteRenderbuffers, d.screen.rb)
	d.screen = texture{w: d.w, h: d.h, srgb: true}
	d.screen.name = d.newTexture(texImage{w: d.w, h: d.h, format: glSRGB8Alpha8, kind: glRGBA, typ: glUnsignedByte})
}
