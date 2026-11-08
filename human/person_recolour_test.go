package human

import (
	"bytes"
	"image"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// pixelsOf is the colour texture of material i of m, nil for none.
func pixelsOf(m *gltf.Model, i int) []byte {
	if t := m.Materials[i].Texture; t >= 0 {
		return m.Images[t].(*image.NRGBA).Pix
	}
	return nil
}

// Recolouring a person gives every material the colours and textures
// that making them in the new look gives, the mesh untouched.
func TestRecolouringAPersonIsMakingThemInTheNewLook(t *testing.T) {
	k := realKit(t)
	for _, c := range baselines {
		a := c.look
		b := c.look
		b.SkinTone, b.HairColour, b.EyeColour = 0.8, RGB(0.85, 0.65, 0.3), RGB(0.2, 0.5, 0.9)
		b.Eyes = "eyes/materials/grey.mhmat"
		if a.Eyes == b.Eyes {
			b.Eyes = "eyes/materials/brownlight.mhmat"
		}
		b.ClothesColours = make([]Colour, len(b.Wear))
		for i := range b.ClothesColours {
			b.ClothesColours[i] = RGB(0.7, 0.15, 0.15)
		}
		for _, steps := range [][2]Look{{a, b}, {b, a}, {b, b}} {
			p, err := k.MakePerson(c.body, steps[0])
			if err != nil {
				t.Fatal(err)
			}
			verts := len(p.Model.Meshes[0].Primitives)
			p.Recolour(steps[1])
			want, err := k.Make(c.body, steps[1])
			if err != nil {
				t.Fatal(err)
			}
			if len(p.Model.Materials) != len(want.Materials) || len(p.Model.Meshes[0].Primitives) != verts {
				t.Fatalf("%s: the person's shape changed", c.name)
			}
			for i := range want.Materials {
				if got, w := p.Model.Materials[i].BaseColor, want.Materials[i].BaseColor; got != w {
					t.Errorf("%s: material %d colour %v, want %v", c.name, i, got, w)
				}
				if !bytes.Equal(pixelsOf(p.Model, i), pixelsOf(want, i)) {
					t.Errorf("%s: material %d (%s) has other pixels", c.name, i, want.Materials[i].Name)
				}
			}
		}
	}
}

// A recolour that changes nothing reports no material changed.
func TestRecolouringToTheSameLookChangesNothing(t *testing.T) {
	k := realKit(t)
	c := baselines[0]
	p, err := k.MakePerson(c.body, c.look)
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Recolour(c.look); len(got) != 0 {
		t.Errorf("materials %v changed", got)
	}
	l := c.look
	l.HairColour = RGB(0.1, 0.1, 0.1)
	if got := p.Recolour(l); len(got) == 0 {
		t.Error("a new hair colour changed no material")
	}
}

// The skin blended for a body is kept for the next person of the same
// shares: the same picture, not a second blend.
func TestTheLastSkinBlendIsKept(t *testing.T) {
	k := realKit(t)
	c := baselines[0]
	shares := skins(c.body)
	a, err := k.blendSkin(shares)
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := k.blendSkin(shares); a != b {
		t.Error("blended twice")
	}
	other := c.body
	other.Age = 0.9
	if b, _ := k.blendSkin(skins(other)); a == b {
		t.Error("another body got the same skin")
	}
}

// Plans made before the earlier ones are applied still end in the look of
// the last, when applied in order, even if a colour went back to what it was.
func TestPlansAppliedInOrderEndInTheLastLook(t *testing.T) {
	k := realKit(t)
	c := baselines[0]
	p, err := k.MakePerson(c.body, c.look)
	if err != nil {
		t.Fatal(err)
	}
	red, blue := c.look, c.look
	red.HairColour, red.SkinTone = RGB(0.8, 0.1, 0.1), 0.7
	blue.HairColour, blue.SkinTone = RGB(0.1, 0.1, 0.8), 0.3
	plans := []*Recolouring{p.Plan(red), p.Plan(blue), p.Plan(c.look)} // and back to the first look
	for _, r := range plans {
		p.Apply(r)
	}
	want, _ := k.Make(c.body, c.look)
	for i := range want.Materials {
		if !bytes.Equal(pixelsOf(p.Model, i), pixelsOf(want, i)) {
			t.Errorf("material %d (%s) is not as first made", i, want.Materials[i].Name)
		}
	}
}
