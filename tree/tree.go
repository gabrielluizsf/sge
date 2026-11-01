// Package tree grows trees in code, as 3D models to bake: pines, ipês,
// cashew trees, coconut palms and mandacaru cacti, each as a sapling, a
// young tree and a mature one, and each seed a tree of its own shape.
// Trunks and branches are tapered tubes along branching curves, clad in
// the bark texture given them (Poly Haven's, say); foliage, flowers and
// fronds are cards with their textures painted here, cut out where they
// are clear, so nothing in a tree but its bark is anyone's but the
// game's.
//
// A tree stands on its origin, its foot, in metres, Y up, grown to a
// real size.
package tree

import (
	"github.com/gabrielluizsf/sge/gltf"
)

// Size is how far grown a tree is.
type Size int

// The sizes a tree is grown at.
const (
	Sapling Size = iota
	Young
	Mature
	Sizes // how many there are
)

// Spec is the tree to grow: how far grown, the seed that shapes it, and
// the barks it may wear.
type Spec struct {
	Size  Size
	Seed  uint32
	Barks Barks
}

// Species is a kind of tree: its name, the range of metres it grows to
// at each size, measured up its height, or across its width where Wide,
// and how it grows.
type Species struct {
	Name   string
	Metres [Sizes][2]float64
	Wide   bool
	grow   func(g *grower)
}

// Grow is the tree sp says: it grows to a size drawn from the species'
// range for sp.Size, so every seed is a different tree and none leaves
// the range.
func (s Species) Grow(sp Spec) *gltf.Model {
	sp.Size = max(Sapling, min(sp.Size, Mature))
	g := newGrower(sp, sp.Seed*2654435761+uint32(sp.Size)*40503+uint32(len(s.Name))*9973)
	r := s.Metres[sp.Size]
	g.size = g.between(r[0], r[1])
	s.grow(g)
	g.fit(s.Wide)
	return g.done()
}
