package farm

import (
	"math"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// herdReach is how far from the farmer, in cells, the herd is lived
// frame by frame: further off its animals stand still and are not drawn.
const herdReach = 24

// creature is an animal of the pen.
type creature struct {
	kind    string
	coat    int
	pos, to vector.Vector3
	walking bool
	wait    float64 // seconds until it sets off again
	grazes  bool
	time    float64
	heading int
	obj     *game.Object
	// milked is one more than the day she was last milked, 0 for never,
	// and held says she stands still, being milked.
	milked int
	held   bool
	pen    *pen
	// hurt is the harm the player's blows have done it.
	hurt float64
	// homing says it is walking back to its pen by route, its rider got down.
	homing bool
	route  *pathfind.Route
	// mounted says it has a rider, going at rideSpeed; walkSpeed is its own.
	mounted              bool
	rideSpeed, walkSpeed float64
}

// pen is a fenced field of animals of a kind, on the farm or, with town
// set, at that town's stable.
type pen struct {
	kind string
	r    rect
	town *town.Town
}

// herd is the animals of the farm's pens and the stables' seen: each wanders
// inside its pen, walks to a spot, stands, and walks again, each drawn as an
// object of the game.
type herd struct {
	f      *farm
	art    AnimalArt
	luck   *random.Source
	pens   []*pen
	crits  []*creature
	stable map[*town.Town]bool
}

// newHerd is the farm's herd with its first cow, or nil when the farm's
// art has no animals to draw.
func newHerd(f *farm) *herd {
	art, ok := f.art.(AnimalArt)
	if !ok {
		return nil
	}
	h := &herd{f: f, art: art, luck: random.New(seed + 47), stable: map[*town.Town]bool{}}
	first := &pen{kind: "cow", r: penRect}
	h.pens = []*pen{first}
	h.add("cow", 0, first)
	return h
}

// add puts an animal of kind and coat in pen p, at a spot of it.
func (h *herd) add(kind string, coat int, p *pen) *creature {
	c := &creature{kind: kind, coat: coat, wait: 1 + 3*h.luck.Float(), heading: 2, pen: p}
	c.pos = h.spot(p)
	c.to = c.pos
	h.crits = append(h.crits, c)
	return c
}

// spot is a place in p's field, a half cell in from its fence.
func (h *herd) spot(p *pen) vector.Vector3 {
	r := p.r.inside()
	x := (float64(r.col) + 0.5 + float64(r.cols-1)*h.luck.Float()) * tile
	y := (float64(r.row) + 0.5 + float64(r.rows-1)*h.luck.Float()) * tile
	return vector.Vector3{X: x, Y: y}
}

// Update lets the animals live, if the farmer is near the pen and out of
// doors, and dresses each as it is.
func (h *herd) Update(o *game.Object, dt float64) {
	h.f.milkTick(dt)
	h.f.rideTick(dt)
	for _, c := range h.crits {
		near := h.f.room == nil && feet(h.f.player).Sub(c.pos).Length() < herdReach*tile
		if !near {
			if c.obj != nil {
				c.obj.Hidden = true
			}
			if c.homing {
				h.backInPen(c)
			}
			continue
		}
		if c.obj == nil {
			c.obj = game.NewObject("animal")
			c.obj.Z = zThings
			o.Game().Add(c.obj)
		}
		c.obj.Hidden = false
		h.live(c, dt)
		h.dress(c)
	}
}

// live moves c: it waits, then walks to a new spot, grazing now and then
// as it stands.
func (h *herd) live(c *creature, dt float64) {
	c.time += dt
	if c.held {
		return
	}
	if c.homing {
		h.goHome(c, dt)
		return
	}
	if !c.walking {
		c.wait -= dt
		if c.wait <= 0 {
			c.to, c.walking = h.spot(c.pen), true
		}
		return
	}
	d := c.to.Sub(c.pos)
	speed, _ := h.art.Pace(c.kind)
	step := speed * tile * dt
	if d.Length() <= step {
		c.pos, c.walking = c.to, false
		c.wait, c.grazes = 2+5*h.luck.Float(), h.luck.Float() < 0.5
		return
	}
	c.pos = c.pos.Add(d.Scale(step / d.Length()))
	c.heading = h.heading(d)
}

// heading is the eighth of a turn, clockwise from the screen's right,
// that an animal walking along dir on the ground faces on the screen.
func (h *herd) heading(dir vector.Vector3) int {
	s := dir
	if p := h.f.player.Game().Projection; p != nil {
		s = p.Project(dir)
	}
	return (int(math.Round(math.Atan2(s.Y, s.X)/(math.Pi/4)))%8 + 8) % 8
}

// dress puts c's object where it is, and gives it the sprite of what it does.
func (h *herd) dress(c *creature) {
	look := AnimalLook{Kind: c.kind, Coat: c.coat, Motion: "stand", Heading: c.heading}
	switch {
	case c.walking:
		_, fps := h.art.Pace(c.kind)
		look.Motion, look.Frame = "walk", int(c.time*fps*c.pace())
	case c.grazes:
		look.Motion = "graze"
	}
	c.obj.Transform.Position = c.pos
	if c.held {
		c.obj.Transform.Position.Y -= heldBehind // whoever is at it is drawn over it
		c.obj.Transform.Position.X -= heldBehind
	}
	if s := h.art.Animal(look); s != nil && s.Image != nil {
		c.obj.Sprite = s
	} else {
		c.obj.Hidden = true
	}
}
