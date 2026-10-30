// Package building builds houses and the things of a farm and a town in
// code, as 3D models to bake: farmhouses, cottages and town houses of a
// few styles, their walls, roofs, doors, windows and chimneys; a market
// stall, a smithy and its anvil, a well, and fences. They are made of
// solids, each closed, clad in tiling PBR textures, such as Poly Haven's,
// at their real size, or in plain colours where there are none.
//
// A building stands on its origin, the middle of its footprint, in
// metres, Y up, its front facing +Z, as glTF models face.
package building

import "github.com/gabrielluizsf/sge/gltf"

// Building is a kind of building or prop: its name, its footprint, the
// metres it covers along X and Z, overhangs aside, and how it is built.
type Building struct {
	Name         string
	Width, Depth float64
	build        func(b *builder)
}

// Model is the building, its materials palette's surfaces.
func (k Building) Model(palette Palette) *gltf.Model {
	b := newBuilder(palette)
	k.build(b)
	return b.done()
}
