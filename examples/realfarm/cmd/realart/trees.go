package main

import (
	"fmt"
	"image"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/tree"
)

// treeSeeds are how many trees of each species and size are baked, so a
// wood of them is not one tree over and over.
const treeSeeds = 3

// barkIDs are the assets, by id, the species' barks are read from.
var barkIDs = struct{ Pine, Smooth, Palm string }{"pine_bark", "bark_platanus", "palm_tree_bark"}

// bakeTrees bakes every species package tree grows, at every size, a few
// trees of each, at every heading, clad in the barks the manifest pins,
// as j says; j.only names species, commas between, or "" for all.
func bakeTrees(j batch, sprites map[string]sprite) {
	barks, err := barksOf(j.manifest)
	if err != nil {
		fail(err)
	}
	for _, s := range tree.All() {
		if j.only != "" && !slices.Contains(strings.Split(j.only, ","), s.Name) {
			continue
		}
		var wg sync.WaitGroup
		var mu sync.Mutex
		for size := tree.Sapling; size < tree.Sizes; size++ {
			for seed := range treeSeeds {
				model := s.Grow(tree.Spec{Size: size, Seed: uint32(seed + 1), Barks: barks})
				// Its headings are baked at once, a core each.
				for i, yaw := range headings {
					wg.Add(1)
					go func() {
						defer wg.Done()
						sp := bake.Render(model, bake.Shot{Yaw: yaw}, options())
						name := fmt.Sprintf("tree_%s_s%d_%d_%d", s.Name, size, seed, i)
						if err := writePNG(filepath.Join(j.out, name+".png"), sp); err != nil {
							fail(err)
						}
						b := sp.Image.Bounds()
						mu.Lock()
						sprites[name] = sprite{W: b.Dx(), H: b.Dy(), PivotX: sp.Pivot.X, PivotY: sp.Pivot.Y}
						mu.Unlock()
					}()
				}
				wg.Wait()
			}
		}
		fmt.Println(s.Name, "grown")
	}
}

// barksOf are the barks the manifest pins, read from their glTF.
func barksOf(m kit.Manifest) (tree.Barks, error) {
	var out tree.Barks
	fe := kit.Fetcher()
	for _, a := range m.Assets {
		var dst *tree.Bark
		switch a.ID {
		case barkIDs.Pine:
			dst = &out.Pine
		case barkIDs.Smooth:
			dst = &out.Smooth
		case barkIDs.Palm:
			dst = &out.Palm
		default:
			continue
		}
		model, err := load(fe.Cache, a)
		if err != nil {
			return out, fmt.Errorf("%s: %w", a.ID, err)
		}
		if len(model.Materials) == 0 {
			return out, fmt.Errorf("%s has no material", a.ID)
		}
		g := model.Materials[0]
		img := func(i int) image.Image {
			if i < 0 || i >= len(model.Images) {
				return nil
			}
			return asNRGBA(model.Images[i])
		}
		*dst = tree.Bark{Albedo: img(g.Texture), ARM: img(g.MetalRough), Normal: img(g.Normal), Span: a.Span}
	}
	return out, nil
}
