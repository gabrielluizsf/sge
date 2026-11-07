package main

import (
	"os"
	"strings"

	"github.com/gabrielluizsf/sge/animal"
)

// horseCoats are the two horses honved444's model holds, side by side, a
// coat each: the one on the +X side, and the one on the -X.
var horseCoats = []struct {
	name string
	side float64
}{{"grey", 1}, {"chestnut", -1}}

// horseSpeed is the cells a second the model's own walk goes at, measured
// by -only stride: 2.96 model metres a second, a unit being 0.7 of a cell.
const horseSpeed = 2.07

// buildHorse reads the model once for each horse and keeps one: without
// the harness, which is black lines drawn across it, and without the nine
// rigid pieces the model places a hundred times too far.
func buildHorse(root *os.Root) ([]coat, error) {
	var out []coat
	for _, h := range horseCoats {
		m, err := open(root, "draft_horse_walk")
		if err != nil {
			return nil, err
		}
		animal.DropRigid(m)
		animal.Filter(m, func(material string, c animal.Vec) bool {
			return !strings.HasPrefix(material, "Default_OBJ") && c[0]*h.side > 0
		})
		r, err := animal.NewRig(m)
		if err != nil {
			return nil, err
		}
		animal.Wrap(m, animal.FrameFor(m, r.Toward(r.Bone("spine.001_00"), r.Bone("skull_015"))))
		out = append(out, coat{name: h.name, model: m, motions: []motion{
			{name: "stand", times: []float64{0}},
			{name: "walk", clip: "draft_Walk", times: loop(8, 1.3333), speed: horseSpeed},
		}})
	}
	return out, nil
}
