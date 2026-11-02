package bakes

import (
	"fmt"
	"image"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/tree"
)

// TreeSeeds are how many trees of each species and size are baked, so a
// wood of them is not one tree over and over.
const TreeSeeds = 3

// FarmTree is the species the farm's trees are, as saplings: the green
// ipê, of the fields.
const FarmTree = "ipe_green"

// TreeHeadings is how many of the headings a tree is baked at: the game
// stands each tree at one, and trees differ by seed.
const TreeHeadings = 1

// barkIDs are the assets, by id, the species' barks are read from.
var barkIDs = struct{ Pine, Smooth, Palm string }{"pine_bark", "bark_platanus", "palm_tree_bark"}

// treeSpot is one tree to bake: the species, how far grown, the seed.
type treeSpot struct {
	species tree.Species
	size    tree.Size
	seed    int
}

// TreeTask is the name of the task that bakes the tree of the species,
// size and seed.
func TreeTask(species string, size tree.Size, seed int) string {
	return fmt.Sprintf("tree %s %d %d", species, size, seed)
}

// treeTasks bake every species package tree grows, at every size, a few
// trees of each, at its headings, clad in the barks the manifest pins:
// tree_<species>_s<size>_<seed>_<heading>.
func (b *Baker) treeTasks() []Task {
	var out []Task
	for _, s := range tree.All() {
		for size := tree.Sapling; size < tree.Sizes; size++ {
			for seed := range TreeSeeds {
				t := treeSpot{s, size, seed}
				out = append(out, Task{Name: TreeTask(s.Name, size, seed), Kind: "tree", ID: s.Name, Core: true, Run: func() error { return b.tree(t) }})
			}
		}
	}
	return out
}

// tree bakes the tree t at its headings, a core each.
func (b *Baker) tree(t treeSpot) error {
	barks, err := b.Barks()
	if err != nil {
		return err
	}
	model := t.species.Grow(tree.Spec{Size: t.size, Seed: uint32(t.seed + 1), Barks: barks})
	return b.each(TreeHeadings, func(i int) error {
		sp := bake.Render(model, bake.Shot{Yaw: Headings[i]}, Options())
		return b.Shelf.Put(fmt.Sprintf("tree_%s_s%d_%d_%d", t.species.Name, t.size, t.seed, i), sp)
	})
}

// Barks are the barks the manifest pins, read from their glTF once.
func (b *Baker) Barks() (tree.Barks, error) {
	b.barkOnce.Do(func() { b.barks, b.barkErr = barksOf(b.Manifest) })
	return b.barks, b.barkErr
}

// barksOf reads the barks m pins.
func barksOf(m kit.Manifest) (tree.Barks, error) {
	var out tree.Barks
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
		model, err := Load(a)
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
			return AsNRGBA(model.Images[i])
		}
		*dst = tree.Bark{Albedo: img(g.Texture), ARM: img(g.MetalRough), Normal: img(g.Normal), Span: a.Span}
	}
	return out, nil
}
