package main

import "github.com/gabrielluizsf/sge/animal"

// wolfBite is the wolf lunging at what it bites, in four frames: crouched,
// leaping, biting, and drawn back.
func wolfBite(r *animal.Rig, legs []animal.Leg) motion {
	type frame struct {
		shift                   animal.Vec
		pitch, neck, head, lift float64
	}
	frames := []frame{
		{animal.Vec{0, -0.5, -0.3}, 8, 10, 5, 0},
		{animal.Vec{0, 0.4, 0.9}, -12, -10, -5, 0.8},
		{animal.Vec{0, 0.0, 1.8}, 10, 20, 25, 0.3},
		{animal.Vec{0, -0.1, 0.6}, 0, 5, 0, 0},
	}
	times := run(4, 0.6)
	x := animal.Vec{1, 0, 0}
	r.Clip("bite", times, func(i int) *animal.Pose {
		f := frames[i]
		p := animal.NewPose()
		p.Shift = f.shift
		p.Turn(r.Bone("root"), animal.About(x, rad(f.pitch)))
		p.Turn(r.Bone("neck"), animal.About(x, rad(f.pitch+f.neck)))
		p.Turn(r.Bone("head"), animal.About(x, rad(f.pitch+f.neck+f.head)))
		for _, l := range legs {
			l.Reach(p, animal.Vec{l.Foot[0], l.Foot[1] + f.lift, l.Foot[2] + f.shift[2]*0.6})
		}
		return p
	})
	return motion{name: "bite", clip: "bite", times: times, once: true}
}

// wolfFall is the wolf dying, in four frames: it staggers, rolls onto its
// side, and lies there.
func wolfFall(r *animal.Rig, legs []animal.Leg) motion {
	times := run(4, 1.0)
	steps := []float64{0, 0.3, 0.65, 1}
	z := animal.Vec{0, 0, 1}
	r.Clip("fall", times, func(i int) *animal.Pose {
		s := steps[i]
		p := animal.NewPose()
		p.Turn(r.Bone("root"), animal.About(z, -rad(90)*s))
		p.Shift = animal.Vec{0, -(3.8 - 1.1) * s, 0}
		p.Turn(r.Bone("neck"), animal.About(z, -rad(90)*s+rad(15)*s))
		p.Turn(r.Bone("head"), animal.About(z, -rad(90)*s+rad(30)*s))
		return p
	})
	return motion{name: "fall", clip: "fall", times: times, once: true}
}
