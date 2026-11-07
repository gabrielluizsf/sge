package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/farm/art"
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
// standing, a pour, two strokes of the axe, two of swimming, or the two that go into the
// water or come out of it.
func (b svgBody) Clip(anim, dir string) Clip {
	switch anim {
	case AnimWalk:
		return b.clip(8, dir+"1", dir, dir+"2", dir)
	case AnimRun:
		return b.clip(13, dir+"1", dir, dir+"2", dir)
	case AnimAct:
		return b.clip(1, dir+"Act")
	case AnimChop:
		return b.clip(2/actFor, dir+"Chop", dir+"Chop1")
	case AnimMine:
		return b.clip(2/MineSeconds, dir+"Mine", dir+"Mine1")
	case AnimEat:
		return b.clip(5/EatSeconds, dir+"Eat1", dir+"Eat", dir+"Eat", dir+"Eat1", dir+"Eat")
	case AnimDrink:
		return b.clip(4/DrinkSeconds, dir+"Drink", dir+"Drink1", dir+"Drink1", dir+"Drink")
	case AnimMount:
		return b.rising(dir, false)
	case AnimRide:
		return b.raised(4, dir, dir+"1", dir, dir+"2")
	case AnimDismount:
		return b.rising(dir, true)
	case AnimWell:
		return b.clip(6/WellSeconds, dir+"Well", dir+"Well1", dir+"Well", dir+"Well1", dir+"Well2", dir+"Well2")
	case AnimDip:
		return b.clip(5/DipSeconds, dir, dir+"Dip", dir+"Dip", dir+"Dip1", dir)
	case AnimMilkDown:
		return b.clip(2/MilkDownSeconds, dir, dir+"Milk")
	case AnimMilk:
		return b.clip(5, dir+"Milk", dir+"Milk1")
	case AnimMilkUp:
		return b.clip(2/MilkDownSeconds, dir+"Milk", dir)
	case AnimSleep, AnimSleepLeft, AnimSleepRight:
		return b.clip(0.5, "Sleep", "Sleep1")
	case AnimHoe, AnimWater, AnimPick:
		return b.clip(1, dir+"Act")
	case AnimForge:
		return b.clip(2/ForgeSeconds, dir+"Chop", dir+"Chop1")
	case AnimSwim:
		return b.clip(3, dir+"Swim", dir+"Swim1")
	case AnimDive:
		return b.clip(2/DiveSeconds, dir, dir+"Swim")
	case AnimRise:
		return b.clip(2/RiseSeconds, dir+"Swim", dir)
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

// raised is a clip of the frames names, fps a second, drawn up on a horse's
// back, in a look that draws the horse and rider apart: their legs fall down
// its flank.
func (b svgBody) raised(fps float64, names ...string) Clip {
	c := b.clip(fps, names...)
	for _, s := range c.Frames {
		b.lift(s, 1)
	}
	return c
}

// rising is a clip of getting up on a horse's back, a step at a time, facing
// dir, or down from it.
func (b svgBody) rising(dir string, down bool) Clip {
	names := []string{dir, dir + "1", dir + "2", dir}
	c := b.clip(float64(len(names))/MountSeconds, names...)
	for i, s := range c.Frames {
		u := float64(i) / float64(len(names)-1)
		if down {
			u = 1 - u
		}
		b.lift(s, u)
	}
	return c
}

// lift raises the frame s from the ground, by u of the way up a horse.
func (b svgBody) lift(s *sprite.Sprite, u float64) {
	_, h := b.drawn()
	s.PivotY = 1 + u*rideLift/float64(h)
}

// rideLift is how far above the ground, in pixels, a rider's feet are drawn:
// their hips at the height of the back of the horse the look draws.
const rideLift = 8 * scale

// sprite is the frame name, standing on the person's feet, drawn at half
// the size of its art, smaller for the young, when it is first shown.
func (b svgBody) sprite(name string) *sprite.Sprite {
	w, h := b.drawn()
	return &sprite.Sprite{PivotX: 0.5, PivotY: 1, Make: func() *canvas.Canvas {
		return picture(b.frame(name), w, h)
	}}
}

// drawn is how many pixels of art wide and tall the person is drawn.
func (b svgBody) drawn() (w, h int) {
	s := b.size
	if s <= 0 {
		s = 1
	}
	return int(art.FarmerW/2*s + 0.5), int(art.FarmerH/2*s + 0.5)
}

// large is the frame name drawn from its art zoom times larger than in the
// world: sharp, as enlarging a small picture is not.
func (b svgBody) large(name string, zoom int) *canvas.Canvas {
	w, h := b.drawn()
	return picture(b.frame(name), w*zoom, h*zoom)
}
