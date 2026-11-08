package human

import (
	"image"
	"testing"
)

// A look without the colour fields makes the very materials and images it
// always did: these fingerprints were taken before they existed.
func TestUnchangedLooksMakeTheSameMaterialsAndImages(t *testing.T) {
	k := realKit(t)
	want := map[string]string{"man": "6c72856e63507715", "old woman": "84f21797c96326bb"}
	for _, c := range baselines {
		m, err := k.Make(c.body, c.look)
		if err != nil {
			t.Fatal(err)
		}
		if got := modelHash(m); got != want[c.name] {
			t.Errorf("%s: fingerprint %s, was %s", c.name, got, want[c.name])
		}
	}
}

// With the colour fields unset or zero, the tiny kit's person is what it was.
func TestUnsetColoursChangeNothing(t *testing.T) {
	k := loaded(t, tinyKit(t))
	body := Body{Gender: 1, Age: 0.5, European: 1, Muscle: 0.5, Weight: 0.5, Height: 0.5}
	base := Look{Skin: "skins/plain/plain.mhmat", Wear: []string{"clothes/shirt/shirt.mhclo"}}
	a, _ := k.Make(body, base)
	withZero := base
	withZero.ClothesColours = []Colour{{}}
	b, _ := k.Make(body, withZero)
	if modelHash(a) != modelHash(b) || len(a.Images) != len(b.Images) {
		t.Error("unset colours changed the person")
	}
	tinted := base
	tinted.ClothesColours = []Colour{RGB(1, 0, 0)}
	c, _ := k.Make(body, tinted)
	shirt := c.Materials[len(c.Materials)-1].BaseColor
	if modelHash(a) == modelHash(c) || shirt[0] <= shirt[2] {
		t.Errorf("a red shirt: %v", shirt)
	}
}

func TestRecolouredHairEyesAndClothesKeepTheirTexture(t *testing.T) {
	k := realKit(t)
	body := Body{Gender: 1, Age: 0.5, European: 1, Muscle: 0.5, Weight: 0.5, Height: 0.5}
	wear := []string{"eyes/high-poly/high-poly.mhclo", "hair/short02/short02.mhclo", "clothes/toigo_basic_tucked_t-shirt/toigo_basic_tucked_t-shirt.mhclo"}
	plain, err := k.Make(body, Look{Eyes: "eyes/materials/brown.mhmat", Wear: wear})
	if err != nil {
		t.Fatal(err)
	}
	blond, shirt, blue := RGB(0.85, 0.65, 0.3), RGB(0.7, 0.15, 0.15), RGB(0.2, 0.5, 0.9)
	dyed, err := k.Make(body, Look{Eyes: "eyes/materials/brown.mhmat", Wear: wear, HairColour: blond, EyeColour: blue, ClothesColours: []Colour{{}, {}, shirt}})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		mat string
		to  Colour
	}{{"hair/short02/short02.mhmat", blond}, {"eyes/materials/brown.mhmat", blue}, {"clothes/toigo_basic_tucked_t-shirt/t_shirt_basic_tucked.mhmat", shirt}} {
		var was, now *image.NRGBA
		for i, mat := range plain.Materials {
			if mat.Name == c.mat {
				was = plain.Images[mat.Texture].(*image.NRGBA)
				now = dyed.Images[dyed.Materials[i].Texture].(*image.NRGBA)
			}
		}
		if was == nil {
			t.Fatalf("no material %s", c.mat)
		}
		wasMean, wasCV := opaqueStats(was)
		mean, cv := opaqueStats(now)
		want := [3]float64{c.to.R, c.to.G, c.to.B}
		if dist(mean, want) >= dist(wasMean, want) {
			t.Errorf("%s: mean %v is no nearer %v than %v", c.mat, mean, want, wasMean)
		}
		if cv < 0.5*wasCV {
			t.Errorf("%s: variation %.3f of %.3f is flat", c.mat, cv, wasCV)
		}
	}
}

func TestSkinToneIsContinuous(t *testing.T) {
	k := realKit(t)
	body := Body{Gender: 0, Age: 0.5, European: 1, Muscle: 0.5, Weight: 0.5, Height: 0.5}
	prev := 2.0
	for _, tone := range []float64{0.02, 0.3, 0.55, 0.8, 1} {
		m, err := k.Make(body, Look{SkinTone: tone})
		if err != nil {
			t.Fatal(err)
		}
		var img *image.NRGBA
		for _, mat := range m.Materials {
			if mat.Name == "skin" {
				img = m.Images[mat.Texture].(*image.NRGBA)
			}
		}
		mean, cv := opaqueStats(img)
		l := 0.2126*mean[0] + 0.7152*mean[1] + 0.0722*mean[2]
		if l >= prev || cv < 0.02 {
			t.Errorf("tone %v: light %.3f after %.3f, variation %.3f", tone, l, prev, cv)
		}
		prev = l
	}
}

// opaqueStats is stats of img's pixels that are at least half opaque.
func opaqueStats(img *image.NRGBA) ([3]float64, float64) {
	out := image.NewNRGBA(image.Rect(0, 0, len(img.Pix)/4, 1))
	n := 0
	for i := 0; i < len(img.Pix); i += 4 {
		if img.Pix[i+3] >= 128 {
			copy(out.Pix[n*4:], img.Pix[i:i+4])
			n++
		}
	}
	return stats(&image.NRGBA{Pix: out.Pix[:n*4], Stride: n * 4, Rect: image.Rect(0, 0, n, 1)})
}
