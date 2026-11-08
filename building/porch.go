package building

import "math"

// porchDepth is how far a porch reaches out of the front.
const porchDepth = 2.2

// porch is a veranda along the house's front: a deck of planks at its
// floor, posts along its edge holding up a lean-to roof of the house's
// covering, under the eaves, and steps up to it in front of the door.
func (r *raising) porch() {
	b := r.b
	w, d := r.h.Width/2-0.2, r.h.Depth/2
	deck := paint{of: Planks}
	b.block(box{lo: vec{-w, 0, d}, hi: vec{w, floor, d + porchDepth}, paint: deck, grain: vec{1, 0, 0}})
	slope := math.Tan(14 * math.Pi / 180)
	top := r.eaves - 0.25
	front := top - (porchDepth+0.3)*slope
	posts := max(2, int(2*w/2.6)+1)
	for i := range posts {
		x := -w + 0.12 + float64(i)*(2*w-0.24)/float64(posts-1)
		z := d + porchDepth - 0.12
		b.block(box{lo: vec{x - 0.08, floor, z - 0.08}, hi: vec{x + 0.08, front + 0.05, z + 0.08}, paint: paint{of: Timber}})
	}
	b.lay(beam{a: vec{-w, front, d + porchDepth - 0.12}, b: vec{w, front, d + porchDepth - 0.12}, w: 0.16, h: 0.18, paint: paint{of: Timber}})
	cover := paint{of: r.h.Roof.Cover}
	t := thickness[cover.of] * 0.7
	end := d + porchDepth + 0.3
	b.convex(hexa([8]vec{
		{-w - 0.2, top, d}, {w + 0.2, top, d}, {w + 0.2, front, end}, {-w - 0.2, front, end},
		{-w - 0.2, top + t, d}, {w + 0.2, top + t, d}, {w + 0.2, front + t, end}, {-w - 0.2, front + t, end},
	}, cover))
	x := r.doorX()
	stone := paint{of: Stone}
	b.block(box{lo: vec{x - 0.8, 0, d + porchDepth}, hi: vec{x + 0.8, floor * 0.5, d + porchDepth + 0.6}, paint: stone})
	b.block(box{lo: vec{x - 0.8, 0, d + porchDepth}, hi: vec{x + 0.8, floor, d + porchDepth + 0.3}, paint: stone})
}

// doorX is how far along the front, from its middle, the door is.
func (r *raising) doorX() float64 {
	at := bays(sides(r.h.Width, r.h.Depth)[0])
	return at[doorBay(len(at))]
}
