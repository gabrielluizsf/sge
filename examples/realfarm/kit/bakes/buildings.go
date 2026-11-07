package bakes

import (
	"fmt"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/building"
)

// buildingTasks bake each house and thing package building builds, clad
// in the manifest's materials; the houses and the stall are core, the
// other props after.
func (b *Baker) buildingTasks() []Task {
	var out []Task
	for _, k := range building.All() {
		core := k.Name == "stall" || k.Name == "anvil" || k.Name == "forge" || k.Name == "stump" || k.Name == "stove" || k.Name == "stable" || IsHouse(k.Name)
		out = append(out, Task{Name: "building " + k.Name, Kind: "building", ID: k.Name, Core: core, Run: func() error { return b.building(k) }})
	}
	return out
}

// IsHouse reports whether name is one of the houses.
func IsHouse(name string) bool {
	for _, h := range building.Houses {
		if h.Name == name {
			return true
		}
	}
	return false
}

// building bakes k at every heading, a core each: building_<name>_<i>.
func (b *Baker) building(k building.Building) error {
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	model := k.Model(palette)
	return b.each(len(Headings), func(i int) error {
		s := bake.Render(model, bake.Shot{Yaw: Headings[i]}, Options())
		return b.Shelf.Put(fmt.Sprintf("building_%s_%d", k.Name, i), s)
	})
}
