package soft

import (
	"fmt"
	"math"

	"github.com/gabrielluizsf/sge/render3d"
)

// texture is a texture's pixels: its colours in linear light, four floats
// a pixel, and its depths, either or both.
type texture struct {
	w, h  int
	srgb  bool
	color []float32
	depth []float32
}

// target gives t a colour and a depth buffer to draw into, if it has
// none, cleared to black and to the farthest.
func (t *texture) target() {
	if t.color == nil {
		t.color = make([]float32, 4*t.w*t.h)
	}
	if t.depth == nil {
		t.depth = make([]float32, t.w*t.h)
		for i := range t.depth {
			t.depth[i] = 1
		}
	}
}

// NewTexture is a texture of desc's size: its pixels decoded to linear
// light when sRGB, or cleared for a target.
func (d *Device) NewTexture(desc render3d.TextureDesc) (render3d.Texture, error) {
	w, h := desc.Width, desc.Height
	if w <= 0 || h <= 0 || w > 1<<14 || h > 1<<14 {
		return 0, fmt.Errorf("soft: a texture %d by %d", w, h)
	}
	t := &texture{w: w, h: h, srgb: desc.SRGB}
	switch {
	case desc.Format == render3d.Depth:
		t.depth = make([]float32, w*h)
	case desc.Pixels == nil:
		t.target()
	default:
		if len(desc.Pixels) < 4*w*h {
			return 0, fmt.Errorf("soft: %d bytes for a texture %d by %d", len(desc.Pixels), w, h)
		}
		t.color = make([]float32, 4*w*h)
		for i := range t.color {
			v := float64(desc.Pixels[i]) / 255
			if desc.SRGB && i%4 != 3 {
				v = linear(v)
			}
			t.color[i] = float32(v)
		}
	}
	d.textures = append(d.textures, t)
	return render3d.Texture(len(d.textures) - 1), nil
}

// at is t's colour at uv, its nearest pixel, repeated past its edges.
func (t *texture) at(uv [2]float64) [4]float64 {
	x := int(math.Floor(uv[0] * float64(t.w)))
	y := int(math.Floor(uv[1] * float64(t.h)))
	x, y = ((x%t.w)+t.w)%t.w, ((y%t.h)+t.h)%t.h
	p := t.color[4*(y*t.w+x):]
	return [4]float64{float64(p[0]), float64(p[1]), float64(p[2]), float64(p[3])}
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
