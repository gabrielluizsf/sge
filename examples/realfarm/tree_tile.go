package main

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
)

// treeSpecies are the species of tree that grow in each of the farm's
// biomes, by the shape, 0 to 2, of the tree: a mix in the biomes that
// have one, as the land does.
var treeSpecies = map[string][3]string{
	"campo":    {"ipe_green", "ipe", "ipe_green"},
	"pinhal":   {"pine", "pine", "pine"},
	"cerrado":  {"ipe", "cashew", "ipe_green"},
	"tropical": {"coconut", "cashew", "coconut"},
	"caatinga": {"mandacaru", "ipe", "mandacaru"},
}

// treeSway is how far the top of a tree leans, in pixels, at the height
// of the whole tree: the wind moves a tree's crown a little however tall.
const treeSway = 2.5

// parseTree is the biome, variant and frame the tile name, "tree_"
// biome, variant and frame, says; ok is false when it is not a tree's.
func parseTree(name string) (biome string, variant, frame int, ok bool) {
	parts := strings.Split(name, "_")
	if len(parts) != 4 || parts[0] != "tree" {
		return "", 0, 0, false
	}
	_, e1 := fmt.Sscanf(parts[2], "%d", &variant)
	_, e2 := fmt.Sscanf(parts[3], "%d", &frame)
	return parts[1], variant, frame, e1 == nil && e2 == nil
}

// treeTile is the tree of the tile name: that biome's species of the size
// and shape its variant says, baked, at the size it really grows, standing
// on its cell and leaning as the frame of its sway says. Nil while it is
// not baked, or when name is no tree's.
func (a *realArt) treeTile(name string) *canvas.Canvas {
	biome, variant, k, ok := parseTree(name)
	species, known := treeSpecies[biome]
	if !ok || !known || variant < 0 || variant >= farm.TreeVariants {
		return nil
	}
	b, ok := a.bake(fmt.Sprintf("tree_%s_s%d_%d_0", species[variant%3], variant/3, variant%3))
	if !ok {
		return nil
	}
	cv := b.standing()
	foot := cv.Height - farm.TileW/4
	return swayAt(cv, swing{px: swayOf(k) * treeSway / 1.5, height: float64(foot), foot: foot})
}

// standing is b's picture in a frame that stands it on a cell, its pivot,
// the tree's foot, on the diamond's middle, the frame as large as the
// picture needs: the tile set places a larger picture by its own bottom
// middle.
func (b *baked) standing() *canvas.Canvas {
	w := max(farm.TileW, 2*max(b.at.PivotX, b.cv.Width-b.at.PivotX))
	h := max(farm.TileH, b.at.PivotY+farm.TileW/4)
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		panic(err)
	}
	b.put(cv, [2]int{w / 2, h - farm.TileW/4})
	return cv
}
