package farm

import (
	"math"

	"github.com/gabrielluizsf/sge/flock"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// birdCount is how many birds live about the view, and birdBeat how many
// wing frames a second a flying one goes through.
const (
	birdCount = 5
	birdBeat  = 9
)

// birds are the farm's birds: a flock, perched on the trees, fences and
// roofs about the farmer and flying off when they or a villager come near,
// each drawn as an object of the game.
type birds struct {
	f      *farm
	art    BirdArt
	flock  *flock.Flock
	things []*game.Object
}

// newBirds is the farm's birds, or nil when its art has none to draw.
func newBirds(f *farm) *birds {
	art, ok := f.art.(BirdArt)
	if !ok {
		return nil
	}
	fl := flock.New(seed+31, birdCount)
	fl.Perches = f.perches
	fl.Around = func() vector.Vector3 { return f.player.Transform.Position }
	fl.Scares = f.scares
	return &birds{f: f, art: art, flock: fl}
}

// scares are where the farmer and the villagers outside stand: what
// startles a bird.
func (f *farm) scares() []vector.Vector3 {
	out := []vector.Vector3{f.player.Transform.Position}
	for _, v := range f.people {
		if v.obj != nil && !v.inside {
			out = append(out, v.obj.Transform.Position)
		}
	}
	return out
}

// Update lets the flock live, and puts each bird where it is, as it
// looks: hidden while the farmer is indoors.
func (b *birds) Update(o *game.Object, dt float64) {
	if b.f.room == nil {
		b.flock.Step(dt)
	}
	for len(b.things) < len(b.flock.Birds()) {
		t := game.NewObject("bird")
		t.Z = zThings
		t.Footing = vector.Vector3{X: tile, Y: tile} // in front of the tree it sits in
		o.Game().Add(t)
		b.things = append(b.things, t)
	}
	for i, bird := range b.flock.Birds() {
		t := b.things[i]
		t.Hidden = b.f.room != nil
		t.Transform.Position = bird.Pos
		b.dress(t, bird)
	}
}

// dress gives t the sprite of bird, lifted by its height.
func (b *birds) dress(t *game.Object, bird *flock.Bird) {
	look := BirdLook{Flying: bird.State == flock.Flying, Heading: b.heading(bird.Dir), Frame: int(bird.Time*birdBeat) % BirdFrames}
	s := b.art.Bird(look)
	if s == nil || s.Image == nil {
		t.Hidden = true
		return
	}
	lifted := *s
	lifted.PivotY += bird.Alt / float64(s.Image.Height)
	t.Sprite = &lifted
}

// heading is the eighth of a turn, clockwise from the screen's right,
// that a bird flying along dir on the ground faces on the screen.
func (b *birds) heading(dir vector.Vector3) int {
	s := dir
	if p := b.f.player.Game().Projection; p != nil {
		s = p.Project(dir)
	}
	return (int(math.Round(math.Atan2(s.Y, s.X)/(math.Pi/4)))%8 + 8) % 8
}
