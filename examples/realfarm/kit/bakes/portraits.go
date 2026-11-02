package bakes

import (
	"fmt"
	"image"
	"math"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/life"
)

// PortraitName is the sprite of villager id's portrait at stage:
// portrait_<id>_<stage>.
func PortraitName(id string, stage life.Stage) string {
	return fmt.Sprintf("portrait_%s_%s", id, stage)
}

// portraitZoom is how many times larger than in the world a portrait is
// rendered: the head and shoulders fill a panel.
const portraitZoom = 12

// portraitFacing is the turn, in degrees, at which a villager looks at the
// player from the screen.
const portraitFacing = 0

// bustShare is how much of a villager's height, from the top of the head,
// a portrait of each stage of life shows: the head and shoulders, of a
// child more of them.
var bustShare = map[life.Stage]float64{life.Baby: 0.55, life.Child: 0.42, life.Teen: 0.32, life.Adult: 0.28, life.Elder: 0.28}

// portraitTasks bake each villager's head and shoulders at each stage of
// life, standing, from the same model as the sprites that walk the farm.
func (b *Baker) portraitTasks() []Task {
	var out []Task
	for _, v := range kit.Villagers {
		for _, g := range v.Stages() {
			name := fmt.Sprintf("portrait %s %s%s", v.ID, g.Stage, v.LookTag())
			out = append(out, Task{Name: name, Kind: "portrait", ID: v.ID, Run: func() error { return b.portrait(v, g) }})
		}
	}
	return out
}

// portrait bakes villager v as they are at g.
func (b *Baker) portrait(v kit.Villager, g kit.Grown) error {
	body, err := b.bodyOf(v, g)
	if err != nil {
		return fmt.Errorf("%s as a %s: %w", v.ID, g.Stage, err)
	}
	o := Options()
	o.Unit *= portraitZoom
	o.Shadow = 0
	s := bake.Render(leaned(body.model), bake.Shot{Animation: human.Stand.Name, Yaw: portraitFacing}, o)
	return b.Shelf.Put(PortraitName(v.ID, g.Stage), Bust(s, bustShare[g.Stage]))
}

// viewTilt is how far, in degrees, the game's view looks down on the
// world: a portrait leans the person back as much, to face the viewer.
const viewTilt = 35.264

// leaned is m under a new root that leans it back by viewTilt, so that the
// corner the game is seen from looks at it level; m itself is not changed.
func leaned(m *gltf.Model) *gltf.Model {
	c := *m
	half := -viewTilt * math.Pi / 360
	c.Nodes = append(append([]gltf.Node{}, m.Nodes...), gltf.Node{Name: "lean", Children: append([]int{}, m.Roots...),
		Mesh: -1, Skin: -1, Rotation: [4]float64{math.Sin(half), 0, 0, math.Cos(half)}, Scale: [3]float64{1, 1, 1}})
	c.Roots = []int{len(m.Nodes)}
	return &c
}

// Bust is the top share of the figure in s, its head and shoulders, cut
// to the figure's width there, but no wider than a little under its height,
// about the head, so the face is large in a portrait.
func Bust(s bake.Sprite, share float64) bake.Sprite { return Cut(s, BustRect(s, share)) }

// BustRect is the part of s Bust keeps; all of it when s is blank.
func BustRect(s bake.Sprite, share float64) image.Rectangle {
	b := s.Image.Bounds()
	top, bottom := -1, b.Min.Y
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if s.Image.NRGBAAt(x, y).A > 0 {
				if top < 0 {
					top = y
				}
				bottom = y + 1
			}
		}
	}
	if top < 0 {
		return b
	}
	cut := top + max(1, int(float64(bottom-top)*share))
	left, right := b.Max.X, b.Min.X
	for y := top; y < cut; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if s.Image.NRGBAAt(x, y).A > 0 {
				left, right = min(left, x), max(right, x+1)
			}
		}
	}
	r := image.Rect(left, top, right, cut)
	if w := r.Dy() * 9 / 10; r.Dx() > w {
		mid := headCentre(s.Image, image.Rect(left, top, right, top+(cut-top)/3))
		x0 := max(left, min(mid-w/2, right-w))
		r = image.Rect(x0, top, x0+w, cut)
	}
	return r
}

// Cut is the part r of s, its pivot at its bottom middle.
func Cut(s bake.Sprite, r image.Rectangle) bake.Sprite {
	out := image.NewNRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	for y := range r.Dy() {
		for x := range r.Dx() {
			out.SetNRGBA(x, y, s.Image.NRGBAAt(r.Min.X+x, r.Min.Y+y))
		}
	}
	return bake.Sprite{Image: out, Pivot: image.Pt(r.Dx()/2, r.Dy())}
}

// headCentre is the middle, across, of what is drawn inside r.
func headCentre(img *image.NRGBA, r image.Rectangle) int {
	sum, n := 0, 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if img.NRGBAAt(x, y).A > 0 {
				sum, n = sum+x, n+1
			}
		}
	}
	if n == 0 {
		return (r.Min.X + r.Max.X) / 2
	}
	return sum / n
}
