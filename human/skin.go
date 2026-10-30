package human

import (
	"image"
	"image/color"
)

// skinShare is one of MakeHuman's skins and how much of it a body has.
type skinShare struct {
	name  string
	share float64
}

// skins are the shares of MakeHuman's eighteen skins, young, middle-aged
// and old, of each sex and ancestry, a body has: as its body is shaped,
// a mix, those of under 2% left out.
func skins(b Body) []skinShare {
	v := b.values()
	ages := ageShares(b.Years())
	var out []skinShare
	var sum float64
	for _, race := range []string{"african", "asian", "caucasian"} {
		for _, sex := range []string{"female", "male"} {
			for i, age := range []string{"young", "middleage", "old"} {
				s := v[race] * v[sex] * ages[i]
				if s < 0.02 {
					continue
				}
				name := age + "_" + race + "_" + sex
				out = append(out, skinShare{"skins/" + name + "/" + name + ".mhmat", s})
				sum += s
			}
		}
	}
	for i := range out {
		out[i].share /= sum
	}
	return out
}

// ageShares are how young, middle-aged and old skin of years looks: young
// to 25, middle-aged at 50, old from 80.
func ageShares(years float64) [3]float64 {
	switch {
	case years <= 25:
		return [3]float64{1, 0, 0}
	case years <= 50:
		t := (years - 25) / 25
		return [3]float64{1 - t, t, 0}
	case years <= 80:
		t := (years - 50) / 30
		return [3]float64{0, 1 - t, t}
	}
	return [3]float64{0, 0, 1}
}

// blendSkin is the skins' textures mixed by their shares, each sampled
// at the first's size.
func (k *Kit) blendSkin(shares []skinShare) (*image.NRGBA, error) {
	var out *image.NRGBA
	var acc [][4]float64
	for _, s := range shares {
		data, err := k.lim.read(k.fsys, s.name)
		if err != nil {
			return nil, err
		}
		m := parseMaterial(data)
		name, err := beside(s.name, m.diffuse)
		if err != nil {
			return nil, err
		}
		img, err := k.texture(name)
		if err != nil {
			return nil, err
		}
		if out == nil {
			out = image.NewNRGBA(img.Rect)
			acc = make([][4]float64, len(img.Pix)/4)
		}
		w, h := out.Rect.Dx(), out.Rect.Dy()
		iw, ih := img.Rect.Dx(), img.Rect.Dy()
		for y := range h {
			for x := range w {
				c := img.NRGBAAt(x*iw/w, y*ih/h)
				a := &acc[y*w+x]
				a[0] += s.share * float64(c.R)
				a[1] += s.share * float64(c.G)
				a[2] += s.share * float64(c.B)
				a[3] += s.share * float64(c.A)
			}
		}
	}
	for i, a := range acc {
		out.SetNRGBA(i%out.Rect.Dx(), i/out.Rect.Dx(), color.NRGBA{uint8(a[0] + 0.5), uint8(a[1] + 0.5), uint8(a[2] + 0.5), uint8(a[3] + 0.5)})
	}
	return out, nil
}
