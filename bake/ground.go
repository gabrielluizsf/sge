package bake

import (
	"image/color"

	"github.com/gabrielluizsf/sge/gltf"
)

// Patch is a cell of ground of a tiling material: Span cells its texture
// spans across and down, and which of those cells, Col and Row from 0 to
// Span-1, this is; laid out by the cells' place, patches meet without a
// seam.
type Patch struct {
	Span     int
	Col, Row int
}

// Ground bakes a cell of ground, one unit square, in material mat of m,
// as Patch p says: a tile of the game's map, its diamond wholly opaque so
// tiles meet without a gap, lit by the sun as everything is.
func Ground(m *gltf.Model, mat int, p Patch, o Options) Sprite {
	span := float32(max(p.Span, 1))
	u0, v0 := float32(p.Col)/span, float32(p.Row)/span
	d := 1 / span
	quad := gltf.Primitive{
		// glTF's axes: the cell lies on its ground, Y up, from 0 to 1 in X
		// and Z, its far corner at the origin.
		Positions: [][3]float32{{0, 0, 0}, {1, 0, 0}, {1, 0, -1}, {0, 0, -1}},
		Normals:   [][3]float32{{0, 1, 0}, {0, 1, 0}, {0, 1, 0}, {0, 1, 0}},
		UVs:       [][2]float32{{u0, v0}, {u0 + d, v0}, {u0 + d, v0 + d}, {u0, v0 + d}},
		Indices:   []uint32{0, 1, 2, 0, 2, 3},
		Material:  mat,
	}
	cell := &gltf.Model{
		Roots:     []int{0},
		Nodes:     []gltf.Node{{Mesh: 0, Skin: -1, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}}},
		Meshes:    []gltf.Mesh{{Primitives: []gltf.Primitive{quad}}},
		Materials: m.Materials,
		Images:    m.Images,
	}
	o.Shadow = 0
	// The cell's far corner at the origin, so a map stands the tile on it;
	// turned so its sides run along the game's rows and columns.
	s := Render(cell, Shot{Yaw: 225}, o)
	b := s.Image.Bounds()
	for y := range b.Dy() {
		for x := range b.Dx() {
			if c := s.Image.NRGBAAt(x, y); c.A >= 96 {
				s.Image.SetNRGBA(x, y, color.NRGBA{c.R, c.G, c.B, 255})
			} else {
				s.Image.SetNRGBA(x, y, color.NRGBA{})
			}
		}
	}
	return s
}
