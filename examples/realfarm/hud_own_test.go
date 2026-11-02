package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/flora"
	"github.com/gabrielluizsf/sge/life"
)

// same reports whether a and b are the same picture: the same one, or of
// one size and every pixel.
func same(a, b *canvas.Canvas) bool {
	if a == b {
		return true
	}
	if a.Width != b.Width || a.Height != b.Height {
		return false
	}
	for y := range a.Height {
		for x := range a.Width {
			if a.At(x, y) != b.At(x, y) {
				return false
			}
		}
	}
	return true
}

// Nothing of the HUD is the stylized art's: every icon, the cursor, the
// panels and the text are the realistic art's own, whether the models
// have been baked or not.
func TestNoHUDPictureIsTheStylizedArts(t *testing.T) {
	real, style := testArt(t), farm.SVGArt()
	names := []string{"icon_coin", "icon_bubble", "icon_can", "icon_hoe", "icon_heart", "icon_heart_empty", "icon_wood",
		"icon_axe_stone", "icon_axe_iron", "icon_axe_gold", "icon_axe_diamond"}
	for _, c := range flora.Crops {
		if c.Name == "withered" {
			continue // a plant's state, not an item
		}
		names = append(names, "icon_"+c.Name, "icon_"+c.Name+"Seeds")
	}
	for _, name := range names {
		if got := real.Icon(name); got == nil || got.Width == 0 || same(got, style.Icon(name)) {
			t.Errorf("%s is not the realistic art's own", name)
		}
	}
	if same(real.Cursor(), style.Cursor()) {
		t.Error("the cursor is the stylized one")
	}
	for _, sel := range []bool{false, true} {
		if same(real.Panel(20, 20, sel), style.Panel(20, 20, sel)) {
			t.Errorf("panel (selected %v) is the stylized one", sel)
		}
	}
	if same(real.Text("Dia 1", 'k'), style.Text("Dia 1", 'k')) {
		t.Error("the text is the stylized face")
	}
	if p := real.Portrait(farm.Who{Stage: life.Adult}); p != nil && same(p, style.Portrait(farm.Who{Stage: life.Adult})) {
		t.Error("the portrait is the stylized one")
	}
}

// The credits name the font, MakeHuman, each Poly Haven asset with its
// authors, each with the text of its licence.
func TestTheCreditsNameEveryComponent(t *testing.T) {
	t.Setenv("SGE_IMPORTS", t.TempDir()) // none of the user's own models
	list := testArt(t).(farm.Credited).Credits()
	var serif, human, haven int
	for _, c := range list {
		if c.Name == "" || c.Author == "" || !strings.HasPrefix(c.URL, "https://") || c.License == nil || len(c.License()) < 200 {
			t.Errorf("credit %+v", c)
		}
		switch {
		case strings.Contains(c.Name, "Source Serif"):
			serif++
			if !strings.Contains(c.License(), "SIL OPEN FONT LICENSE") {
				t.Error("the font's licence is not the OFL")
			}
		case strings.Contains(c.Name, "MakeHuman"):
			human++
			if !strings.Contains(c.License(), "CC0 1.0 Universal") {
				t.Error("MakeHuman's licence is not CC0")
			}
		case strings.HasPrefix(c.Name, "Poly Haven"):
			haven++
		default:
			t.Errorf("%s is not one of the components", c.Name)
		}
	}
	if serif != 1 || human != 1 || haven < 40 {
		t.Errorf("%d fonts, %d MakeHuman, %d Poly Haven", serif, human, haven)
	}
}

// The user's Sketchfab models are credited from the licence files in their
// own folder, read when asked for, and not at all where there is none.
func TestSketchfabCreditsComeFromTheUsersFolder(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SGE_IMPORTS", dir)
	if got := sketchfab(); len(got) != 0 {
		t.Fatalf("credits with no folder: %v", got)
	}
	if err := os.MkdirAll(filepath.Join(dir, "water"), 0o755); err != nil {
		t.Fatal(err)
	}
	text := "Model Information:\n* title:\tRealistic Water\n* source:\thttps://sketchfab.com/3d-models/x\n* author:\tNegitropic (https://sketchfab.com/negitropic)\n\nModel License:\n* license type:\tSKETCHFAB Standard\n"
	if err := os.WriteFile(filepath.Join(dir, "water", licenceFile), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	got := sketchfab()
	if len(got) != 1 || got[0].Name != "Sketchfab - Realistic Water" || !strings.HasPrefix(got[0].Author, "Negitropic") || got[0].URL != "https://sketchfab.com/3d-models/x" {
		t.Fatalf("credits %+v", got)
	}
	if got[0].License() != text {
		t.Error("the licence shown is not the file's text")
	}
}
