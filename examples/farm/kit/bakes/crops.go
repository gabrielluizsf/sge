package bakes

import (
	"fmt"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/flora"
)

// cropSeeds are how many plants of each crop and stage are baked, so a
// field of them is not one plant over and over.
const cropSeeds = 3

// cropTasks bake each crop flora grows.
func (b *Baker) cropTasks() []Task {
	var out []Task
	for _, c := range flora.Crops {
		out = append(out, Task{Name: "crop " + c.Name, Kind: "crop", ID: c.Name, Core: true, Run: func() error { return b.crop(c) }})
	}
	return out
}

// crop bakes crop c at every stage, a few plants of each:
// crop_<name>_s<stage>_<seed>.
func (b *Baker) crop(c flora.Crop) error {
	for stage := range c.Stages {
		for seed := range cropSeeds {
			s := bake.Render(c.Grow(stage, uint32(seed+1)), bake.Shot{Yaw: float64(seed) * 120}, Options())
			if err := b.Shelf.Put(fmt.Sprintf("crop_%s_s%d_%d", c.Name, stage, seed), s); err != nil {
				return err
			}
		}
	}
	return nil
}
