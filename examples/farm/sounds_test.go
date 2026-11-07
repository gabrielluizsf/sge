package main

import (
	"strings"
	"testing"
)

// Samples are read in order, signed, from -1 up to nearly 1, and a stray
// last byte is no sample.
func TestClipOfReadsSamples(t *testing.T) {
	c := clipOf([]byte{0x00, 0x00, 0xFF, 0x7F, 0x00, 0x80, 0x01})
	if len(c.Samples) != 3 || c.Samples[0] != 0 || c.Samples[1] < 0.99 || c.Samples[2] != -1 {
		t.Errorf("%v", c.Samples)
	}
}

// The look brings soft, short sounds for everything it has one of: none so
// long or loud that a few together pile up, and a menu's tick shortest of all.
func TestTheLooksSoundsAreShortAndSoft(t *testing.T) {
	a := testArt(t).(*realArt)
	got := a.Sounds()
	for _, name := range []string{"ui_tick", "ui_open", "ui_close", "deny", "coin", "sell", "level_up",
		"step_grass", "step_soil", "step_path", "chop", "fell", "clang", "plant", "harvest"} {
		if len(got[name]) == 0 {
			t.Errorf("no %s", name)
		}
	}
	for name, versions := range got {
		for i, c := range versions {
			var peak float32
			for _, s := range c.Samples {
				peak = max(peak, s, -s)
			}
			seconds := float64(len(c.Samples)) / float64(c.Rate)
			if seconds <= 0 || seconds > 1 {
				t.Errorf("%s %d lasts %.2f s", name, i, seconds)
			}
			if peak < 0.2 || peak > 0.55 {
				t.Errorf("%s %d peaks at %.2f", name, i, peak)
			}
		}
	}
	tick := got["ui_tick"][0]
	if seconds := float64(len(tick.Samples)) / float64(tick.Rate); seconds > 0.3 {
		t.Errorf("a menu's tick lasts %.2f s", seconds)
	}
	if len(got["step_grass"]) < 3 {
		t.Error("a step has too few versions: it would repeat")
	}
}

// The name of a sound comes from its file, whatever underscores it has.
func TestSoundName(t *testing.T) {
	if got := soundName("step_grass_2.pcm"); got != "step_grass" {
		t.Errorf("%q", got)
	}
}

// The credits say whose the sounds are.
func TestTheSoundsAreCredited(t *testing.T) {
	for _, c := range testArt(t).(*realArt).Credits() {
		if strings.Contains(c.Name, "Kenney") {
			return
		}
	}
	t.Error("no credit for the sounds")
}
