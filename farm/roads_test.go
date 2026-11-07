package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/town"
)

// Walking the roads from the farm's gate finds a town, and from it another.
func TestTheFarmsRoadLeadsOnFromTownToTown(t *testing.T) {
	towns := newTowns(SVGArt())
	towns.Keep = []town.Rect{{C0: -2, R0: -2, C1: 32, R1: 22}}
	paved := func(c, r int) bool {
		if towns.OnRoad(c, r, 0) {
			return true
		}
		tw := towns.At(c, r)
		return tw != nil && tw.Road(c, r)
	}
	seen := map[town.Cell]bool{{Col: 31, Row: pathRow}: true}
	todo := []town.Cell{{Col: 31, Row: pathRow}}
	found := map[*town.Town]bool{}
	for len(todo) > 0 {
		c := todo[0]
		todo = todo[1:]
		if tw := towns.At(c.Col, c.Row); tw != nil && tw.Center == c {
			found[tw] = true
		}
		for _, d := range []town.Cell{{Col: 1}, {Col: -1}, {Row: 1}, {Row: -1}} {
			n := town.Cell{Col: c.Col + d.Col, Row: c.Row + d.Row}
			if !seen[n] && paved(n.Col, n.Row) && abs(n.Col) < 600 && abs(n.Row) < 600 {
				seen[n] = true
				todo = append(todo, n)
			}
		}
	}
	if len(found) < 3 {
		t.Errorf("the road from the farm reaches %d towns", len(found))
	}
}

// Beside a road, the wild lays the fence's tile that runs along that road,
// and a fence turned either way is solid.
func TestTheWildTurnsAFenceToItsRoad(t *testing.T) {
	_, f := build(SVGArt(), idle{})
	if !f.tiles.IsSolid(tileID[fenceDown]) || !isFence(tileID[fenceDown]) || !isFence(tileID["fence"]) {
		t.Fatal("a turned fence is not a solid fence")
	}
	if SVGArt().Tile(fenceDown) == nil {
		t.Fatal("no picture of a turned fence")
	}
	var both [2]int
	for r := -150; r < 150; r++ {
		for c := -150; c < 150; c++ {
			if !f.towns.Fence(c, r) {
				continue
			}
			want := tileID["fence"]
			if f.towns.FenceDown(c, r) {
				want = tileID[fenceDown]
				both[1]++
			} else {
				both[0]++
			}
			if got := f.thingAt(c, r); got != want && isFence(got) {
				t.Errorf("cell %d,%d has tile %d, want %d", c, r, got, want)
			}
		}
	}
	if both[0] == 0 || both[1] == 0 {
		t.Errorf("fences across, down: %v", both)
	}
}
