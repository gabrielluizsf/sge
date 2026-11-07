package human

import (
	"image"
	"math"
)

// tint is how a texture is recoloured: the colour, in linear light, its
// average becomes; how much of its own hues stays, 0 to 1; and, if not
// nil, where it applies, 0 to 1, at a place across the image given as
// fractions of its width.
type tint struct {
	to   [3]float64
	keep float64
	mask func(x, y float64) float64
	// flat, if over 0, is the most the light of the texture may vary
	// about its average, as a share of it: a print more contrasty than
	// that, once one colour, would be patches of light and dark, so it is
	// eased to flat, and the weave, which varies less, is left as it is.
	flat float64
}

// tinted is img recoloured: in linear light, each pixel is as much lighter
// or darker than the image's average as it was, in the colour asked, so
// strands, weave and grain stay, and the average is that colour. Alpha
// is kept, and img is not changed.
func tinted(img *image.NRGBA, t tint) *image.NRGBA {
	out := image.NewNRGBA(img.Rect)
	copy(out.Pix, img.Pix)
	w, h := img.Rect.Dx(), img.Rect.Dy()
	at := func(x, y int) float64 {
		if t.mask == nil {
			return 1
		}
		return t.mask(float64(x)/float64(w), float64(y)/float64(h))
	}
	var sum, sum2, weight float64
	for y := range h {
		for x := range w {
			i := y*img.Stride + x*4
			a := float64(img.Pix[i+3]) / 255 * at(x, y)
			l := luma([3]float64{decode8[img.Pix[i]], decode8[img.Pix[i+1]], decode8[img.Pix[i+2]]})
			sum += a * l
			sum2 += a * l * l
			weight += a
		}
	}
	if weight == 0 || sum == 0 {
		return out
	}
	mean, goal := sum/weight, luma(t.to)
	ease := 1.0
	if sd := math.Sqrt(max(sum2/weight-mean*mean, 0)); t.flat > 0 && sd > t.flat*mean {
		ease = t.flat * mean / sd
	}
	for y := range h {
		for x := range w {
			m := at(x, y)
			if m == 0 {
				continue
			}
			i := y*img.Stride + x*4
			var l, c [3]float64
			for k := range 3 {
				l[k] = decode8[img.Pix[i+k]]
			}
			ratio, peak := luma(l)/mean, 0.0
			ratio = 1 + (ratio-1)*ease
			for k := range 3 {
				c[k] = (1-t.keep)*t.to[k]*ratio + t.keep*l[k]*goal/mean
				peak = max(peak, c[k])
			}
			for k := range 3 {
				if peak > 1 {
					c[k] /= peak
				}
				v := l[k] + (c[k]-l[k])*m
				out.Pix[i+k] = encode8[int(min(max(v, 0), 1)*encodeSteps+0.5)]
			}
		}
	}
	return out
}
