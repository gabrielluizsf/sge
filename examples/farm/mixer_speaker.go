package main

import (
	"math"

	"github.com/gabrielluizsf/sge/examples/farm/sfx"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sound"
	"github.com/gabrielluizsf/sge/vector"
)

// mixerSpeaker plays the farm's sounds on the game's mixer. A sound in the
// world is heard from where it is on the screen, left or right, and each
// play takes the next of its versions, a little higher or lower and
// louder or softer than the last, as no two steps sound the same.
type mixerSpeaker struct {
	g     *game.Game
	clips map[string][]*sound.Clip
	turn  map[string]int
	luck  uint32
}

// gain is how loud each sound is played, the HUD's quietest: they are
// heard most often.
var gain = map[string]float64{
	"step_grass": 0.35, "step_soil": 0.4, "step_path": 0.35, "stroke": 0.45, "splash": 0.6, "chirp": 0.22,
	"ui_tick": 0.25, "ui_open": 0.3, "ui_close": 0.3, "coin": 0.3, "sell": 0.35, "deny": 0.3,
	"plant": 0.4, "pour": 0.35, "harvest": 0.35, "chop": 0.4, "fell": 0.5, "clang": 0.25,
}

func newMixerSpeaker(g *game.Game) *mixerSpeaker {
	return &mixerSpeaker{g: g, clips: sfx.Clips(), turn: map[string]int{}, luck: 1}
}

func (s *mixerSpeaker) play(name string, at *vector.Vector3) {
	versions := s.clips[name]
	if len(versions) == 0 || s.g.Sound == nil {
		return
	}
	i := s.turn[name] % len(versions)
	s.turn[name]++
	v := sound.Voice{Gain: gain[name] * (0.9 + 0.2*s.chance()), Pitch: 0.94 + 0.12*s.chance()}
	if at != nil {
		v.Pan = s.pan(*at)
		if v.Gain *= s.near(*at); v.Gain < 0.01 {
			return
		}
	}
	s.g.Sound.Play(versions[i], v)
}

// pan is where across the screen a point of the world is heard from,
// softened so nothing is only in one ear.
func (s *mixerSpeaker) pan(at vector.Vector3) float64 {
	w, _ := s.g.View()
	if w == 0 {
		return 0
	}
	p := (s.g.ToScreen(at).X/float64(w))*2 - 1
	return max(-0.7, min(0.7, p*0.7))
}

// chance is the next of the speaker's random numbers, from 0 to 1.
func (s *mixerSpeaker) chance() float64 {
	s.luck = s.luck*1664525 + 1013904223
	return float64(s.luck>>8) / float64(1<<24)
}

// earshot is how far from the middle of the screen, in pixels, a sound
// in the world fades away to nothing.
const earshot = 900

// near is how loud a sound at a point of the world is heard: in full on
// the screen's middle, fading to nothing earshot pixels from it, so a
// villager's axe across town is faint.
func (s *mixerSpeaker) near(at vector.Vector3) float64 {
	w, h := s.g.View()
	if w == 0 {
		return 1
	}
	p := s.g.ToScreen(at)
	dx, dy := p.X-float64(w)/2, p.Y-float64(h)/2
	return max(0, 1-math.Sqrt(dx*dx+dy*dy)/earshot)
}
