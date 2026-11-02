package building

import "fmt"

// Insides are the pieces a house's inside is laid with, each on one cell
// a metre square about its origin, named as the farm's tiles are: the
// floor of boards, the doormat, the rugs, and each style of Houses' back
// walls, in their order.
func Insides() []Building {
	out := []Building{floorPiece("floor", nil), floorPiece("mat", &spread{w: 0.32, d: 0.22, of: Hessian})}
	for i, m := range rugs {
		out = append(out, floorPiece(fmt.Sprint("rug", i), &spread{w: 0.44, d: 0.44, of: m}))
	}
	for k, h := range houseStyles {
		out = append(out, innerWalls(k, h)...)
	}
	return out
}

// rugs are what each of the rugs is woven of.
var rugs = []Material{Check, Plaid, Jacquard, Linen}

// spread is what lies on a floor: w and d its half width and depth, and
// what it is made of.
type spread struct {
	w, d float64
	of   Material
}

// floorThick is how thick the boards are, under the ground's level; a
// spread lies on them spreadThick thick.
const floorThick, spreadThick = 0.03, 0.012

// floorPiece is a cell of floor of boards running across, a spread on it
// if c is set.
func floorPiece(name string, c *spread) Building {
	return Building{Name: name, Width: 1, Depth: 1, build: func(b *builder) {
		b.block(box{lo: vec{-0.5, -floorThick, -0.5}, hi: vec{0.5, 0, 0.5}, paint: paint{of: Floorboards}, grain: vec{1, 0, 0}})
		if c != nil {
			b.block(box{lo: vec{-c.w, 0, -c.d}, hi: vec{c.w, spreadThick, c.d}, paint: paint{of: c.of}})
		}
	}}
}
