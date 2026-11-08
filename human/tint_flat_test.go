package human

import (
	"image"
	"image/color"
	"math"
	"testing"
)

// blotchy is a texture of two halves, one dark and one light, as a print
// or washed denim is.
func blotchy(dark, light uint8) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	for y := range 32 {
		for x := range 32 {
			v := dark
			if (x/4+y/4)%2 == 0 {
				v = light
			}
			img.SetNRGBA(x, y, color.NRGBA{v, v, v, 255})
		}
	}
	return img
}

// spread is the standard deviation of the light of img, over its average.
func spread(img *image.NRGBA) float64 {
	var sum, sum2, n float64
	for i := 0; i < len(img.Pix); i += 4 {
		a := float64(img.Pix[i+3]) / 255
		l := luma([3]float64{decode8[img.Pix[i]], decode8[img.Pix[i+1]], decode8[img.Pix[i+2]]})
		sum += a * l
		sum2 += a * l * l
		n += a
	}
	m := sum / n
	return math.Sqrt(sum2/n-m*m) / m
}

// A print of strong light and dark, once one colour, is eased to a mild
// variation, not left as patches; the colour is still the one asked, in
// one hue.
func TestAStrongPrintIsEasedWhenRecoloured(t *testing.T) {
	src := blotchy(20, 230)
	to := RGB(0.3, 0.5, 0.25)
	flat := tinted(src, tint{to: to.linear(), flat: clothContrast})
	if got := spread(flat); got > clothContrast+0.02 {
		t.Errorf("light varies by %.2f of the average, over %.2f", got, clothContrast)
	}
	if spread(flat) == 0 {
		t.Error("the print is gone altogether")
	}
	if loose := spread(tinted(src, tint{to: to.linear()})); loose < 0.6 {
		t.Errorf("without easing the print varies by %.2f: the test's print is too mild", loose)
	}
	for i := 0; i < len(flat.Pix); i += 4 {
		r, g, b := float64(flat.Pix[i]), float64(flat.Pix[i+1]), float64(flat.Pix[i+2])
		if !(g > r && g > b) {
			t.Fatalf("pixel %d is %v,%v,%v: not the green asked", i/4, r, g, b)
		}
	}
}

// A weave that varies less than the limit is left exactly as it was
// recoloured before.
func TestAMildWeaveIsNotEased(t *testing.T) {
	src := blotchy(110, 140)
	to := RGB(0.6, 0.2, 0.2).linear()
	a := tinted(src, tint{to: to})
	b := tinted(src, tint{to: to, flat: clothContrast})
	for i := range a.Pix {
		if a.Pix[i] != b.Pix[i] {
			t.Fatalf("byte %d: %d, eased %d", i, a.Pix[i], b.Pix[i])
		}
	}
}

// textureOf is the colour texture of the proxy item, as the pack has it.
func textureOf(t *testing.T, k *Kit, item string) *image.NRGBA {
	t.Helper()
	data, err := k.lim.read(k.fsys, item)
	if err != nil {
		t.Fatal(err)
	}
	p, err := parseProxy(data, len(k.base.pos))
	if err != nil {
		t.Fatal(item, err)
	}
	matName, err := beside(item, p.material)
	if err != nil {
		t.Fatal(err)
	}
	mdata, err := k.lim.read(k.fsys, matName)
	if err != nil {
		t.Fatal(err)
	}
	file, err := beside(matName, parseMaterial(mdata).diffuse)
	if err != nil {
		t.Fatal(err)
	}
	img, err := k.texture(file)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

// Every cloth the wardrobe has, in two colours of those offered, is recoloured to
// a texture of mild variation, whatever its pack's print.
func TestEveryColourableClothRecoloursCleanly(t *testing.T) {
	k := realKit(t)
	for _, list := range [][]garment{tops, bottoms, dresses} {
		for _, g := range list {
			src := textureOf(t, k, g.item)
			for _, c := range []Colour{ClothColours()[1], ClothColours()[5]} {
				tt, _ := tintTo(c, nil)
				tt.flat = clothContrast
				if s := spread(tinted(src, tt)); s > clothContrast+0.02 {
					t.Errorf("%s in %v: light varies by %.2f of the average", g.item, c, s)
				}
			}
		}
	}
}
