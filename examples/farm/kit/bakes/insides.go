package bakes

import (
	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/building"
)

// InsideYaw is the turn the pieces of a house's inside are baked at: the
// ground's, their sides along the map's rows and columns.
const InsideYaw = 225

// InsideName is the sprite piece name, one of building.Insides, is kept
// as.
func InsideName(name string) string { return "inside_" + name }

// insideTasks bake each piece of a house's inside, clad in the
// manifest's materials.
func (b *Baker) insideTasks() []Task {
	var out []Task
	for _, k := range building.Insides() {
		out = append(out, Task{Name: "inside " + k.Name, Kind: "inside", ID: k.Name, Core: true, Run: func() error { return b.inside(k) }})
	}
	return out
}

// inside bakes piece k standing on its cell, casting no shadow, for no
// sun reaches in; a floor, every piece but the walls, wholly opaque, so
// floors meet without a seam.
func (b *Baker) inside(k building.Building) error {
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	o := Options()
	o.Shadow = 0 // no sun reaches in
	floor := !isWall(k.Name)
	s := bake.Render(k.Model(palette), bake.Shot{Yaw: InsideYaw}, o)
	if floor {
		bake.Opaque(s)
	}
	return b.Shelf.Put(InsideName(k.Name), s)
}

// isWall reports whether the piece name is a back wall.
func isWall(name string) bool { return len(name) > 5 && name[:5] == "wall_" }
