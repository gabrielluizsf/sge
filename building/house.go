package building

import "image/color"

// House is a style of house: how wide its front is and how deep, in
// metres; how many storeys it has and how tall each is; what its walls
// are made of, and the colour plaster is painted; whether a frame of
// dark timber shows on them; its roof; whether a chimney rises from it;
// whether a porch runs along its front; and the colours its window frames
// and its door and shutters are painted, none for bare wood, no shutters
// for none.
type House struct {
	Width, Depth float64
	Storeys      int
	Storey       float64
	Wall         Material
	Tint         color.NRGBA
	Timbered     bool
	Roof         Roof
	Chimney      bool
	Porch        bool
	Trim, Door   color.NRGBA
	Shutters     *color.NRGBA
}

// Roof is a house's roof: what covers it, tiles, thatch or slates;
// whether it is hipped, sloping on all four sides, or gabled, on two,
// its ridge along the front; its pitch, degrees above level; and how far
// its eaves reach past the walls, in metres.
type Roof struct {
	Cover    Material
	Hipped   bool
	Pitch    float64
	Overhang float64
}

// floor is how high a house's ground floor stands on its plinth.
const floor = 0.35

// Building is the house, by name.
func (h House) Building(name string) Building {
	if h.Storey <= 0 {
		h.Storey = 2.8
	}
	h.Storeys = max(h.Storeys, 1)
	return Building{Name: name, Width: h.Width, Depth: h.Depth, build: func(b *builder) { raise(b, h) }}
}

// raising is a house going up: its builder, its style, and the height of
// its eaves, where its walls meet its roof.
type raising struct {
	b     *builder
	h     House
	eaves float64
}

// raise builds house h: its plinth and walls, the timbers on them, its
// doors and windows, its porch, roof and chimney.
func raise(b *builder, h House) {
	r := &raising{b: b, h: h, eaves: floor + float64(h.Storeys)*h.Storey}
	r.walls()
	for i, s := range sides(h.Width, h.Depth) {
		if h.Timbered {
			r.timbers(s)
		}
		r.openings(s, i == 0)
	}
	if h.Porch {
		r.porch()
	}
	r.roof()
	if h.Chimney {
		r.chimney()
	}
}

// wall is the paint of the house's walls.
func (r *raising) wall() paint { return paint{of: r.h.Wall, tint: r.h.Tint} }

// walls are the house's walls, a block from the ground to its eaves, on
// a stone plinth a little wider.
func (r *raising) walls() {
	w, d := r.h.Width/2, r.h.Depth/2
	r.b.block(box{lo: vec{-w - 0.06, 0, -d - 0.06}, hi: vec{w + 0.06, floor, d + 0.06}, paint: paint{of: Stone}})
	r.b.block(box{lo: vec{-w, floor - 0.01, -d}, hi: vec{w, r.eaves, d}, paint: r.wall()})
}
