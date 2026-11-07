package main

import (
	"math"

	"github.com/gabrielluizsf/sge/animal"
	"github.com/gabrielluizsf/sge/gltf"
)

// beast is how an animal with no skeleton is given one: bones along its
// body, in the model's axes once it faces +Z and stands on Y=0, and four
// legs, each two bars from the hip to the knee and the knee to the foot.
type beast struct {
	body    []animal.Bone // the torso, neck, head and tail, parents first
	feet    [4]animal.Vec // front left, front right, back left, back right
	hipY    float64       // how high the hips are
	knee    [2]float64    // how far the knees stand from the leg's line, front and back
	chest   string        // the bone the front legs hang from; the back legs hang from "root"
	gait    animal.Gait
	crouch  float64
	forward animal.Vec // the way the model faces as it is made
}

var legNames = [4]string{"fl", "fr", "bl", "br"}

// phases make a walk of one foot after another: front left, back right,
// front right, back left.
var phases = [4]float64{0, 0.5, 0.75, 0.25}

// bones are the body's, then each leg's three.
func (b beast) bones() []animal.Bone {
	out := append([]animal.Bone(nil), b.body...)
	for i, foot := range b.feet {
		parent, bend, off := "root", -1.0, b.knee[1]
		if i < 2 {
			parent, bend, off = b.chest, 1.0, b.knee[0]
		}
		hip := animal.Vec{foot[0], b.hipY, foot[2]}
		knee := animal.Vec{foot[0], b.hipY / 2, foot[2] + bend*off}
		n := legNames[i]
		out = append(out,
			animal.Bone{Name: n + "_up", Parent: parent, Head: hip, Tail: knee},
			animal.Bone{Name: n + "_low", Parent: n + "_up", Head: knee, Tail: foot},
			animal.Bone{Name: n + "_toe", Parent: n + "_low", Head: foot, Tail: foot},
		)
	}
	return out
}

// rig stands m, a model with no skeleton, in the frame and gives it the
// beast's, returning the rig and its four legs.
func (b beast) rig(m *gltf.Model) (*animal.Rig, []animal.Leg, error) {
	animal.Bake(m, animal.FrameFor(m, b.forward))
	r, err := animal.Skeleton(m, b.bones())
	if err != nil {
		return nil, nil, err
	}
	legs := make([]animal.Leg, 4)
	for i, n := range legNames {
		bend := -1.0
		if i < 2 {
			bend = 1
		}
		legs[i] = r.NewLeg([]string{n + "_up"}, []string{n + "_low"}, n+"_toe").At(bend, phases[i])
	}
	return r, legs, nil
}

// walk writes the walk, six frames, into the rig, and the motion that
// plays it.
func (b beast) walk(r *animal.Rig, legs []animal.Leg, unit float64) motion {
	times := loop(6, b.gait.Period)
	r.Clip("walk", times, func(i int) *animal.Pose { return animal.Walk(legs, b.gait, times[i], b.crouch) })
	return motion{name: "walk", clip: "walk", times: times, speed: b.gait.Speed() * unit}
}

// graze writes the head lowered to the grass, in two frames.
func graze(r *animal.Rig, legs []animal.Leg, neck, head string, down float64) motion {
	times := loop(2, 1)
	x := animal.Vec{1, 0, 0}
	r.Clip("graze", times, func(i int) *animal.Pose {
		p := animal.NewPose()
		bob := 1 + 0.12*float64(i)
		p.Turn(r.Bone(neck), animal.About(x, down*0.55*bob))
		p.Turn(r.Bone(head), animal.About(x, down*bob))
		for _, l := range legs {
			l.Reach(p, l.Foot)
		}
		return p
	})
	return motion{name: "graze", clip: "graze", times: times}
}

func rad(deg float64) float64 { return deg * math.Pi / 180 }
