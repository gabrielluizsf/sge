package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/town"
)

// props are the objects standing on chunk cx, cy of the wild: the houses
// and stalls of the towns whose cells, their far corner, lie on it.
func (f *farm) props(cx, cy int) []*game.Object {
	const size = 16
	in := func(c town.Cell) bool {
		return floorDiv(c.Col, size) == cx && floorDiv(c.Row, size) == cy
	}
	var objs []*game.Object
	mid := town.Cell{Col: cx*size + size/2, Row: cy*size + size/2}
	for _, t := range f.towns.Near(mid.Col, mid.Row, size/2+town.Reach+2) {
		for _, h := range t.Houses {
			if in(town.Cell{Col: h.Col, Row: h.Row}) {
				objs = append(objs, f.newHouse(h.Kind, town.Cell{Col: h.Col, Row: h.Row}, h.Door.Row > h.Row))
			}
		}
		if in(t.Stall) {
			objs = append(objs, f.newStall(t.Stall))
		}
	}
	return objs
}

func floorDiv(a, b int) int {
	q := a / b
	if a%b != 0 && (a < 0) != (b < 0) {
		q--
	}
	return q
}
