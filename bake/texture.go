package bake

import (
	"image"
	"math"
)

// texture is an image's pixels in linear light, red, green, blue and
// alpha, for sampling.
type texture struct {
	w, h int
	px   [][4]float64
}

// textureOf reads img: into linear light, its sRGB decoded, for colour;
// as its values are, for a map of data, normals or roughness.
func textureOf(img image.Image, colour bool) *texture {
	b := img.Bounds()
	t := &texture{w: b.Dx(), h: b.Dy(), px: make([][4]float64, b.Dx()*b.Dy())}
	for y := range t.h {
		for x := range t.w {
			r, g, bl, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			if a == 0 {
				continue
			}
			// RGBA is premultiplied: undo it before decoding the curve.
			af := float64(a)
			c := [3]float64{float64(r) / af, float64(g) / af, float64(bl) / af}
			if colour {
				c = [3]float64{linear(c[0]), linear(c[1]), linear(c[2])}
			}
			t.px[y*t.w+x] = [4]float64{c[0], c[1], c[2], af / 0xffff}
		}
	}
	return t
}

// linear is an sRGB value, 0 to 1, in linear light; srgb is the reverse.
func linear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func srgb(c float64) float64 {
	c = math.Max(0, math.Min(1, c))
	if c <= 0.0031308 {
		return c * 12.92
	}
	return 1.055*math.Pow(c, 1/2.4) - 0.055
}

// at is the texture at uv, repeated past its edges, blended between its
// four nearest pixels.
func (t *texture) at(uv [2]float64) [4]float64 {
	x := uv[0]*float64(t.w) - 0.5
	y := uv[1]*float64(t.h) - 0.5
	x0, y0 := math.Floor(x), math.Floor(y)
	fx, fy := x-x0, y-y0
	px := func(x, y float64) [4]float64 {
		xi := ((int(x) % t.w) + t.w) % t.w
		yi := ((int(y) % t.h) + t.h) % t.h
		return t.px[yi*t.w+xi]
	}
	a, b, c, d := px(x0, y0), px(x0+1, y0), px(x0, y0+1), px(x0+1, y0+1)
	var out [4]float64
	for i := range out {
		out[i] = (a[i]*(1-fx)+b[i]*fx)*(1-fy) + (c[i]*(1-fx)+d[i]*fx)*fy
	}
	return out
}
