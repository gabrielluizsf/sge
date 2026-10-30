package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/sprite"
)

// svgBody is a person drawn from SVG, facing four ways: frame is the SVG
// of each of their frames by name, the way they face and then what they
// do, as "left", "left1" or "leftChop1"; size is how tall they are drawn,
// as a part of a grown person's height, 0 for whole.
type svgBody struct {
	frame func(name string) string
	size  float64
}

func (svgBody) Facings() int { return 4 }

// Clip is anim facing dir: standing still, a step with each foot between
// standing, a pour, two strokes of the axe, or two of swimming.
func (b svgBody) Clip(anim, dir string) Clip {
	switch anim {
	case AnimWalk:
		return b.clip(8, dir+"1", dir, dir+"2", dir)
	case AnimAct:
		return b.clip(1, dir+"Act")
	case AnimChop:
		return b.clip(2/actFor, dir+"Chop", dir+"Chop1")
	case AnimSwim:
		return b.clip(3, dir+"Swim", dir+"Swim1")
	}
	return b.clip(1, dir)
}

// clip is the frames names shown fps a second.
func (b svgBody) clip(fps float64, names ...string) Clip {
	frames := make([]*sprite.Sprite, len(names))
	for i, n := range names {
		frames[i] = b.sprite(n)
	}
	return Clip{Frames: frames, FPS: fps}
}

// sprite is the frame name, standing on the person's feet, drawn at half
// the size of its art, smaller for the young, when it is first shown.
func (b svgBody) sprite(name string) *sprite.Sprite {
	s := b.size
	if s <= 0 {
		s = 1
	}
	w, h := int(art.FarmerW/2*s+0.5), int(art.FarmerH/2*s+0.5)
	return &sprite.Sprite{PivotX: 0.5, PivotY: 1, Make: func() *canvas.Canvas {
		return picture(b.frame(name), w, h)
	}}
}
