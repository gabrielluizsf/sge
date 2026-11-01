package bakes

import (
	"fmt"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
)

// groundTasks bake each ground of the manifest: <id>_g0 to _g3.
func (b *Baker) groundTasks() []Task {
	var out []Task
	for _, a := range b.Manifest.Assets {
		if a.Kind == "ground" {
			out = append(out, Task{Name: "ground " + a.ID, Kind: "ground", ID: a.ID, Core: true, Run: func() error { return b.ground(a) }})
		}
	}
	return out
}

// ground bakes the four patches of ground a.
func (b *Baker) ground(a kit.Asset) error {
	model, err := Load(a)
	if err != nil {
		return fmt.Errorf("%s: %w", a.ID, err)
	}
	for i, p := range Patches {
		if err := b.Shelf.Put(fmt.Sprintf("%s_g%d", a.ID, i), bake.Ground(model, p, Options())); err != nil {
			return err
		}
	}
	return nil
}
