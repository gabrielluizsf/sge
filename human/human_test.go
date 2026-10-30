package human

import (
	"image"
	"math"
	"testing"
	"testing/fstest"
)

// A person stands on the ground, in metres, grown as their body says,
// in their skin, wearing what they wear over what it hides.
func TestMakeShapesDressesAndGroundsAPerson(t *testing.T) {
	k := loaded(t, tinyKit(t))
	look := Look{Skin: "skins/plain/plain.mhmat", Wear: []string{"clothes/shirt/shirt.mhclo"}}
	woman, err := k.Make(Body{Gender: 0, Age: 0.5, European: 1, Muscle: 0.5, Weight: 0.5, Height: 0.5}, look)
	if err != nil {
		t.Fatal(err)
	}
	man, err := k.Make(Body{Gender: 1, Age: 0.5, European: 1, Muscle: 0.5, Weight: 0.5, Height: 0.5}, look)
	if err != nil {
		t.Fatal(err)
	}
	top := func(prims int) float32 {
		var y float32
		for _, p := range man.Meshes[0].Primitives[:prims] {
			for _, v := range p.Positions {
				y = max(y, v[1])
			}
		}
		return y
	}
	if got := top(1); math.Abs(float64(got)-1.0) > 1e-5 {
		t.Errorf("the man's body is %v m tall, not 1: the quad under the shirt is hidden, the target lifts the rest", got)
	}
	body := woman.Meshes[0].Primitives[0]
	if len(body.Indices) != 6 {
		t.Errorf("the body has %d indices, not the 6 of the quad the shirt does not hide", len(body.Indices))
	}
	for _, v := range body.Positions {
		if v[1] < 0 {
			t.Errorf("a vertex below the ground: %v", v)
		}
	}
	if len(woman.Meshes[0].Primitives) != 2 {
		t.Fatalf("%d primitives, not the body and the shirt", len(woman.Meshes[0].Primitives))
	}
	shirt := woman.Meshes[0].Primitives[1]
	// Its first vertex is between the body's 2 and 3, a unit in front.
	if p := shirt.Positions[0]; math.Abs(float64(p[0])) > 1e-6 || math.Abs(float64(p[1])-1) > 1e-6 || math.Abs(float64(p[2])-0.1) > 1e-6 {
		t.Errorf("the shirt's first vertex is at %v, not (0, 1, 0.1)", p)
	}
	if c := woman.Materials[shirt.Material].BaseColor; c != [4]float64{0.2, 0.3, 0.8, 1} {
		t.Errorf("the shirt is %v", c)
	}
	if w := shirt.Weights[0]; math.Abs(float64(w[0]+w[1]+w[2]+w[3])-1) > 1e-6 {
		t.Errorf("the shirt's weights sum to %v", w)
	}
	if len(woman.Skins) != 1 || len(woman.Skins[0].Joints) != 2 {
		t.Errorf("skins %v", woman.Skins)
	}
}

// MakeHuman's weights: at the average, a young woman's and a young man's
// builds weigh half each; a woman of 90 is all old.
func TestBodyWeighsTargetsAsMakeHumanDoes(t *testing.T) {
	v := Average.values()
	if w := weight("universal-female-young-averagemuscle-averageweight", v); math.Abs(w-0.5) > 1e-9 {
		t.Errorf("average young woman %v", w)
	}
	if w := weight("african-male-young", v); math.Abs(w-0.5/3) > 1e-9 {
		t.Errorf("average young African man %v", w)
	}
	old := Body{Gender: 0, Age: 1, Asian: 1, Muscle: 1, Weight: 0, Height: 1, Proportions: 1}.values()
	for name, want := range map[string]float64{
		"asian-female-old": 1, "asian-female-young": 0, "universal-female-old-maxmuscle-minweight": 1,
		"female-old-maxmuscle-minweight-maxheight": 1, "female-old-maxmuscle-minweight-idealproportions": 1,
	} {
		if w := weight(name, old); math.Abs(w-want) > 1e-9 {
			t.Errorf("%s: %v, not %v", name, w, want)
		}
	}
	if y := (Body{Age: 1}).Years(); y != 90 {
		t.Errorf("age 1 is %v years", y)
	}
}

// Skins blend as the body is: a young European woman is her skin alone; a
// man between Africa and Asia, of 50, is half of each, middle-aged.
func TestSkinsAreOfAgeSexAndAncestry(t *testing.T) {
	s := skins(Body{Gender: 0, Age: 0.5, European: 1})
	if len(s) != 1 || s[0].name != "skins/young_caucasian_female/young_caucasian_female.mhmat" || s[0].share != 1 {
		t.Errorf("young European woman: %v", s)
	}
	s = skins(Body{Gender: 1, Age: 0.5 + 25.0/130, African: 1, Asian: 1})
	if len(s) != 2 || math.Abs(s[0].share-0.5) > 1e-9 || s[0].name != "skins/middleage_african_male/middleage_african_male.mhmat" {
		t.Errorf("middle-aged man of Africa and Asia: %v", s)
	}
}

// A file's licence is what its header says, CC0 however written, and a
// "cc0" pack's file that says AGPL3 is AGPL3.
func TestLicenceIsWhatTheFileSays(t *testing.T) {
	fsys := tinyKit(t)
	for name, want := range map[string]string{
		"# This asset was explicitly released as CC0 in september 2020. The license\n": "CC0",
		"# license CC0\nname x\n":      "CC0",
		"# license CC-0\n":             "CC0",
		"# license: CC0\n":             "CC0",
		"# license AGPL3 (see also)\n": "AGPL3 (see also)",
		"# license CC BY 4.0\n":        "CC BY 4.0",
		"name nothing said\nverts 0\n": "",
	} {
		fsys["x.mhclo"] = &fstest.MapFile{Data: []byte(name)}
		got, err := loaded(t, fsys).Licence("x.mhclo")
		if err != nil || got != want {
			t.Errorf("%q: %q %v, not %q", name, got, err, want)
		}
	}
}

// A person's licences are of what they wear, and of the skins blended for
// them when they have none of their own.
func TestLicencesAreOfAllAPersonIsMadeOf(t *testing.T) {
	fsys := tinyKit(t)
	fsys["skins/young_african_male/young_african_male.mhmat"] = &fstest.MapFile{Data: []byte("# license CC0\n")}
	got, err := loaded(t, fsys).Licences(Body{Gender: 1, Age: 0.5, African: 1}, Look{Wear: []string{"clothes/shirt/shirt.mhclo"}})
	want := map[string]string{"clothes/shirt/shirt.mhclo": "", "skins/young_african_male/young_african_male.mhmat": "CC0"}
	if err != nil || len(got) != len(want) {
		t.Fatalf("%v %v", got, err)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: %q, not %q", k, got[k], v)
		}
	}
}

// Hair greys with age, for that person alone: the young keep its colour,
// an old one's turns toward silver; a shirt does not grey.
func TestHairGreysWithAge(t *testing.T) {
	fsys := tinyKit(t)
	fsys["hair/dark/dark.mhclo"] = &fstest.MapFile{Data: []byte("obj_file dark.obj\nmaterial dark.mhmat\nverts 0\n0\n1\n2\n")}
	fsys["hair/dark/dark.obj"] = &fstest.MapFile{Data: []byte("v 0 0 0\nv 0 0 0\nv 0 0 0\nf 1 2 3\n")}
	fsys["hair/dark/dark.mhmat"] = &fstest.MapFile{Data: []byte("diffuseTexture dark.png\n")}
	fsys["hair/dark/dark.png"] = &fstest.MapFile{Data: pngOf(t, 2, colorOf(20))}
	k := loaded(t, fsys)
	hairOf := func(years float64) uint8 {
		m, err := k.Make(Body{Age: (years-25)/130 + 0.5, European: 1}, Look{Skin: "skins/plain/plain.mhmat", Wear: []string{"hair/dark/dark.mhclo"}})
		if err != nil {
			t.Fatal(err)
		}
		hair := m.Meshes[0].Primitives[1]
		return m.Images[m.Materials[hair.Material].Texture].(*image.NRGBA).Pix[0]
	}
	young, old := hairOf(30), hairOf(80)
	if young != 20 || old < 120 {
		t.Errorf("hair at 30 is %d, at 80 %d", young, old)
	}
	if again := hairOf(30); again != 20 {
		t.Errorf("an old one's grey hair greyed the young's: %d", again)
	}
	if Greying(45) != 0 || Greying(100) != 0.85 {
		t.Error(Greying(45), Greying(100))
	}
}
