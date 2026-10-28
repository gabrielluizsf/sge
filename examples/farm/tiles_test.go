package main

import "testing"

// A tile's soft pixels, as a tree's shadow, stay as soft in the tile set
// as in its drawing: the sheet they are cut from is copied onto, not
// blended over, which would make them opaque.
func TestTilesKeepTheirSoftPixels(t *testing.T) {
	tiles := newTiles()
	want, got := tileLook("tree0"), tiles.Tiles[tileID["tree0"]].Image
	for y := range want.Height {
		for x := range want.Width {
			if w, g := want.At(x, y), got.At(x, y); w != g {
				t.Fatalf("pixel %d,%d of tree0 is %v in the tile set, %v in its drawing", x, y, g, w)
			}
		}
	}
}
