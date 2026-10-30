package farm

import (
	"strings"

	"github.com/gabrielluizsf/sge/game"
)

// heard is what of the farmer's animation was last heard: the steps and
// strokes are played as the frames they land on are first shown.
type heard struct {
	anim  string
	frame int
	swam  bool
}

// footsteps plays the farmer's sounds as they move: a step on what is
// under them as a foot lands, a stroke as an arm sweeps the water, and a
// splash going into it.
func (p *farmer) footsteps(o *game.Object) {
	anim, frame := o.Animator.Current(), o.Animator.FrameIndex()
	at := feet(o)
	if p.swims && !p.heard.swam {
		p.f.sound("splash", &at)
	}
	if anim != p.heard.anim || frame != p.heard.frame {
		switch {
		case strings.HasSuffix(anim, "_swim") && frame == 0:
			p.f.sound("stroke", &at)
		case strings.HasSuffix(anim, "_walk") && frame%2 == 0:
			p.f.sound(p.f.stepOn(at.X, at.Y), &at)
		}
	}
	p.heard = heard{anim: anim, frame: frame, swam: p.swims}
}

// stepOn is the sound of a step at the world point x, y: on the path, on
// tilled soil, or in grass, as the wild is all over; a room's floor
// sounds as the path.
func (f *farm) stepOn(x, y float64) string {
	if f.room != nil {
		return "step_path"
	}
	col, row := f.wild.CellOf(x, y)
	switch f.groundAt(col, row) {
	case tileID["path"]:
		return "step_path"
	case tileID["soil_dry"], tileID["soil_wet"]:
		return "step_soil"
	}
	return "step_grass"
}
