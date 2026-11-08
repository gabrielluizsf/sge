package building

import "math"

// chimney is a stack rising through the roof's back slope, near its
// ridge and one end: of the walls' brick or stone, or stone on walls of
// plaster or planks; a band corbelled out under its top, and its flue,
// sooty, at the top.
func (r *raising) chimney() {
	stack := paint{of: Stone}
	if r.h.Wall == Brick {
		stack = paint{of: Brick}
	}
	slope := math.Tan(r.h.Roof.Pitch * math.Pi / 180)
	ridge := r.eaves + min(r.h.Depth, r.h.Width)/2*slope
	x, z := r.h.Width/2-1.1, -0.75
	top := ridge + 0.9
	w, d := 0.32, 0.38
	b := r.b
	b.block(box{lo: vec{x - w, r.eaves - 0.5, z - d}, hi: vec{x + w, top, z + d}, paint: stack})
	b.block(box{lo: vec{x - w - 0.05, top - 0.3, z - d - 0.05}, hi: vec{x + w + 0.05, top - 0.18, z + d + 0.05}, paint: stack})
	b.block(box{lo: vec{x - w + 0.08, top - 0.2, z - d + 0.08}, hi: vec{x + w - 0.08, top + 0.005, z + d - 0.08}, paint: paint{of: Soot}})
}
