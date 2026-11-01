package bakes

import (
	"fmt"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/gltf"
)

// modelTasks bake each model of the manifest, a thing standing on the
// ground: trees, rocks, plants and the town's props.
func (b *Baker) modelTasks() []Task {
	var out []Task
	for _, a := range b.Manifest.Assets {
		if a.Kind == "model" {
			out = append(out, Task{Name: "model " + a.ID, Kind: "model", ID: a.ID, Core: true, Run: func() error { return b.model(a) }})
		}
	}
	return out
}

// model bakes asset a at every heading, or each of its variants at two.
func (b *Baker) model(a kit.Asset) error {
	model, err := Load(a)
	if err != nil {
		return fmt.Errorf("%s: %w", a.ID, err)
	}
	for _, sh := range shots(model, a.Whole) {
		if err := b.Shelf.Put(a.ID+sh.suffix, bake.Render(model, sh.Shot, Options())); err != nil {
			return err
		}
	}
	return nil
}

// shot is a picture to bake of a model, and the suffix of its name.
type shot struct {
	bake.Shot
	suffix string
}

// shots are the pictures baked of m: each heading of the whole model, if
// it has one root or is whole, its roots the parts of one thing; or, for
// a model of variants side by side, each variant alone at two.
func shots(m *gltf.Model, whole bool) []shot {
	var out []shot
	if len(m.Roots) == 1 || whole {
		for i, yaw := range Headings {
			out = append(out, shot{bake.Shot{Yaw: yaw}, fmt.Sprintf("_%d", i)})
		}
		return out
	}
	for r := range m.Roots {
		for i, yaw := range Headings[:2] {
			out = append(out, shot{bake.Shot{Yaw: yaw, Only: true, Root: r}, fmt.Sprintf("_v%d_%d", r, i)})
		}
	}
	return out
}
