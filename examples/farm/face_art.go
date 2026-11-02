package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/face"
)

// FaceArt is what an Art may also have whose portraits live: who's head
// and shoulders with their face as look has it, blinking and speaking,
// or nil when it has no such picture, and the still Portrait is shown.
type FaceArt interface {
	Face(who Who, look face.Look) *canvas.Canvas
}

// talkPortrait is the conversation's portrait of v: their face as it is
// now, if the art draws faces, else their still one.
func (f *farm) talkPortrait(v *villager) *canvas.Canvas {
	if fa, ok := f.art.(FaceArt); ok && f.talk.blink != nil {
		look := face.Look{Lids: f.talk.blink.Lids(), Mouth: f.talk.mouth.Shape()}
		if p := fa.Face(v.who, look); p != nil {
			return p
		}
	}
	return f.art.Portrait(v.who)
}

// faceUpdate lets dt seconds go by on the face of who is talking: their
// eyes blink, and their mouth says the letter being shown.
func (f *farm) faceUpdate(dt float64) {
	t := f.talk
	if t.blink == nil {
		t.blink = face.NewBlinker(seedOf(t.v.name))
	}
	t.blink.Update(dt)
	r, saying := f.talkBox().Saying()
	t.mouth.Update(dt, r, saying)
}

// seedOf is a seed of its own for name, so each villager blinks in their
// own time.
func seedOf(name string) uint32 {
	h := uint32(2166136261)
	for _, r := range name {
		h = (h ^ uint32(r)) * 16777619
	}
	return h
}
