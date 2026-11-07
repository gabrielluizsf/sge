package main

import (
	"os"

	"github.com/gabrielluizsf/sge/animal"
)

// theCow is how the cow walks: a lateral sequence, a foot after another.
var theCow = animal.Gait{Stride: 0.26, Period: 1.1, Duty: 0.65, Lift: 0.045}

// buildCow reads Rukh3D's cow, which has bones but no walk, and writes the
// walk on them. Its grazing was left out: the head did not lower.
func buildCow(root *os.Root) ([]coat, error) {
	m, err := open(root, "farm_cow_animated_dairy_cattle")
	if err != nil {
		return nil, err
	}
	r, err := animal.NewRig(m)
	if err != nil {
		return nil, err
	}
	animal.Wrap(m, animal.FrameFor(m, r.Toward(r.Bone("spine_01"), r.Bone("head_012"))))
	if r, err = animal.NewRig(m); err != nil {
		return nil, err
	}
	legs := []animal.Leg{
		r.NewLeg([]string{"lfl1_017", "lfl2_018"}, []string{"twf.L_019", "feetf.L_020"}, "feetf.L_end_045").At(1, 0),
		r.NewLeg([]string{"lfr1_021", "lfr2_022"}, []string{"twf.R_023", "feetf.R_024"}, "feetf.R_end_046").At(1, 0.5),
		r.NewLeg([]string{"lbl1_02", "lbl2_03"}, []string{"twb.L_04", "feetb.L_00"}, "feetb.L_end_040").At(-1, 0.75),
		r.NewLeg([]string{"lbr1_031", "lbr2_032"}, []string{"twb.R_033", "feetb.R_034"}, "feetb.R_end_048").At(-1, 0.25),
	}
	times := loop(6, theCow.Period)
	r.Clip("walk", times, func(i int) *animal.Pose { return animal.Walk(legs, theCow, times[i], 0.015) })
	return []coat{{name: "dairy", model: m, motions: []motion{
		{name: "stand", times: []float64{0}},
		{name: "walk", clip: "walk", times: times, speed: theCow.Speed() * 2.6},
	}}}, nil
}
