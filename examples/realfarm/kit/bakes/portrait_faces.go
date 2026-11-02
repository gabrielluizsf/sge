package bakes

import (
	"fmt"
	"slices"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/life"
)

// PortraitFace is a face a portrait is baked with besides its still one:
// its name, and the expression of MakeHuman's face pose units it is.
type PortraitFace struct {
	Name string
	Face human.Expression
}

// PortraitFaces are the faces a villager talks with: the eyelids half
// shut and shut, for a blink, and the mouth ajar, open and round, as
// letters are said.
var PortraitFaces = []PortraitFace{
	{"lids1", human.Lids(0.55)}, {"lids2", human.Lids(1)},
	{"say1", human.Mouth(0.25, 0)}, {"say2", human.Mouth(0.55, 0)}, {"say3", human.Mouth(0.3, 0.8)},
}

// facesVersion counts the changes in how the talking faces are made that
// bake them all again: 2 shaded the teeth as the inside of a mouth.
const facesVersion = 2

// PortraitFaceName is the sprite of villager id's portrait at stage with
// the face called face.
func PortraitFaceName(id string, stage life.Stage, face string) string {
	return PortraitName(id, stage) + "_" + face
}

// portraitFaceTasks bake each villager's talking faces at each stage of
// life, after their still portraits.
func (b *Baker) portraitFaceTasks() []Task {
	var out []Task
	for _, v := range kit.Villagers {
		for _, g := range v.Stages() {
			name := fmt.Sprintf("portrait face %s %s%s v%d", v.ID, g.Stage, v.LookTag(), facesVersion)
			out = append(out, Task{Name: name, Kind: "portrait", ID: v.ID, Run: func() error { return b.portraitFaces(v, g) }})
		}
	}
	return out
}

// portraitFaces bakes villager v's faces as they are at g, each cut where
// the still portrait is, so the face moves and the picture does not.
func (b *Baker) portraitFaces(v kit.Villager, g kit.Grown) error {
	makers, err := b.humanKit()
	if err != nil {
		return err
	}
	body, err := b.talker(makers, v, g)
	if err != nil {
		return fmt.Errorf("%s as a %s: %w", v.ID, g.Stage, err)
	}
	o := Options()
	o.Unit *= portraitZoom
	o.Shadow = 0
	shot := bake.Shot{Animation: human.Stand.Name, Yaw: portraitFacing}
	rect := BustRect(bake.Render(leaned(body), shot, o), bustShare[g.Stage])
	for _, f := range PortraitFaces {
		b.mu.Lock() // a kit is for one goroutine at a time
		m, err := makers.Express(body, f.Face)
		b.mu.Unlock()
		if err != nil {
			return err
		}
		if err := b.Shelf.Put(PortraitFaceName(v.ID, g.Stage, f.Name), Cut(bake.Render(leaned(m), shot, o), rect)); err != nil {
			return err
		}
	}
	return nil
}

// teeth are the teeth a talking mouth shows, MakeHuman's CC0 ones.
const teeth = "teeth/teeth_base/teeth_base.mhclo"

// talker is villager v as they are at g, with teeth in their mouth, which
// a shut one hides, so their still portrait is as it was.
func (b *Baker) talker(makers *human.Kit, v kit.Villager, g kit.Grown) (*gltf.Model, error) {
	b.mu.Lock() // a kit is for one goroutine at a time
	defer b.mu.Unlock()
	v.Look.Wear = append(slices.Clone(v.Look.Wear), teeth)
	if err := onlyCC0(makers, v); err != nil {
		return nil, err
	}
	m, err := makers.Make(v.At(g.Years), v.Look)
	if err == nil {
		err = human.Animate(m, g.Motions...)
	}
	return m, err
}
