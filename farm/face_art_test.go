package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/dialogue"
	"github.com/gabrielluizsf/sge/face"
)

// facing is the stylized art with faces that keep the looks asked for.
type facing struct {
	Art
	looks []face.Look
}

func (a *facing) Face(_ Who, l face.Look) *canvas.Canvas {
	a.looks = append(a.looks, l)
	cv, _ := canvas.NewCanvas(4, 4)
	return cv
}

// The face of who is talking lives: their mouth moves while the words
// appear and shuts when all is said, and their eyes blink every few
// seconds, words or none.
func TestTheTalkingFaceLives(t *testing.T) {
	art := &facing{Art: SVGArt()}
	_, f := build(art, idle{})
	v := tradesman(f, eastTown(f), roleLumberjack)
	f.openTalk(v)
	screen, _ := canvas.NewCanvas(1280, 720)
	f.talkBox().Paint(screen, f.art, nil) // lays the words out
	const dt = 1.0 / 60
	spoke := false
	for !f.talkBox().Done() {
		f.talkUpdate(dt, dialogue.Input{})
		f.talkPortrait(v)
		spoke = spoke || art.looks[len(art.looks)-1].Mouth != face.Shut
	}
	if !spoke {
		t.Error("the mouth never moved while speaking")
	}
	art.looks = nil
	for range int(7 / dt) {
		f.talkUpdate(dt, dialogue.Input{})
		f.talkPortrait(v)
	}
	blinked := false
	for _, l := range art.looks[2:] {
		blinked = blinked || l.Lids == face.LidsShut
		if l.Mouth != face.Shut {
			t.Fatal("the mouth moves with everything said")
		}
	}
	if !blinked {
		t.Error("no blink in seven seconds")
	}
}
