package farm

import (
	"math"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/random"
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
}

// herd is the animals of the farm's pen: each wanders inside it, walks to
// a spot, stands, and walks again, each drawn as an object of the game.
type herd struct {
	f     *farm
	art   AnimalArt
	luck  *random.Source
	field rect
	crits []*creature
}

// newHerd is the farm's herd with its first cow, or nil when the farm's
// art has no animals to draw.
func newHerd(f *farm) *herd {
	art, ok := f.art.(AnimalArt)
	if !ok {
		return nil
	}
	h := &herd{f: f, art: art, luck: random.New(seed + 47), field: penRect.inside()}
	h.add("cow", 0)
	return h
}

// add puts an animal of kind and coat in the pen, at a spot of it.
func (h *herd) add(kind string, coat int) {
	c := &creature{kind: kind, coat: coat, wait: 1 + 3*h.luck.Float(), heading: 2}
	c.pos = h.spot()
	c.to = c.pos
	h.crits = append(h.crits, c)
}

// spot is a place in the pen's field, a half cell in from its fence.
func (h *herd) spot() vector.Vector3 {
	r := h.field
	x := (float64(r.col) + 0.5 + float64(r.cols-1)*h.luck.Float()) * tile
	y := (float64(r.row) + 0.5 + float64(r.rows-1)*h.luck.Float()) * tile
	return vector.Vector3{X: x, Y: y}
}

// Update lets the animals live, if the farmer is near the pen and out of
// doors, and dresses each as it is.
func (h *herd) Update(o *game.Object, dt float64) {
	near := h.f.room == nil && h.f.player.Transform.Position.Sub(h.f.ground.CellCenter(h.field.col+h.field.cols/2, h.field.row+h.field.rows/2)).Length() < herdReach*tile
	for _, c := range h.crits {
		if c.obj == nil {
			c.obj = game.NewObject("animal")
			c.obj.Z = zThings
			o.Game().Add(c.obj)
		}
		c.obj.Hidden = !near
		if !near {
			continue
		}
		h.live(c, dt)
		h.dress(c)
	}
}

// live moves c: it waits, then walks to a new spot, grazing now and then
// as it stands.
func (h *herd) live(c *creature, dt float64) {
	c.time += dt
	if !c.walking {
		c.wait -= dt
		if c.wait <= 0 {
			c.to, c.walking = h.spot(), true
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
		look.Motion, look.Frame = "walk", int(c.time*fps)
	case c.grazes:
		look.Motion = "graze"
	}
	c.obj.Transform.Position = c.pos
	if s := h.art.Animal(look); s != nil && s.Image != nil {
		c.obj.Sprite = s
	} else {
		c.obj.Hidden = true
	}
}
