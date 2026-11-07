package main

import (
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/farm"
)

// A building's footprint covers its walls and porch, rounded up to whole
// cells, so what stands on it and what is solid are the building the
// player sees; the houses' bigger than the stylized look's two by two.
func TestFootprintsCoverTheBakedBuildings(t *testing.T) {
	a := &realArt{}
	for name, want := range map[string][2]int{"stall": {3, 2}, "well": {2, 2}, "house0": {8, 9}, "house1": {7, 5}, "house5": {8, 7}} {
		if cols, rows := a.Footprint(name); cols != want[0] || rows != want[1] {
			t.Errorf("%s: %d by %d cells, not %d by %d", name, cols, rows, want[0], want[1])
		}
	}
	if cols, rows := a.Footprint("house0"); cols <= 2 || rows <= 2 {
		t.Errorf("a house of %d by %d cells is no bigger than the stylized one", cols, rows)
	}
}

// A porch's cells are in front of the walls, so the walls' middle, where
// the picture's pivot is, is not the middle of the footprint: the picture
// is shifted along the diagonal by the half of the porch, and the other
// way when the house is turned to face away.
func TestPivotShiftFollowsThePorch(t *testing.T) {
	a := &realArt{}
	dx, dy := a.pivotShift("house0", frontHeading)
	bx, by := a.pivotShift("house0", backHeading)
	if dx >= 0 || dy <= 0 || bx != -dx || by != -dy {
		t.Errorf("front shift %v, %v; back %v, %v", dx, dy, bx, by)
	}
	if sx, sy := a.pivotShift("house1", frontHeading); sx != 0 || sy != 0 {
		t.Errorf("a house with no porch moves by %v, %v", sx, sy)
	}
}

// What stands on a building's cells shows nothing of its own, and a tree's
// name is read into the biome, variant and sway frame the farm asked for.
func TestBlocksAreBlankAndTreesAreRead(t *testing.T) {
	a := newArt(nil).(*realArt)
	cv := a.Tile("block")
	for _, p := range cv.Pixels {
		if p>>24 != 0 {
			t.Fatal("a building's cells draw a picture of their own")
		}
	}
	biome, variant, frame, ok := parseTree("tree_cerrado_7_3")
	if !ok || biome != "cerrado" || variant != 7 || frame != 3 {
		t.Errorf("parsed %q, %d, %d, %v", biome, variant, frame, ok)
	}
	for _, b := range farm.TreeBiomes {
		if _, ok := treeSpecies[b]; !ok {
			t.Errorf("no species for the biome %s", b)
		}
	}
	if _, _, _, ok := parseTree("grass0"); ok {
		t.Error("grass0 read as a tree")
	}
}

// The facings a person is baked milking are those the farm shows them in.
func TestMilkingIsBakedInTheFacingsShown(t *testing.T) {
	if turnOf("right") != bakes.MilkFacingRight || turnOf("left") != bakes.MilkFacingLeft {
		t.Errorf("right %d, left %d", turnOf("right"), turnOf("left"))
	}
}
