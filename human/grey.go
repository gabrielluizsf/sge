package human

import (
	"image"
	"strings"
)

// Greying is how grey hair is at years: none to 45, most of it by 80.
func Greying(years float64) float64 { return min(max((years-45)/35, 0), 0.85) }

// greys is whether the proxy in file name greys with age: hair and brows.
func greys(name string) bool {
	return strings.HasPrefix(name, "hair/") || strings.HasPrefix(name, "eyebrows/")
}

// greyOut turns material mat's colour texture grey by d's greying, as a
// texture of its own: others wearing the same hair keep theirs.
func (d *dresser) greyOut(mat int) {
	m := &d.model.Materials[mat]
	if d.grey <= 0 || m.Texture < 0 {
		return
	}
	src, ok := d.model.Images[m.Texture].(*image.NRGBA)
	if !ok {
		return
	}
	d.model.Images = append(d.model.Images, greyed(src, d.grey))
	m.Texture = len(d.model.Images) - 1
}

// greyed is img grey by amount, from 0 to 1: each pixel toward a silver
// as light as it is, its alpha kept.
func greyed(img *image.NRGBA, amount float64) *image.NRGBA {
	out := image.NewNRGBA(img.Rect)
	copy(out.Pix, img.Pix)
	for i := 0; i+3 < len(out.Pix); i += 4 {
		p := out.Pix[i : i+3]
		lum := 0.3*float64(p[0]) + 0.59*float64(p[1]) + 0.11*float64(p[2])
		silver := 150 + 0.4*lum
		for c := range p {
			p[c] = uint8(float64(p[c])*(1-amount) + silver*amount + 0.5)
		}
	}
	return out
}
