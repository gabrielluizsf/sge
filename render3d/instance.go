package render3d

import (
	"math"

	"github.com/gabrielluizsf/sge/pose"
)

// Instance is a model shown somewhere: its chain of detail, Place, where
// in the game's world it stands (a translation, rotation and scale, as At
// makes), and the Animation it plays (by name, "" for none) and how many
// seconds into it, looped.
type Instance struct {
	Model     *LOD
	Place     pose.Mat
	Animation string
	Time      float64

	level int
	anim  animRef
}

// animRef is an animation's name looked up in a model, kept until either
// changes.
type animRef struct {
	model *Model
	name  string
	index int
}

// Level is the level of detail the instance was last drawn at.
func (in *Instance) Level() int { return in.level }

// animation is the index of the instance's animation in m, or -1.
func (in *Instance) animation(m *Model) int {
	if in.anim.model != m || in.anim.name != in.Animation {
		in.anim = animRef{model: m, name: in.Animation, index: m.poser.Animation(in.Animation)}
	}
	return in.anim.index
}

// At is the place of a model standing at pos, turned yaw degrees round
// the vertical, counterclockwise seen from above, and scaled by scale.
func At(pos [3]float64, yaw, scale float64) pose.Mat {
	a := yaw * math.Pi / 180
	h := math.Sin(a / 2)
	return pose.TRS(pos, [4]float64{0, 0, h, math.Cos(a / 2)}, [3]float64{scale, scale, scale})
}
