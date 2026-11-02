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

// Each motion is baked facing eight ways, a frame at a time.
const (
	Facings = 8
	Frames  = 8
)

// water is the pond a swimmer is baked in: dark, clear a hand deep.
var water = &bake.Water{Color: [3]float64{0.01, 0.06, 0.07}, Clear: 0.35}

// stageKey names a villager at a stage of life.
type stageKey struct {
	id    string
	stage life.Stage
	tag   string // of the character the player made, if the villager is the player
}

// stageBody is the model of a villager at a stage, animated.
type stageBody struct{ model *gltf.Model }

// PersonName is the name of the frames of villager id at stage, in the
// motion: person_<id>_<stage>_<motion>, then _d<facing>_f<frame>.
func PersonName(id string, stage life.Stage, motion string) string {
	return fmt.Sprintf("person_%s%s_%s_%s", id, kit.SpriteTag(id), stage, motion)
}

// look is a villager at a stage of life, in a motion, facing a way: what
// a task of people bakes, its eight frames.
type look struct {
	v      kit.Villager
	g      kit.Grown
	mo     human.Motion
	facing int
}

// peopleTasks bake each villager at each stage of life, a motion at a
// time and a facing at a time. Core are the grown-up player in all of
// their motions, and the other grown-ups standing and walking.
func (b *Baker) peopleTasks() []Task {
	var out []Task
	for _, v := range kit.Villagers {
		for _, g := range v.Stages() {
			for _, mo := range g.Motions {
				core := g.Stage == life.Adult && (v.ID == "player" || mo.Name == human.Stand.Name || mo.Name == human.Walk.Name)
				for d := range Facings {
					l := look{v, g, mo, d}
					name := fmt.Sprintf("person %s %s %s %s %d%s", versionOf(mo), v.ID, g.Stage, mo.Name, d, kit.PlayerTag(v))
					out = append(out, Task{Name: name, Kind: "person", ID: v.ID, Core: core, Run: func() error { return b.person(l) }})
				}
			}
		}
	}
	return out
}

// person bakes the frames of l, on every core.
func (b *Baker) person(l look) error {
	body, err := b.bodyOf(l.v, l.g)
	if err != nil {
		return fmt.Errorf("%s as a %s: %w", l.v.ID, l.g.Stage, err)
	}
	frames := framesOf(l)
	return b.each(len(frames), func(i int) error {
		return b.Shelf.Put(frames[i].name, shoot(body.model, frames[i]))
	})
}

// frame is one picture of a villager to bake: its name and shot.
type frame struct {
	name string
	shot bake.Shot
	wet  bool
}

// framesOf are every frame of l.
func framesOf(l look) []frame {
	var out []frame
	prefix := PersonName(l.v.ID, l.g.Stage, l.mo.Name)
	for k := range Frames {
		out = append(out, frame{
			name: fmt.Sprintf("%s_d%d_f%d", prefix, l.facing, k),
			shot: bake.Shot{Animation: l.mo.Name, Time: timeOf(l.mo, k), Yaw: float64(l.facing) * 360 / Facings},
			wet:  slices.Contains(wetMotions, l.mo.Name),
		})
	}
	return out
}

// wetMotions are the motions in the water, baked with it round them.
var wetMotions = []string{human.Swim.Name, human.Tread.Name, human.Dive.Name, human.Rise.Name}

// timeOf is the moment of motion m that frame k of a bake shows: a loop's
// frames are evenly round it; a motion run through once shows its first
// and its last, which poses just short of the loop's wrap to the start.
func timeOf(m human.Motion, k int) float64 {
	if m.Once {
		return m.Length * float64(k) / (Frames - 1) * (1 - 1e-9)
	}
	return m.Length * float64(k) / Frames
}

// shoot bakes frame f of m, in the water if f is wet.
func shoot(m *gltf.Model, f frame) bake.Sprite {
	o := Options()
	if f.wet {
		o.Water = water
	}
	return bake.Render(m, f.shot, o)
}

// versionOf is the version of the motion mo as baked: the motions', but
// the chop's, the forge's, the hoe's and the pour's own, so each alone is
// baked again when it changes.
func versionOf(mo human.Motion) string {
	switch mo.Name {
	case human.Chop.Name:
		return fmt.Sprintf("v%dc%d", human.MotionsVersion, human.ChopVersion)
	case human.Forge.Name:
		return fmt.Sprintf("v%df%d", human.MotionsVersion, human.ForgeVersion)
	case human.Hoe.Name:
		return fmt.Sprintf("v%dh%d", human.MotionsVersion, human.HoeVersion)
	case human.Water.Name:
		return fmt.Sprintf("v%dw%d", human.MotionsVersion, human.WaterVersion)
	case human.Pick.Name:
		return fmt.Sprintf("v%dp%d", human.MotionsVersion, human.PickVersion)
	}
	return fmt.Sprintf("v%d", human.MotionsVersion)
}
