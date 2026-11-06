package farm

import (
	"github.com/gabrielluizsf/sge/flock"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/vector"
)

// The heights, in screen pixels, birds sit at on what they perch on: a
// fence's rail and a roof.
const (
	fenceTop = 26
	roofTop  = 120
)

// perches are the places birds may sit within radius of c: the tops of
// trees, fences and the roofs of the towns' houses.
func (f *farm) perches(c vector.Vector3, radius float64) []flock.Perch {
	col0, row0 := f.cellOf(c)
	r := int(radius / tile)
	treeTop := 0.35 * float64(f.art.Reach())
	if treeTop == 0 {
		treeTop = 44
	}
	var out []flock.Perch
	for row := row0 - r; row <= row0+r; row++ {
		for col := col0 - r; col <= col0+r; col++ {
			at := pathfind.Cell{Col: col, Row: row}
			switch id := f.thingAt(col, row); {
			case isTree(id):
				out = append(out, flock.Perch{At: f.center(at), Height: treeTop})
			case isFence(id):
				out = append(out, flock.Perch{At: f.center(at), Height: fenceTop})
			}
		}
	}
	for _, t := range f.towns.Near(col0, row0, r) {
		for _, h := range t.Houses {
			out = append(out, flock.Perch{At: houseRect(h).middle(), Height: roofTop})
		}
	}
	return out
}
