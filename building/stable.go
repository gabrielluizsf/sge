package building

// Stable is a town's stable, 5 by 4 m: the same open shed of timber the
// smith works under, with bales of hay stacked at its back, a water trough
// and a counter at the front where the dealer sells.
var Stable = Building{Name: "stable", Width: 5, Depth: 4, build: func(b *builder) {
	shed(b)
	hay := paint{of: Thatch}
	for _, at := range []vec{{-1.5, 0, -1.3}, {-0.4, 0, -1.3}, {0.7, 0, -1.3}, {-1, 0.55, -1.3}, {0.1, 0.55, -1.3}} {
		b.block(box{lo: add(at, vec{-0.5, 0, -0.35}), hi: add(at, vec{0.5, 0.55, 0.35}), paint: hay})
	}
	trough(b, vec{1.3, 0, 0.1})
	boards := paint{of: Planks}
	b.block(box{lo: vec{-1.6, 0, 1.45}, hi: vec{0.4, 0.95, 1.6}, paint: boards})
	b.block(box{lo: vec{-1.7, 0.95, 1.35}, hi: vec{0.5, 1.01, 1.7}, paint: paint{of: Timber}})
}}
