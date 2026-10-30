package farm

import (
	"github.com/gabrielluizsf/sge/examples/farm/sfx"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sound"
)

// addSound gives the game its sounds: a mixer Run plays, the farm a
// speaker on it, and the ambience of wherever the farmer goes.
func addSound(g *game.Game, f *farm) {
	g.Sound = sound.NewMixer(sound.Rate)
	f.speaker = newMixerSpeaker(g)
	loops := sfx.Loops()
	a := &ambience{f: f, farmer: g.Find("farmer"), luck: 7, bird: 3}
	a.leaves, a.water = g.Sound.Loop(loops["rustle"]), g.Sound.Loop(loops["lap"])
	o := game.NewObject("ambience")
	o.AddScript(a)
	g.Add(o)
}
