package soft

import (
	"fmt"
	"image"
	"math"

	"github.com/gabrielluizsf/sge/render3d"
)

// Read is t's pixels: its colours, sRGB encoded when t is sRGB, or, for a
// texture of depth only, its depths as greys, white farthest.
func (d *Device) Read(t render3d.Texture) (*image.NRGBA, error) {
	tx := d.texture(t)
	if tx == nil {
		return nil, fmt.Errorf("soft: no texture %d", t)
	}
	img := image.NewNRGBA(image.Rect(0, 0, tx.w, tx.h))
	for i := range tx.w * tx.h {
		p := img.Pix[4*i : 4*i+4]
		if tx.color == nil {
			g := byte(math.Round(float64(tx.depth[i]) * 255))
			p[0], p[1], p[2], p[3] = g, g, g, 255
			continue
		}
		for k := range 4 {
			v := float64(tx.color[4*i+k])
			if tx.srgb && k < 3 {
				v = srgb(v)
			}
			p[k] = byte(math.Round(math.Max(0, math.Min(1, v)) * 255))
		}
	}
	return img, nil
}
