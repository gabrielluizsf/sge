package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm/sfx"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sound"
)

// loudArt is an art that brings one sound of its own, and one of no versions.
type loudArt struct {
	Art
	clip *sound.Clip
}

func (l loudArt) Sounds() map[string][]*sound.Clip {
	return map[string][]*sound.Clip{"ui_tick": {l.clip}, "deny": nil}
}

// An art's sound plays instead of the one made up, and the others stay.
func TestAnArtsSoundsReplaceOnlyTheirOwn(t *testing.T) {
	clip := &sound.Clip{Rate: sound.Rate, Samples: []float32{0.5}}
	s := newMixerSpeaker(game.New())
	before := s.clips["chop"]
	s.useSounds(loudArt{Art: SVGArt(), clip: clip})
	if got := s.clips["ui_tick"]; len(got) != 1 || got[0] != clip {
		t.Errorf("ui_tick is %v", got)
	}
	if got := s.clips["deny"]; len(got) == 0 || got[0] == nil {
		t.Error("a sound of no versions took the place of deny")
	}
	if len(s.clips["chop"]) != len(before) || s.clips["chop"][0] != before[0] {
		t.Error("chop changed")
	}
}

// An art with no sounds of its own keeps them all, a level up among them,
// and every sound the speaker has a volume for.
func TestAPlainArtKeepsItsSounds(t *testing.T) {
	s := newMixerSpeaker(game.New())
	s.useSounds(SVGArt())
	if len(s.clips) != len(sfx.Clips()) {
		t.Errorf("%d sounds of %d", len(s.clips), len(sfx.Clips()))
	}
	if len(s.clips["level_up"]) == 0 || gain["level_up"] == 0 {
		t.Error("no level up to hear")
	}
}
