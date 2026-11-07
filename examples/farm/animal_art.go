package main

import (
	"bytes"
	"image/png"
	"sync"

	"github.com/gabrielluizsf/sge/examples/farm/animals"
	"github.com/gabrielluizsf/sge/examples/farm/kit"
	"github.com/gabrielluizsf/sge/farm"
	"github.com/gabrielluizsf/sge/sprite"
)

// theAnimals are the baked animals' motions and places, read once, and
// the sprites made of them, each when it is first asked for.
var theAnimals struct {
	once    sync.Once
	infos   map[string]animals.Info
	places  map[string]animals.Sprite
	sprites map[string]*sprite.Sprite
}

func loadAnimals() {
	theAnimals.once.Do(func() {
		theAnimals.infos, _ = animals.Infos()
		theAnimals.places, _ = animals.Sprites()
		theAnimals.sprites = map[string]*sprite.Sprite{}
	})
}

// Animal is the baked animal as look says, in its coat, motion and
// facing; a motion it has no frames of, grazing for the cow, is shown as
// standing.
func (a *realArt) Animal(look farm.AnimalLook) *sprite.Sprite {
	loadAnimals()
	in, ok := theAnimals.infos[look.Kind]
	if !ok || len(in.Coats) == 0 {
		return nil
	}
	coat := in.Coats[(look.Coat%len(in.Coats)+len(in.Coats))%len(in.Coats)]
	motion, ok := in.Motions[look.Motion]
	if !ok || motion.Frames == 0 {
		motion = in.Motions["stand"]
		look.Motion = "stand"
	}
	frame := 0
	if motion.Frames > 0 {
		frame = (look.Frame%motion.Frames + motion.Frames) % motion.Frames
	}
	name := animals.Frame(look.Kind, coat, look.Motion, (look.Heading+turnOffset)%8, frame)
	if s, ok := theAnimals.sprites[name]; ok {
		return s
	}
	at, ok := theAnimals.places[name]
	data, err := animals.PNG(name)
	if !ok || err != nil || at.W == 0 {
		return nil
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	s, err := sprite.FromImage(img)
	if err != nil {
		return nil
	}
	s.PivotX, s.PivotY = float64(at.PivotX)/float64(at.W), float64(at.PivotY)/float64(at.H)
	theAnimals.sprites[name] = s
	return s
}

// Milker is by the hind legs of the cow as she is baked, a cell being a
// metre, so the hands of the one sitting there, measured on their pictures
// and hers, come to her teats under her belly, she in profile.
func (a *realArt) Milker() (back, near float64) { return 1.75 * kit.Cell, 0.1 * kit.Cell }

// Frames is how many pictures of the motion were baked, at the least one.
func (a *realArt) Frames(kind, motion string) int {
	loadAnimals()
	return max(theAnimals.infos[kind].Motions[motion].Frames, 1)
}

// Pace is the speed the animal was baked to walk at, and the frames a
// second that plays its loop in its own time.
func (a *realArt) Pace(kind string) (speed, fps float64) {
	loadAnimals()
	m := theAnimals.infos[kind].Motions["walk"]
	if m.Length <= 0 {
		return 0.5, 6
	}
	return m.Speed, float64(m.Frames) / m.Length
}
