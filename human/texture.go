package human

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/gabrielluizsf/sge/gltf"
)

// texture is the image in file name, decoded no larger than the kit's
// limit allows, and shrunk to its Texture pixels across: read once, for
// every person who wears it.
func (k *Kit) texture(name string) (*image.NRGBA, error) {
	if img, ok := k.pictures[name]; ok {
		return img, nil
	}
	img, err := k.decode(name)
	if err != nil {
		return nil, err
	}
	k.pictures[name] = img
	return img, nil
}

// decode is the image in file name, as texture keeps it.
func (k *Kit) decode(name string) (*image.NRGBA, error) {
	data, err := k.lim.read(k.fsys, name)
	if err != nil {
		return nil, err
	}
	img, err := gltf.DecodeImage(data, k.lim.Pixels)
	if err != nil {
		return nil, err
	}
	return shrink(img, k.lim.Texture), nil
}

// shrink is img no more than most pixels across, each pixel the mean of
// the square of img's it covers.
func shrink(img image.Image, most int) *image.NRGBA {
	b := img.Bounds()
	f := 1
	for most > 0 && max(b.Dx(), b.Dy())/f > most {
		f++
	}
	src := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(src, src.Rect, img, b.Min, draw.Src)
	if f == 1 {
		return src
	}
	w, h := b.Dx()/f, b.Dy()/f
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			var s [4]int
			for dy := range f {
				for dx := range f {
					c := src.NRGBAAt(x*f+dx, y*f+dy)
					s[0], s[1], s[2], s[3] = s[0]+int(c.R), s[1]+int(c.G), s[2]+int(c.B), s[3]+int(c.A)
				}
			}
			n := f * f
			out.SetNRGBA(x, y, color.NRGBA{uint8(s[0] / n), uint8(s[1] / n), uint8(s[2] / n), uint8(s[3] / n)})
		}
	}
	return out
}
