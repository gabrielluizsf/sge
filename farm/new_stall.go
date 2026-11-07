package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/town"
)

// newStall is a shop's stall, standing on one cell, at, among the people
// and the trees.
func (f *farm) newStall(at town.Cell) *game.Object {
	o := game.NewObject("stall")
	o.Z = zThings
	o.Sprite = f.art.Building("stall")
	o.Transform.Position = stallRect(f.art, at).middle()
	return o
}

// newShed is the shop's shed name, on the cells r.
func (f *farm) newShed(name string, r rect) *game.Object {
	o := game.NewObject(name)
	o.Z = zThings
	o.Sprite = f.art.Building(name)
	o.Transform.Position = r.middle()
	return o
}

// newHouse is a house of the given kind on its cells from
// corner, its door toward the viewer when front is set.
func (f *farm) newHouse(kind int, corner town.Cell, front bool) *game.Object {
	o := game.NewObject("house")
	o.Z = zThings
	o.Sprite = f.art.Building(spriteName("house", kind, front))
	o.Fades = true
	cols, rows := f.art.Footprint(spriteName("house", kind, true))
	o.Transform.Position = rect{corner.Col, corner.Row, cols, rows}.middle()
	return o
}

// spriteName is the name of a house's sprite of the kind, from the front
// or the back.
func spriteName(base string, kind int, front bool) string {
	s := base + itoa(kind)
	if !front {
		s += "_back"
	}
	return s
}

// newWell is a town's well, on the cells r, among the people and the trees.
func (f *farm) newWell(r rect) *game.Object {
	o := game.NewObject("well")
	o.Z = zThings
	o.Sprite = f.art.Building("well")
	o.Transform.Position = r.middle()
	return o
}
