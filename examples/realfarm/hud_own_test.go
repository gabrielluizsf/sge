package main

import (
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
	names := []string{"icon_coin", "icon_bubble", "icon_can", "icon_hoe", "icon_heart", "icon_heart_empty", "icon_food", "icon_wood",
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
	list := testArt(t).(farm.Credited).Credits()
	var serif, human, haven, sounds, animals int
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
		case strings.Contains(c.Name, "Kenney"):
			sounds++
			if !strings.Contains(c.License(), "CC0 1.0 Universal") {
				t.Error("the sounds' licence is not CC0")
			}
		case strings.Contains(c.Name, "changed:"):
			animals++
			if !strings.Contains(c.License(), "Attribution 4.0") {
				t.Errorf("%s's licence is not CC-BY 4.0", c.Name)
			}
		default:
			t.Errorf("%s is not one of the components", c.Name)
		}
	}
	if serif != 1 || human != 1 || haven < 40 || sounds != 1 || animals != 4 {
		t.Errorf("%d fonts, %d MakeHuman, %d Poly Haven, %d sound packs, %d animals", serif, human, haven, sounds, animals)
	}
}

// No bird of the stylized farm flies over the realistic one: until a
// shippable model of one is pinned, its art draws no birds at all.
func TestNoBirdIsTheStylizedArts(t *testing.T) {
	if _, ok := testArt(t).(farm.BirdArt); ok {
		t.Error("the realistic art draws birds with no model of its own")
	}
}
