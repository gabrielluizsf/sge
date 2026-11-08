package building

// A fence is laid a cell at a time, each piece a metre: a post in the
// middle of its cell, rails running from it, at two heights, to the
// cell's edges, where the next piece's meet them. Turned, a piece runs
// or turns any way.
var (
	// FenceStraight runs along X, through its cell.
	FenceStraight = Building{Name: "fence_straight", Width: 1, Depth: 0.12, build: func(b *builder) {
		b.fence([][2]vec{{{-0.5, 0, 0}, {0.5, 0, 0}}})
	}}
	// FenceCorner turns from +X to +Z.
	FenceCorner = Building{Name: "fence_corner", Width: 0.55, Depth: 0.55, build: func(b *builder) {
		b.fence([][2]vec{{{-0.05, 0, 0}, {0.5, 0, 0}}, {{0, 0, -0.05}, {0, 0, 0.5}}})
	}}
	// FenceEnd ends a fence coming from +X.
	FenceEnd = Building{Name: "fence_end", Width: 0.55, Depth: 0.12, build: func(b *builder) {
		b.fence([][2]vec{{{-0.05, 0, 0}, {0.5, 0, 0}}})
	}}
)

// fence is a fence's post, split oak, its top cut to a point to shed the
// rain, and its rails along each of runs, at the knee and at the hip.
func (b *builder) fence(runs [][2]vec) {
	wood := paint{of: Timber}
	b.block(box{lo: vec{-0.06, 0, -0.06}, hi: vec{0.06, 1.0, 0.06}, paint: wood, grain: up})
	b.raiseTaper(taper{at: vec{0, 1.0, 0}, h: 0.07, base: [2]float64{0.06, 0.06}, top: [2]float64{0.025, 0.025}, paint: wood})
	for _, run := range runs {
		for _, y := range []float64{0.42, 0.82} {
			b.lay(beam{a: add(run[0], vec{0, y, 0}), b: add(run[1], vec{0, y, 0}), w: 0.05, h: 0.11, paint: wood})
		}
	}
}
